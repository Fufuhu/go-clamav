package dynamodb

import (
	"context"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/logging"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Client struct {
	conf    config.Configuration
	service *awsDynamodb.Client
}

func NewClient(conf config.Configuration, ctx context.Context) (*Client, error) {
	logger := logging.GetLogger()
	defer logger.Sync()

	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(conf.Region))
	if err != nil {
		logger.Warn("DynamoDBクライアントの生成に失敗しました")
	}
	if conf.DynamoDBBaseUrl != "" {
		cfg.BaseEndpoint = &conf.DynamoDBBaseUrl
	}

	return &Client{
		conf:    conf,
		service: awsDynamodb.NewFromConfig(cfg),
	}, nil
}
