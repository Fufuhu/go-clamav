package config

type Configuration struct {
	QueueURL string `envconfig:"QUEUE_URL" required:"true"`
}
