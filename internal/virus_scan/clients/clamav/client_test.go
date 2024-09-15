package clamav

import (
	"fmt"
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

func TestClient_GetAddress(t *testing.T) {
	conf, err := config.GetConfig()
	assert.Nil(t, err)
	assert.NotNil(t, conf)

	client := NewClient(*conf)
	assert.NotNil(t, client)

	address := client.GetAddress()
	assert.NotEmpty(t, address)
	assert.Equal(t, "localhost:3310", address)
}

type MockReader struct {
	Count int
}

func (m MockReader) Read(buf []byte) (n int, err error) {
	return 0, nil
}

func TestClient_Scan(t *testing.T) {
	conf, err := config.GetConfig()
	assert.Nil(t, err)
	assert.NotNil(t, conf)

	client := NewClient(*conf)
	assert.NotNil(t, client)

	t.Run("正常ファイル", func(t *testing.T) {

		mockReader := MockReader{Count: 1}

		// テスト用のバイト列を作成
		result, err := client.Scan(mockReader)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		fmt.Println(result)
	})

}
