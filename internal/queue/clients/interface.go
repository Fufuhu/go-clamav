package clients

import "context"

type ClientInterface interface {
	Poll(ctx context.Context, process func(QueueMessageInterface) error) error
	SendMessage(ctx context.Context, message string) (*string, error)
	ReceiveMessages(ctx context.Context) ([]QueueMessageInterface, error)
	DeleteMessage(ctx context.Context, receiptHandle string) error
}
