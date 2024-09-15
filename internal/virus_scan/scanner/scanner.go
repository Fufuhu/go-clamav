package scanner

import (
	"context"
	"github.com/Fufuhu/go-clamav/internal/db/clients/dynamodb"
	"github.com/Fufuhu/go-clamav/internal/logging"
	"github.com/Fufuhu/go-clamav/internal/model"
	"github.com/Fufuhu/go-clamav/internal/objects/clients/s3"
	"github.com/Fufuhu/go-clamav/internal/queue/clients"
	"github.com/Fufuhu/go-clamav/internal/virus_scan/clients/clamav"
	"go.uber.org/zap"
	"io"
	"time"
)

type Scanner struct {
	dynamodbClient *dynamodb.Client
	s3Client       *s3.Client
	clamdClient    *clamav.Client
}

func (s *Scanner) Scan() (model.ScanResult, error) {
	return model.ScanResult{}, nil
}

func (s *Scanner) Process(message clients.QueueMessageInterface, ctx context.Context) error {
	logger := logging.GetLogger()
	defer logger.Sync()

	file, err := s.s3Client.GetObject(ctx, message)
	if err != nil {
		logger.Warn("S3オブジェクトの取得に失敗しました")
		logger.Error(err.Error())
		return err
	}
	defer func(file io.ReadCloser) {
		err := file.Close()
		if err != nil {
			logger.Error("ファイルのクローズに失敗しました")
			logger.Error(err.Error())
		}
	}(file)

	result, err := s.clamdClient.Scan(file)
	if err != nil {
		logger.Warn("ファイルのスキャンに失敗しました")
		logger.Error(err.Error())
		return err
	}

	scanResult := &model.ScanResult{}
	scanResult.Key = message.GetKey()
	scanResult.Bucket = message.GetBucket()
	scanResult.ReceiptHandle = message.GetReceiptHandle()
	scanResult.ScannedAt = time.Now().Format(time.RFC3339)

	if result.Message == clamav.ResultOK {
		logger.Info("ファイルは感染していません",
			zap.String("Bucket", message.GetBucket()),
			zap.String("Key", message.GetKey()))
		scanResult.ScanResult = model.ScanResultClean
	} else {
		logger.Info("ファイルは感染しています",
			zap.String("Bucket", message.GetBucket()),
			zap.String("Key", message.GetKey()),
			zap.String("Message", result.Message),
		)
		scanResult.ScanResult = model.ScanResultInfected
	}

	return nil
}

func NewScanner(dynamodbClient *dynamodb.Client, s3Client *s3.Client, clamdClient *clamav.Client) Scanner {
	return Scanner{
		dynamodbClient: dynamodbClient,
		s3Client:       s3Client,
		clamdClient:    clamdClient,
	}
}
