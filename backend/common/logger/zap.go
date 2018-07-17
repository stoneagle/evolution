package logger

import (
	"evolution/backend/common/config"
	"runtime"
	"time"

	"github.com/davecgh/go-spew/spew"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level int8

const (
	DebugLevel Level = iota + 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Logger struct {
	logger   *zap.SugaredLogger
	Project  string
	Resource string
	Caller   int
}

type detail struct {
	file     string      `json:"file"`
	line     int         `json:"line"`
	project  string      `json:"project"`
	resource string      `json:"resource"`
	action   string      `json:"func"`
	msg      interface{} `json:"msg"`
	err      string      `json:"err"`
}

var (
	onceLogger    *Logger = &Logger{}
	DefaultCaller int     = 2
)

func (l *Logger) Log(level Level, msg interface{}, err error) {
	fpcs := make([]uintptr, 1)
	n := runtime.Callers(l.Caller, fpcs)
	if n == 0 {
		l.logger.Error("there is no logger caller")
		return
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		l.logger.Error("logger caller is nil")
		return
	}

	if len(l.Resource) == 0 {
		l.Resource = "server"
	}

	file, line := caller.FileLine(fpcs[0] - 1)
	message := detail{
		msg:      msg,
		file:     file,
		line:     line,
		project:  l.Project,
		resource: l.Resource,
		action:   caller.Name(),
	}
	if err != nil {
		message.err = err.Error()
	}

	res := spew.Sdump(message)
	switch level {
	case DebugLevel:
		l.logger.Debug(res)
	case InfoLevel:
		l.logger.Info(res)
	case WarnLevel:
		l.logger.Warn(res)
	case ErrorLevel:
		l.logger.Error(res)
	}
}

func Get() *Logger {
	if (Logger{}) == *onceLogger {
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
			DisableCaller:    true,
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    encoder_cfg,
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}

		conf := config.Get()
		if conf.App.Mode == "debug" {
			customCfg.OutputPaths = []string{"stderr"}
		} else {
			logPath := conf.App.Log
			customCfg.OutputPaths = []string{"stderr", logPath + "/server.log"}
		}

		logger, _ := customCfg.Build()
		newLogger := logger.Named("evolution")
		defer newLogger.Sync()
		onceLogger = &Logger{
			Caller: DefaultCaller,
		}
		onceLogger.logger = newLogger.Sugar()
	}
	return onceLogger
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}
