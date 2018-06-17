package common

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var onceLogger *zap.SugaredLogger = &zap.SugaredLogger{}

func GetLogger() *zap.SugaredLogger {
	if (zap.SugaredLogger{}) == *onceLogger {
		encoder_cfg := zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeTime:     TimeEncoder,
		}

		currlevel := zap.NewAtomicLevelAt(zap.DebugLevel)
		customCfg := zap.Config{
			Level:            currlevel,
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    encoder_cfg,
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}

		logger, _ := customCfg.Build()
		newLogger := logger.Named("quant")
		defer newLogger.Sync()
		onceLogger = newLogger.Sugar()
	}

	return onceLogger
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}
