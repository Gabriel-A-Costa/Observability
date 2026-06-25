package config

import "go.uber.org/zap"

func NewLogger(env string) *zap.Logger {
	logger := zap.Must(zap.NewProduction())

	if env == "development" {
		logger = zap.Must(zap.NewDevelopment())
	}

	return logger
}
