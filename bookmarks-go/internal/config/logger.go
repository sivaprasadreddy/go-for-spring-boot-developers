package config

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(cfg AppConfig) *Logger {
	logFile := cfg.Logging.FileName
	logLevel, err := zap.ParseAtomicLevel(cfg.Logging.Level)
	if err != nil {
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	hook := lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    1024,
		MaxBackups: 30,
		MaxAge:     7,
		Compress:   true,
	}

	encoder := getEncoder()
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		logLevel)
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	}
	if cfg.Environment != "prod" {
		options = append(options, zap.Development())
	}
	sugaredLogger := zap.New(core, options...).With(zap.String("env", cfg.Environment)).Sugar()
	return &Logger{sugaredLogger}
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}
