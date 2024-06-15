package s3

import (
	"context"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	conf := config.Configuration{
		S3BaseUrl: "http://localhost:9001",
	}

	ctx := context.Background()
	client, err := NewClient(conf, ctx)
	assert.NotNil(t, client)
	assert.Nil(t, err)
}
