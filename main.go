/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"github.com/Fufuhu/go-clamav/cmd"
	"github.com/Fufuhu/go-clamav/config"
)

func main() {

	conf, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(conf.QueueURL)

	cmd.Execute()
}
