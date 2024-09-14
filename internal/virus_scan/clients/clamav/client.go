package clamav

import "github.com/Fufuhu/go-clamav/config"

type Client struct {
	conf config.Configuration
}

func NewClient(conf config.Configuration) *Client {
	return &Client{
		conf: conf,
	}
}
