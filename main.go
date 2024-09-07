/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"github.com/Fufuhu/go-clamav/cmd"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/logging"
)

func main() {

	conf, err := config.GetConfig()
	logger := logging.GetLogger()
	if err != nil {
		logger.Error("設定ファイルの読み込みに失敗しました")
		logger.Error(err.Error())
		panic(err)
	}

	fmt.Println(conf.QueueURL)

	cmd.Execute()
}
