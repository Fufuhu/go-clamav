package scanner

import (
	"context"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/db/clients/dynamodb"
	"github.com/Fufuhu/go-clamav/internal/objects/clients/s3"
	"github.com/Fufuhu/go-clamav/internal/virus_scan/clients/clamav"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewScanner(t *testing.T) {
	cfg, err := config.GetConfig()
	assert.Nil(t, err)
	assert.NotNil(t, cfg)

	ctx := context.Background()
	s3Client, err := s3.NewClient(*cfg, ctx)
	assert.Nil(t, err)
	assert.NotNil(t, s3Client)

	dynamoClient, err := dynamodb.NewClient(*cfg, ctx)
	assert.Nil(t, err)
	assert.NotNil(t, dynamoClient)

	clamdClient := clamav.NewClient(*cfg)
	assert.NotNil(t, clamdClient)

	scanner := NewScanner(dynamoClient, s3Client, clamdClient)
	assert.NotNil(t, scanner)
}
