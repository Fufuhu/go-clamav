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
	conf config.Configuration
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

	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(c.conf.Region))
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	svc := awsSqs.NewFromConfig(cfg)

	receiveMessageInput := &awsSqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.conf.QueueURL),
		MaxNumberOfMessages: c.conf.MaxNumberOfMessages,
		WaitTimeSeconds:     c.conf.WaitTimeSeconds,
	}

	for {
		s3Objects, err := c.ReceiveMessages(ctx, receiveMessageInput, svc)
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
		} `json:"s3"`
		Object struct {
			Key string `json:"key"`
		} `json:"object"`
	} `json:"Records"`
}

// ReceiveMessages キューからメッセージを取得する
func (c *Client) ReceiveMessages(ctx context.Context, receiveMessageInput *awsSqs.ReceiveMessageInput, svc *awsSqs.Client) ([]clients.S3Object, error) {
	logger := logging.GetLogger()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)

	result, err := svc.ReceiveMessage(ctx, receiveMessageInput)

	if err != nil {
		logger.Warn("SQSキューからのメッセージの取得に失敗しました")
		logger.Error(err.Error())
		return []clients.S3Object{}, err
	}

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
		var s3Objects []clients.S3Object
		for _, record := range event.Records {
			s3Objects = append(s3Objects, clients.S3Object{
				Bucket: record.S3.Bucket.Name,
				Key:    record.Object.Key,
			})
		}
	}

	return nil, nil
}

func NewClient(conf config.Configuration) *Client {
	return &Client{
		conf: conf,
	}
}
