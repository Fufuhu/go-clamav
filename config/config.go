package config

import "github.com/kelseyhightower/envconfig"

type Configuration struct {
	QueueURL string `envconfig:"QUEUE_URL" required:"true"`
}

var conf *Configuration

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
