package dynamodb

import (
	"context"

	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/logging"
	"github.com/Fufuhu/go-clamav/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
)

type Client struct {
	conf    config.Configuration
	service *awsDynamodb.Client
}

// PutScanResult PutScanResult関数はスキャン結果をDynamoDBに追加する
func (c *Client) PutScanResult(ctx context.Context, result *model.ScanResult) (*model.ScanResult, error) {
	logger := logging.GetLogger()
	defer logger.Sync()

	objectPath := result.GetObjectPath()

	item := map[string]types.AttributeValue{}

	item["ObjectPath"] = &types.AttributeValueMemberS{
		Value: objectPath,
	}

	item["Status"] = &types.AttributeValueMemberS{
		Value: result.ScanResult,
	}

	item["ScannedAt"] = &types.AttributeValueMemberS{
		Value: result.ScannedAt,
	}

	item["Message"] = &types.AttributeValueMemberS{
		Value: result.Message,
	}

	_, err := c.service.PutItem(ctx, &awsDynamodb.PutItemInput{
		TableName: aws.String(c.conf.DynamoDBTable),
		Item:      item,
	})

	if err != nil {
		logger.Warn("DynamoDBへのスキャン結果の追加に失敗しました",
			zap.String("ObjectPath", objectPath))
		logger.Error(err.Error())
		return nil, err
	}

	// 感染していた場合
	if result.ScanResult == model.ScanResultInfected {
		logger.Info("感染しているファイルが検出されました",
			zap.String("ObjectPath", objectPath))

		// 情報を追加
		_, err := c.service.PutItem(ctx, &awsDynamodb.PutItemInput{
			TableName: aws.String(c.conf.DynamoDBTableInfected),
			Item:      item,
		})

		if err != nil {
			logger.Warn("DynamoDBへの感染ファイル情報の追加に失敗しました",
				zap.String("ObjectPath", objectPath))
			logger.Error(err.Error())
			return nil, err
		}
	}

	return result, nil
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
