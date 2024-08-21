package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func Init(appName, appVersion, level, format string) error {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg.Encoding = format
	switch format {
	case "console":
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.EncoderConfig.EncodeLevel = nil
	}

	var lvl zapcore.Level
	lvl.Set(level)
	cfg.Level = zap.NewAtomicLevelAt(lvl)

	l, err := cfg.Build()
	if err != nil {
		return err
	}

	if format == "console" {
		logger = l.Sugar()
		return nil
	}

	logger = l.Sugar().With(
		"app_name", appName,
		"app_version", appVersion,
	)

	return nil
}

func G() *zap.SugaredLogger {
	if logger == nil {
		return zap.S()
	}
	return logger
}
