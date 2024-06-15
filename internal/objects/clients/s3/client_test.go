package s3

import (
	"context"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/queue/clients"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	conf := config.Configuration{
		S3BaseUrl: "http://localhost:9000",
	}

	ctx := context.Background()
	client, err := NewClient(conf, ctx)
	assert.NotNil(t, client)
	assert.Nil(t, err)
}

// TestClient_ManipulateObject PutObject関数にてS3オブジェクトをアップロード、削除するテスト
func TestClient_ManipulateObject(t *testing.T) {
	conf := config.Configuration{
		S3BaseUrl: "http://127.0.0.1:9000",
		Region:    "ap-northeast-1",
	}

	ctx := context.Background()
	client, err := NewClient(conf, ctx)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	t.Run("PutObject", func(t *testing.T) {
		objectBody := []byte("test")
		s3Object := clients.S3Object{
			Bucket: "test",
			Key:    "test",
		}
		err = client.PutObject(ctx, objectBody, s3Object)
		assert.Nil(t, err)
	})

	t.Run("DeleteObject", func(t *testing.T) {
		s3Object := clients.S3Object{
			Bucket: "test",
			Key:    "test",
		}
		err = client.DeleteObject(ctx, s3Object)
		assert.Nil(t, err)
	})
}
