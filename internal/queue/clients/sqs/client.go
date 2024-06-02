package sqs

import (
	"context"
	"encoding/json"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/logging"
	"github.com/Fufuhu/go-clamav/internal/queue/clients"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsSqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"go.uber.org/zap"
)

type Client struct {
	conf    config.Configuration
	service *awsSqs.Client
}

// Poll SQSにポーリングする
func (c *Client) Poll(ctx context.Context, process func(clients.S3Object) error) error {
	logger := logging.GetLogger()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)

	for {
		s3Objects, err := c.ReceiveMessages(ctx)
		if err != nil {
			logger.Warn("SQSのメッセージ取得に失敗しました")
			logger.Error(err.Error())
			continue
		}
		for _, s3Object := range s3Objects {
			err = process(s3Object)
			if err != nil {
				logger.Warn("S3から取得したオブジェクトの処理に失敗しました")
				logger.Error(err.Error())
				continue
			}
		}
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
func (c *Client) ReceiveMessages(ctx context.Context) ([]clients.S3Object, error) {
	logger := logging.GetLogger()
	defer logger.Sync()

	receiveMessageInput := &awsSqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.conf.QueueURL),
		MaxNumberOfMessages: c.conf.MaxNumberOfMessages,
		WaitTimeSeconds:     c.conf.WaitTimeSeconds,
	}

	result, err := c.service.ReceiveMessage(ctx, receiveMessageInput)

	if err != nil {
		logger.Warn("SQSキューからのメッセージの取得に失敗しました")
		logger.Error(err.Error())
		return []clients.S3Object{}, err
	}

	var s3Objects []clients.S3Object
	for _, message := range result.Messages {
		logger.Info("SQSメッセージを処理中です",
			zap.String("MessageID", *message.MessageId),
			zap.String("MessageBody", *message.Body))

		var event = &S3Event{}
		err = json.Unmarshal([]byte(*message.Body), event)
		if err != nil {
			logger.Warn("SQSメッセージのjson.Unmarshalに失敗しました")
			logger.Error(err.Error())
			continue
		}

		// eventのRecordsからイベント情報を取り出してQueueMessageのフォーマットにして格納する
		for _, record := range event.Records {
			s3Objects = append(s3Objects, clients.S3Object{
				Bucket: record.S3.Bucket.Name,
				Key:    record.S3.Object.Key,
			})
		}

		deleteMessageInput := &awsSqs.DeleteMessageInput{
			QueueUrl:      aws.String(c.conf.QueueURL),
			ReceiptHandle: message.ReceiptHandle,
		}

		_, err = c.service.DeleteMessage(ctx, deleteMessageInput)
		if err != nil {
			logger.Warn("SQSメッセージの削除に失敗しました",
				zap.String("MessageID", *message.MessageId),
				zap.String("ReceiptHandle", *message.ReceiptHandle))
			logger.Error(err.Error())
			continue
		}
		logger.Info("SQSメッセージを削除しました",
			zap.String("MessageID", *message.MessageId),
			zap.String("ReceiptHandle", *message.ReceiptHandle))
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
