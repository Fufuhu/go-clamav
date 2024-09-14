package clients

import "context"

type QueueMessageInterface interface {
	DeleteMessage(ctx context.Context, client ClientInterface) error
	SetReceiptHandle(receiptHandle string)
	GetReceiptHandle() string
	GetBucket() string
	SetBucket(bucket string)
	GetKey() string
	SetKey(key string)
}
