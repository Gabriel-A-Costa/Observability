package config

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(env string) *zap.Logger {

	// Condiguração do lumberjack para definir local, tamanho e numero de backups dos arquivos de logs
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "/app/logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
	})

	// Configuração personalizada do log (JSON)
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// Core que escreve em arquivo
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		fileWriter,
		zapcore.InfoLevel,
	)

	// Core que escreve no terminal
	var consoleCore zapcore.Core
	if env == "development" {
		consoleCore = zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel,
		)
	} else {
		consoleCore = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		)
	}

	return zap.New(zapcore.NewTee(fileCore, consoleCore), zap.AddCaller())
}
