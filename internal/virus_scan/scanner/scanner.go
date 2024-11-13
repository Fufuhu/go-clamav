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
	"github.com/Fufuhu/go-clamav/config"
)

type Scanner struct {
	dynamodbClient *dynamodb.Client
	s3Client       *s3.Client
	clamdClient    *clamav.Client
}

func (s *Scanner) Scan() (model.ScanResult, error) {
	return model.ScanResult{}, nil
}

func (s *Scanner) Process(message clients.QueueMessageInterface, conf config.Configuration, ctx context.Context) error {
	logger := logging.GetLogger()
	defer logger.Sync()

	// 対象外ファイルをスキップ
	if ok, err := message.IsTargetFile(conf); err != nil {
		logger.Warn("ファイルの対象外判定に失敗しました")
		logger.Error(err.Error())
		return err
	} else if !ok {
		logger.Info("ファイルが対象外のためスキップします",
			zap.String("Bucket", message.GetBucket()),
			zap.String("Key", message.GetKey()))
		return nil
	}

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

	logger.Info("ファイルの取得に成功しました")
	logger.Info("ファイルのスキャンを開始します")

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

	_, err = s.dynamodbClient.PutScanResult(ctx, scanResult)
	if err != nil {
		logger.Error("スキャン結果の保存に失敗しました",
			zap.String("Bucket", message.GetBucket()),
			zap.String("Key", message.GetKey()))
		logger.Error(err.Error())
		return err
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
