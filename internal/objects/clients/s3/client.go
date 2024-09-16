package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/logging"
	"github.com/Fufuhu/go-clamav/internal/queue/clients"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
	"io"
)

type Client struct {
	conf    config.Configuration
	service *awsS3.Client
}

// PutObject PutObject関数はS3オブジェクトをアップロードする
func (c *Client) PutObject(ctx context.Context, objectBody []byte, s3Object clients.QueueMessageInterface) error {
	logger := logging.GetLogger()
	defer logger.Sync()

	_, err := c.service.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket: aws.String(s3Object.GetBucket()),
		Key:    aws.String(s3Object.GetKey()),
		Body:   bytes.NewReader(objectBody),
	})
	if err != nil {
		logger.Warn("S3オブジェクトのアップロードに失敗しました")
		logger.Error(err.Error())
		return err
	}
	return nil
}

// DeleteObject DeleteObject関数はS3オブジェクトを削除する
func (c *Client) DeleteObject(ctx context.Context, s3Object clients.QueueMessageInterface) error {
	logger := logging.GetLogger()
	defer logger.Sync()

	_, err := c.service.DeleteObject(ctx, &awsS3.DeleteObjectInput{
		Bucket: aws.String(s3Object.GetBucket()),
		Key:    aws.String(s3Object.GetKey()),
	})
	if err != nil {
		logger.Warn("S3オブジェクトの削除に失敗しました")
		logger.Error(err.Error())
		return err
	}
	return nil
}

// GetObject GetObject関数はS3オブジェクトを取得する
func (c *Client) GetObject(ctx context.Context, s3Object clients.QueueMessageInterface) (io.ReadCloser, error) {
	// 本当なら取得したオブジェクトのボディを[]byteで取りたいが、メモリ上にすべて展開するのは安全ではないので
	// io.ReadCloserを渡すようにしている
	logger := logging.GetLogger()
	defer logger.Sync()

	getObjectOutput, err := c.service.GetObject(ctx, &awsS3.GetObjectInput{
		Bucket: aws.String(s3Object.GetBucket()),
		Key:    aws.String(s3Object.GetKey()),
	})

	if err != nil {
		logger.Warn("S3オブジェクトの取得に失敗しました",
			zap.String("bucket", s3Object.GetBucket()),
			zap.String("key", s3Object.GetKey()))
		return nil, err
	}

	if getObjectOutput == nil {
		logger.Warn("S3オブジェクトが取得できません",
			zap.String("bucket", s3Object.GetBucket()),
			zap.String("key", s3Object.GetKey()))
		return nil, errors.New("s3オブジェクトが取得できず、オブジェクトボディを返せません")
	}

	logger.Info("S3オブジェクトを取得しました",
		zap.String("bucket", s3Object.GetBucket()),
		zap.String("key", s3Object.GetKey()))

	return getObjectOutput.Body, nil
}

// NewClient NewClient関数はS3クライアントを生成する
func NewClient(conf config.Configuration, ctx context.Context) (*Client, error) {
	logger := logging.GetLogger()
	defer logger.Sync()

	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(conf.Region))

	var svc *awsS3.Client

	if err != nil {
		logger.Warn("AWS S3クライアントの設定作成に失敗しました")
	}
	if conf.S3BaseUrl != "" {
		logger.Info("S3のカスタムエンドポイントを設定します", zap.String("S3BaseUrl", conf.S3BaseUrl))
		logger.Info("S3 Base URLを設定します", zap.String("S3BaseUrl", conf.S3BaseUrl))
		cfg.BaseEndpoint = aws.String(conf.S3BaseUrl)
		svc = awsS3.NewFromConfig(cfg, func(o *awsS3.Options) {
			o.BaseEndpoint = aws.String(conf.S3BaseUrl)
			o.UsePathStyle = true
		})
		logger.Info(fmt.Sprintf("s3クライアントのエンドポイントは、%sです", *cfg.BaseEndpoint))
	} else {
		svc = awsS3.NewFromConfig(cfg)
	}
	logger.Info(fmt.Sprintf("クライアントのリージョンは、%sです", cfg.Region))

	return &Client{
		conf:    conf,
		service: svc,
	}, nil
}
