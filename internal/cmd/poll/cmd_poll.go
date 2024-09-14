package poll

import (
	"context"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/cmd"
	"github.com/Fufuhu/go-clamav/internal/logging"
	"github.com/Fufuhu/go-clamav/internal/queue/clients"
	"github.com/Fufuhu/go-clamav/internal/queue/clients/sqs"
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
	// エラーは拾ったところでどうにもならないのでpanicする
	err = sqsClient.Poll(ctx, func(object clients.QueueMessageInterface) error {
		// TODO 取得したS3オブジェクトを処理するための関数を準備する
		return nil
	})
	if err != nil {
		logger.Error("SQSのポーリングに失敗しました")
		logger.Error(err.Error())
		panic(err)
	}
}

func NewCommand(cfg config.Configuration) cmd.CommandInterface {
	return &CommandPoll{}
}
