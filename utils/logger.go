package utils

import (
	"nam-club/NumBuyer_back/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

var (
	Log *zap.Logger
)

func init() {
	// Log, _ = zap.NewDevelopment()
	var err error
	if Log, err = setUp(); err != nil {
		panic(err)
	}
	// defer Log.Sync()
}

func setUp() (*zap.Logger, error) {

	debugLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	if config.Env.LogLevel == LogLevelDebug {
		debugLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else if config.Env.LogLevel == LogLevelInfo {
		debugLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else if config.Env.LogLevel == LogLevelWarn {
		debugLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	} else if config.Env.LogLevel == LogLevelError {
		debugLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	} else {
		debugLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	conf := zap.Config{
		Level:             debugLevel,
		Development:       true,
		Encoding:          "console",
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stdout"},
		DisableCaller:     true,
		DisableStacktrace: true,
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			CallerKey:    "caller",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	var err error
	l, err := conf.Build()
	if err != nil {
		return nil, err
	}

	return l, nil
}
