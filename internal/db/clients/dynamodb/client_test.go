package dynamodb

import (
	"context"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/model"
	"github.com/Fufuhu/go-clamav/internal/queue/clients"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestPutScanResult PutScanResult関数にてスキャン結果をDynamoDBに追加するテスト
func TestPutScanResult(t *testing.T) {
	conf := config.Configuration{
		DynamoDBBaseUrl:       "http://localhost:8000",
		DynamoDBTable:         "ScanResults",
		DynamoDBTableInfected: "InfectedScanResults",
	}

	ctx := context.Background()
	client, err := NewClient(conf, ctx)
	assert.Nil(t, err)

	t.Run("clean", func(t *testing.T) {

		result := &model.ScanResult{
			S3Object: clients.S3Object{
				Bucket: "test-bucket",
				Key:    "test-key-clean",
			},
			ScanResult: model.ScanResultClean,
			ScannedAt:  "2021-01-01T00:00:00Z",
		}

		result, err = client.PutScanResult(ctx, result)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("infected", func(t *testing.T) {

		result := &model.ScanResult{
			S3Object: clients.S3Object{
				Bucket: "test-bucket",
				Key:    "test-key-infected",
			},
			ScanResult: model.ScanResultInfected,
			ScannedAt:  "2021-01-01T00:00:00Z",
		}

		result, err = client.PutScanResult(ctx, result)
		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

}

// TestNewClient NewClient関数にてDynamoDBクライアントを生成するテスト
func TestNewClient(t *testing.T) {
	conf := config.Configuration{
		DynamoDBBaseUrl: "http://localhost:8000",
		DynamoDBTable:   "ScanResults",
	}

	ctx := context.Background()
	client, err := NewClient(conf, ctx)
	assert.Nil(t, err)
	assert.NotNil(t, client)
}
