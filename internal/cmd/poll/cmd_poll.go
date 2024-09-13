package poll

import (
	"fmt"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/cmd"
	"github.com/spf13/cobra"
)

type CommandPoll struct {
}

func (p *CommandPoll) Run(cmd *cobra.Command, args []string) {
	fmt.Println("poll called")

	// SQSからメッセージを取得する

	// SQSから取得したメッセージをパースする

	// S3からファイルを取得する

	// ファイルをスキャンする

	// スキャン結果をDynamoDBに保存する

	// SQSからメッセージを削除する
}

func NewCommand(cfg config.Configuration) cmd.CommandInterface {
	return &CommandPoll{}
}
