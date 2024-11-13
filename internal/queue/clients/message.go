package clients

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"github.com/Fufuhu/go-clamav/config"
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

func (m *QueueMessage) GetObject() string {
	return m.GetObjectPath()
}

func (m *QueueMessage) IsTargetFile(conf config.Configuration) (bool, error) {
	// ファイル名のパターンが指定されていない場合は全てのファイルを対象とする
	if conf.ScanningTargetFilePatterns == "" {
		return true, nil
	}

	patterns := strings.Split(conf.ScanningTargetFilePatterns, ",")

	// 指定されている場合は、いずれかのパターンにマッチするファイルのみを対象とする
	for _, pattern := range patterns {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return false, fmt.Errorf("正規表現のコンパイルに失敗しました: %v", err)
		}
		if regex.MatchString(m.Key) {
			return true, nil
		}
	}
	return false, nil
}