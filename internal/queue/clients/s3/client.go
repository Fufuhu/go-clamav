package s3

import (
	"bytes"
	"context"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/logging"
	"github.com/Fufuhu/go-clamav/internal/queue/clients"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	conf    config.Configuration
	service *awsS3.Client
}

// PutObject PutObject関数はS3オブジェクトをアップロードする
func (c *Client) PutObject(ctx context.Context, objectBody []byte, s3Object clients.S3Object) error {
	logger := logging.GetLogger()
	defer logger.Sync()

	_, err := c.service.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket: aws.String(s3Object.Bucket),
		Key:    aws.String(s3Object.Key),
		Body:   bytes.NewReader(objectBody),
	})
	if err != nil {
		logger.Warn("S3オブジェクトのアップロードに失敗しました")
		logger.Error(err.Error())
		return err
	}
	return nil
}

// NewClient NewClient関数はS3クライアントを生成する
func NewClient(conf config.Configuration, ctx context.Context) (*Client, error) {
	logger := logging.GetLogger()
	defer logger.Sync()

	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(conf.Region))
	if err != nil {
		logger.Warn("AWS S3クライアントの設定作成に失敗しました")
	}
	if conf.S3BaseUrl != "" {
		cfg.BaseEndpoint = aws.String(conf.S3BaseUrl)
	}

	return &Client{
		conf:    conf,
		service: awsS3.NewFromConfig(cfg),
	}, nil
}
