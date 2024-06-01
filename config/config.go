package config

import "github.com/kelseyhightower/envconfig"

type Configuration struct {
	QueueURL string `envconfig:"QUEUE_URL" required:"true"`
}

func GetConfig() (*Configuration, error) {
	conf := &Configuration{}
	err := envconfig.Process("", conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
