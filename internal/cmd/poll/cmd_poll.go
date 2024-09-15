package poll

import (
	"context"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/cmd"
	"github.com/Fufuhu/go-clamav/internal/db/clients/dynamodb"
	"github.com/Fufuhu/go-clamav/internal/logging"
	"github.com/Fufuhu/go-clamav/internal/objects/clients/s3"
	"github.com/Fufuhu/go-clamav/internal/queue/clients/sqs"
	"github.com/Fufuhu/go-clamav/internal/virus_scan/clients/clamav"
	"github.com/Fufuhu/go-clamav/internal/virus_scan/scanner"
	"github.com/spf13/cobra"
)

type CommandPoll struct {
}

func (p *CommandPoll) Run(cmd *cobra.Command, args []string) {
	logger := logging.GetLogger()
	defer logger.Sync()

	ctx := context.Background()
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Error("設定ファイルの読み込みに失敗しました")
		logger.Error(err.Error())
		panic(err)
	}

	sqsClient, err := sqs.NewClient(*cfg, ctx)
	if err != nil {
		logger.Error("SQSクライアントの作成に失敗しました")
		logger.Error(err.Error())
		panic(err)
	}

	// S3クライアントの作成
	s3Client, err := s3.NewClient(*cfg, ctx)
	if err != nil {
		logger.Error("S3クライアントの作成に失敗しました")
		logger.Error(err.Error())
		panic(err)
	}

	// DynamoDBクライアントの作成
	dynamoClient, err := dynamodb.NewClient(*cfg, ctx)
	if err != nil {
		logger.Error("DynamoDBクライアントの作成に失敗しました")
		logger.Error(err.Error())
		panic(err)
	}

	clamdClient := clamav.NewClient(*cfg)

	virusScanner := scanner.NewScanner(dynamoClient, s3Client, clamdClient)

	// エラーは拾ったところでどうにもならないのでpanicする
	err = sqsClient.Poll(ctx, virusScanner.Process)
	if err != nil {
		logger.Error("SQSのポーリングに失敗しました")
		logger.Error(err.Error())
		panic(err)
	}
}

func NewCommand(cfg config.Configuration) cmd.CommandInterface {
	return &CommandPoll{}
}
