package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/logging"
	"github.com/Fufuhu/go-clamav/internal/queue/clients"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsSqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"go.uber.org/zap"
	"time"
)

type Client struct {
	conf    config.Configuration
	service *awsSqs.Client
}

// Poll SQSにポーリングする。processには、S3Objectをどう処理するかを表す関数を渡す
func (c *Client) Poll(ctx context.Context, process func(clients.QueueMessageInterface, context.Context) error) error {
	logger := logging.GetLogger()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)

	for {
		messages, err := c.ReceiveMessages(ctx)
		if err != nil {
			logger.Warn("SQSのメッセージ取得に失敗しました")
			logger.Error(err.Error())
			continue
		}
		logger.Info(fmt.Sprintf("%d個のメッセージを取得しました", len(messages)))
		for _, message := range messages {
			logger.Info("個別メッセージの処理を開始します")
			err = process(message, ctx)
			if err != nil {
				logger.Warn("SQSのメッセージ処理に失敗しました",
					zap.String("Bucket", message.GetBucket()),
					zap.String("Key", message.GetKey()))
				logger.Error(err.Error())
				continue
			}

			// SQSのメッセージ処理に成功したらメッセージを削除する
			if err = c.DeleteMessage(ctx, message.GetReceiptHandle()); err != nil {
				logger.Warn("SQSメッセージの削除に失敗しました",
					zap.String("ReceiptHandle", message.GetReceiptHandle()))
				logger.Error(err.Error())
				continue
			}
			logger.Info("SQSメッセージの削除が完了しました")
		}

		time.Sleep(5 * time.Second)
	}
}

type S3Event struct {
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

// ReceiveMessages キューからメッセージを取得する
func (c *Client) ReceiveMessages(ctx context.Context) ([]clients.QueueMessageInterface, error) {
	logger := logging.GetLogger()
	defer logger.Sync()

	receiveMessageInput := &awsSqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.conf.QueueURL),
		MaxNumberOfMessages: c.conf.MaxNumberOfMessages,
		WaitTimeSeconds:     c.conf.WaitTimeSeconds,
		VisibilityTimeout:   c.conf.VisibilityTimeout,
	}

	logger.Info("SQSキューからメッセージを取得します")

	result, err := c.service.ReceiveMessage(ctx, receiveMessageInput)

	if err != nil {
		logger.Warn("SQSキューからのメッセージの取得に失敗しました")
		logger.Error(err.Error())
		return []clients.QueueMessageInterface{}, err
	}

	var s3Objects []clients.QueueMessageInterface
	for _, message := range result.Messages {
		logger.Info("SQSメッセージを処理中です",
			zap.String("MessageID", *message.MessageId),
			zap.String("MessageBody", *message.Body))

		var event = S3Event{}
		err = json.Unmarshal([]byte(*message.Body), &event)
		if err != nil {
			logger.Warn("SQSメッセージのjson.Unmarshalに失敗しました")
			logger.Error(err.Error())
			continue
		}

		// eventのRecordsからイベント情報を取り出してQueueMessageのフォーマットにして格納する
		for _, record := range event.Records {

			s3Objects = append(s3Objects, &clients.QueueMessage{
				Bucket:        record.S3.Bucket.Name,
				Key:           record.S3.Object.Key,
				ReceiptHandle: *message.ReceiptHandle,
			})
		}
	}

	return s3Objects, nil
}

// SendMessage メッセージを送信する。送信したメセージIDとエラーを返す
func (c *Client) SendMessage(ctx context.Context, message string) (*string, error) {
	logger := logging.GetLogger()
	defer logger.Sync()

	sendMessageInput := &awsSqs.SendMessageInput{
		QueueUrl:    aws.String(c.conf.QueueURL),
		MessageBody: aws.String(message),
	}

	sendMessageOutput, err := c.service.SendMessage(ctx, sendMessageInput)
	if err != nil {
		logger.Warn("SQSメッセージの送信に失敗しました")
		logger.Error(err.Error())
		return nil, err
	}
	logger.Info("SQSメッセージを送信しました",
		zap.String("MessageID", *sendMessageOutput.MessageId))

	return sendMessageOutput.MessageId, err
}

func (c *Client) DeleteMessage(ctx context.Context, receiptHandle string) error {
	logger := logging.GetLogger()
	defer logger.Sync()

	deleteMessageInput := &awsSqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.conf.QueueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := c.service.DeleteMessage(ctx, deleteMessageInput)
	if err != nil {
		logger.Warn("SQSメッセージの削除に失敗しました",
			zap.String("ReceiptHandle", receiptHandle))
		logger.Error(err.Error())
		return err
	}
	logger.Info("SQSメッセージを削除しました",
		zap.String("ReceiptHandle", receiptHandle))

	return nil
}

func NewClient(conf config.Configuration, ctx context.Context) (*Client, error) {
	logger := logging.GetLogger()
	defer logger.Sync()

	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(conf.Region))
	if err != nil {
		logger.Warn("AWSクライアントの設定作成に失敗しました")
		logger.Error(err.Error())
		return nil, err
	}
	if conf.BaseUrl != "" {
		cfg.BaseEndpoint = aws.String(conf.BaseUrl)
	}
	svc := awsSqs.NewFromConfig(cfg)

	return &Client{
		conf:    conf,
		service: svc,
	}, nil
}
