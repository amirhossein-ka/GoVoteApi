package zap

import (
	"GoVoteApi/config"
	"GoVoteApi/pkg/logger"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger_impl struct {
	z *zap.SugaredLogger
}


func New(cfg *config.Log) (logger.Logger, error) {
	writer, err := loggerWriter(cfg.FilePath)
	if err != nil {
		return nil, err
	}
	enc := encoder()

	core := zapcore.NewCore(enc, writer, zapcore.InfoLevel)
	sgl := zap.New(core).Sugar()

	return &logger_impl{z: sgl}, nil
}

func loggerWriter(path string) (zapcore.WriteSyncer, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return zapcore.AddSync(f), nil
}

func encoder() zapcore.Encoder {
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(config)
}
