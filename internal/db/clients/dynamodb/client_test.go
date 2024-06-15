package dynamodb

import (
	"context"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
