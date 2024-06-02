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

func TestClient_ReceiveMessagesWithNullMessage(t *testing.T) {
	conf := config.Configuration{
		Region:              "ap-northeast-1",
		QueueURL:            "http://localhost:9324/000000000000/queue1",
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     2,
	}
	client := NewClient(conf)
	assert.NotNil(t, client)

	ctx := context.Background()
	receiveMessageInput := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String("http://localhost:9324/000000000000/queue1"),
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     2,
	}

	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(conf.Region))

	// TODO: BaseEndpointをconfigに設定するようにして指定がない場合は、configに設定しないようにする
	cfg.BaseEndpoint = aws.String("http://localhost:9324")
	assert.Nil(t, err)

	svc := awsSqs.NewFromConfig(cfg)

	_, err = svc.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String("http://localhost:9324/000000000000/queue1"),
		MessageBody: aws.String("test"),
	})
	assert.Nil(t, err)

	messages, err := client.ReceiveMessages(ctx, receiveMessageInput, svc)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(messages))
}
