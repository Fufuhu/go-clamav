package scanner

import (
	"github.com/Fufuhu/go-clamav/internal/db/clients/dynamodb"
	"github.com/Fufuhu/go-clamav/internal/model"
	"github.com/Fufuhu/go-clamav/internal/objects/clients/s3"
	"github.com/Fufuhu/go-clamav/internal/virus_scan/clients/clamav"
)

type Scanner struct {
	dynamodbClient *dynamodb.Client
	s3Client       *s3.Client
	clamdClient    *clamav.Client
}

func (s *Scanner) Scan() (model.ScanResult, error) {
	return model.ScanResult{}, nil
}

func NewScanner(dynamodbClient *dynamodb.Client, s3Client *s3.Client, clamdClient *clamav.Client) Scanner {
	return Scanner{
		dynamodbClient: dynamodbClient,
		s3Client:       s3Client,
		clamdClient:    clamdClient,
	}
}
