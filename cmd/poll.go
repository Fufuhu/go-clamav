/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/cmd/poll"
	"github.com/Fufuhu/go-clamav/internal/logging"
	"github.com/spf13/cobra"
)

// pollCmd represents the poll command
var pollCmd = &cobra.Command{
	Use:   "poll",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.GetLogger()
		defer logger.Sync()
		cfg, err := config.GetConfig()
		if err != nil {
			logger.Error("設定ファイルの読み込みに失敗しました")
			logger.Error(err.Error())
			panic(err)
		}
		command := poll.NewCommand(*cfg)
		command.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(pollCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pollCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pollCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
