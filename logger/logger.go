package logger

import (
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	fileName       = "metr.log"
	rotationCount  = 10
	rotateSuffix   = ".%Y%m%d"
	rotateDuration = time.Duration(24 * time.Hour)
)

// NewLogger returns Logger
func NewLogger(dir string) (*zap.Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	cores := []zapcore.Core{}
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	path, err := filepath.Abs(filepath.Join(dir, fileName))
	if err != nil {
		return nil, err
	}

	options := []rotatelogs.Option{
		rotatelogs.WithClock(rotatelogs.Local),
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(rotationCount),
		rotatelogs.WithLinkName(path),
		rotatelogs.WithRotationTime(rotateDuration),
	}

	w, err := rotatelogs.New(
		path+rotateSuffix,
		options...,
	)
	if err != nil {
		return nil, err
	}

	logCore := zapcore.NewCore(
		encoder,
		zapcore.AddSync(w),
		zapcore.InfoLevel,
	)
	cores = append(cores, logCore)
	logger := zap.New(zapcore.NewTee(cores...))

	return logger, nil
}

// NewNoLogger return logger does not output
func NewNoLogger() (*zap.Logger, error) {
	logger := zap.New(zapcore.NewTee([]zapcore.Core{}...))

	return logger, nil
}
