package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestGetConfig GetConfig関数にてQUEUE_URLの環境変数の値がConfiguration構造体に格納されていることを確認するテスト
func TestGetConfig(t *testing.T) {
	expected := "test"
	err := os.Setenv("QUEUE_URL", expected)
	assert.Nil(t, err)
	config, err := GetConfig()

	assert.Nil(t, err)
	assert.Equal(t, expected, config.QueueURL)
	assert.Equal(t, DefaultRegion, config.Region)
	assert.Equal(t, DefaultMaxNumberOfMessages, config.MaxNumberOfMessages)
	assert.Equal(t, DefaultWaitTimeSeconds, config.WaitTimeSeconds)
	assert.Equal(t, DefaultClamdHost, config.ClamdHost)
	assert.Equal(t, DefaultClamdPort, config.ClamdPort)
}

// TestInitialize Initialize関数にてconf変数がnilになることを確認するテスト
func TestInitialize(t *testing.T) {
	Initialize()
	assert.Nil(t, conf)
}
