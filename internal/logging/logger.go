package logging

import "go.uber.org/zap"

var logger *zap.Logger

// GetLogger GetLogger関数はロガーを取得する
func GetLogger() *zap.Logger {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return logger
}
