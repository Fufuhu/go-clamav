package config

import "github.com/kelseyhightower/envconfig"

type Configuration struct {
	QueueURL            string `envconfig:"QUEUE_URL" required:"true"`
	Region              string `envconfig:"REGION" required:"true" default:"ap-northeast-1"`
	MaxNumberOfMessages int32  `envconfig:"MAX_NUMBER_OF_MESSAGES" required:"true" default:"1"`
	WaitTimeSeconds     int32  `envconfig:"WAIT_TIME_SECONDS" required:"true" default:"20"`
	BaseUrl             string `envconfig:"BASE_URL" required:"false" default:""`
	S3BaseUrl           string `envconfig:"S3_BASE_URL" required:"false" default:""`
}

var conf *Configuration

const DefaultRegion = "ap-northeast-1"
const DefaultMaxNumberOfMessages = int32(1)
const DefaultWaitTimeSeconds = int32(20)

// Initialize Initialize関数はconf変数を初期化する
// 環境変数を変更したあとにプロセスを再起動することなしに設定を再取得したい場合などに利用する
func Initialize() {
	conf = nil
}

// GetConfig GetConfig関数は環境変数から設定を取得しConfiguration構造体に格納する
// Initialize実行後は環境変数の取得からやり直す
func GetConfig() (*Configuration, error) {
	if conf != nil {
		return conf, nil
	}
	conf = &Configuration{}
	err := envconfig.Process("", conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
