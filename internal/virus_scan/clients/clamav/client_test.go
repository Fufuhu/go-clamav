package clamav

import (
	"github.com/Fufuhu/go-clamav/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	conf, err := config.GetConfig()
	assert.Nil(t, err)
	assert.NotNil(t, conf)

	client := NewClient(*conf)
	assert.NotNil(t, client)
}
