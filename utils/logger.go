package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log *zap.Logger
)

func init() {
	// Log, _ = zap.NewDevelopment()
	Log, _ = setUp()
	// defer Log.Sync()
}

func setUp() (*zap.Logger, error) {

	conf := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
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
