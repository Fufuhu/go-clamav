package clamav

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Fufuhu/go-clamav/config"
	"github.com/Fufuhu/go-clamav/internal/logging"
	"io"
	"net"
)

type Client struct {
	conf config.Configuration
}

const InstreamCommand = "nINSTREAM\n"

// GetAddress GetAddress関数はClamdのアドレスを取得する
func (c *Client) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.conf.ClamdHost, c.conf.ClamdPort)
}

// Scan Scan関数はio.Readerで取得されるバイト列をスキャンする
func (c *Client) Scan(reader io.Reader) (Result, error) {
	logger := logging.GetLogger()
	defer logger.Sync()

	// サーバーに接続
	conn, err := net.Dial("tcp", c.GetAddress())
	if err != nil {
		logger.Error("Clamdクライアントの作成に失敗しました")
		return Result{}, err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Error("Clamdクライアントのクローズに失敗しました")
			logger.Error(err.Error())
		}
	}(conn)

	// INSTREAMコマンドの送信
	_, err = conn.Write([]byte(InstreamCommand))
	if err != nil {
		logger.Error("INSTREAMコマンドの送信に失敗しました")
		logger.Error(err.Error())
		return Result{}, err
	}

	// バイト列を分割して送信するためのチャンクバッファを作成する
	buf := make([]byte, 1024)
	// バイト列を読み込んでclamdに送信する
	for {
		n, err := reader.Read(buf)
		if n == 0 {
			break
		}
		if n > 0 {
			// チャンクサイズの送信
			size := uint32(n)
			sizeBuf := new(bytes.Buffer)
			if err := binary.Write(sizeBuf, binary.BigEndian, size); err != nil {
				logger.Error("チャンクサイズのバイト列の作成に失敗しました")
				logger.Error(err.Error())
				return Result{}, err
			}

			// チャンクデータの送信
			_, err = conn.Write(buf[:n])
			if err != nil {
				logger.Error("チャンクデータの送信に失敗しました")
				logger.Error(err.Error())
				return Result{}, err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error("バイト列の読み込みに失敗しました")
			logger.Error(err.Error())
			return Result{}, err
		}
	}

	// データ終了を示すための0バイトチャンクを送信する
	_, err = conn.Write([]byte{0, 0, 0, 0})
	if err != nil {
		logger.Error("データ終了のチャンクの送信に失敗しました")
		logger.Error(err.Error())
		return Result{}, err
	}

	// レスポンスの読み取り
	responseBuf := make([]byte, 4096)
	n, err := conn.Read(responseBuf)
	if err != nil && err != io.EOF {
		logger.Error("レスポンスの読み取りに失敗しました")
		logger.Error(err.Error())
		return Result{}, err
	}
	response := string(responseBuf[:n])

	return Result{
		Message: response,
	}, nil
}

func NewClient(conf config.Configuration) *Client {
	return &Client{
		conf: conf,
	}
}
