package clients

import (
	"context"
	"fmt"
)

type QueueMessage struct {
	Bucket        string
	Key           string
	ReceiptHandle string
}

func (m *QueueMessage) GetObjectPath() string {
	path := fmt.Sprintf("s3://%s/%s", m.Bucket, m.Key)
	return path
}

func (m *QueueMessage) DeleteMessage(ctx context.Context, client ClientInterface) error {

	return nil
}

func (m *QueueMessage) SetReceiptHandle(receiptHandle string) {
	m.ReceiptHandle = receiptHandle
}

func (m *QueueMessage) GetReceiptHandle() string {
	return m.ReceiptHandle
}

func (m *QueueMessage) GetBucket() string {
	return m.Bucket
}

func (m *QueueMessage) SetBucket(bucket string) {
	m.Bucket = bucket
}

func (m *QueueMessage) GetKey() string {
	return m.Key
}

func (m *QueueMessage) SetKey(key string) {
	m.Key = key
}
