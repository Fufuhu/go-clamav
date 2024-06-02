package sqs

import (
	"context"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	awsSqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	conf := config.Configuration{
		Region:              "ap-northeast-1",
		QueueURL:            "https://sqs.ap-northeast-1.amazonaws.com/123456789012/queue",
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
	}
	client := NewClient(conf)
	assert.NotNil(t, client)
}

func TestClient_ReceiveMessages(t *testing.T) {
	conf := config.Configuration{
		Region:              "ap-northeast-1",
		QueueURL:            "http://localhost:9324/000000000000/queue1",
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     20,
		BaseUrl:             "http://localhost:9324",
	}
	client := NewClient(conf)
	assert.NotNil(t, client)

	ctx := context.Background()
	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(conf.Region))

	// TODO: BaseEndpointをconfigに設定するようにして指定がない場合は、configに設定しないようにする
	if conf.BaseUrl != "" {
		cfg.BaseEndpoint = aws.String(conf.BaseUrl)
	}
	assert.Nil(t, err)

	svc := awsSqs.NewFromConfig(cfg)

	_, err = svc.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl: aws.String("http://localhost:9324/000000000000/queue1"),
		MessageBody: aws.String(`{
  "Records": [
    {
      "eventVersion": "2.1",
      "eventSource": "aws:s3",
      "awsRegion": "ap-northeast-1",
      "eventTime": "2022-08-07T14:33:59.870Z",
      "eventName": "ObjectCreated:Put",
      "userIdentity": {
        "principalId": "AWS:AIDAVMRY2N7OKTN33RYNV"
      },
      "requestParameters": {
        "sourceIPAddress": "60.95.0.122"
      },
      "responseElements": {
        "x-amz-request-id": "Q73VJ1CPJ64CKJQ0",
        "x-amz-id-2": "jqP4VGy4ubSEOvB+XRCdTjWUJEuCkkWRyiRlxdKCNqjP8cTjRUg0JGhDYsW9RprSsQPqdnlOviWD11mpmynwSJzlRyzzT8rgCka5XEnLzq8="
      },
      "s3": {
        "s3SchemaVersion": "1.0",
        "configurationId": "SQS-Event",
        "bucket": {
          "name": "20220807-sqs-test",
          "ownerIdentity": {
            "principalId": "A2B5KBXGR14B9R"
          },
          "arn": "arn:aws:s3:::20220807-sqs-test"
        },
        "object": {
          "key": "hane.jpg",
          "size": 9846,
          "eTag": "ad1cdeed43375dca5b5e892be0968525",
          "sequencer": "0062EFCD57CFFC5419"
        }
      }
    }
  ]
}
`),
	})
	assert.Nil(t, err)

	messages, err := client.ReceiveMessages(ctx, svc)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(messages))
	assert.Equal(t, "hane.jpg", messages[0].Key)
	assert.Equal(t, "20220807-sqs-test", messages[0].Bucket)
}
