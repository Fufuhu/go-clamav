package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	expected := "test"
	err := os.Setenv("QUEUE_URL", expected)
	assert.Nil(t, err)
	config, err := GetConfig()

	assert.Nil(t, err)
	assert.Equal(t, expected, config.QueueURL)
}
