package clients

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestS3Object_GetObjectPath GetObjectPath関数にてS3ObjectのBucketとKeyを結合した文字列が返却されることを確認するテスト
func TestQueueMessage_GetObjectPath(t *testing.T) {
	s3Object := QueueMessage{
		Bucket: "test",
		Key:    "object",
	}

	expected := "s3://test/object"
	actual := s3Object.GetObjectPath()
	assert.Equal(t, expected, actual)
}
