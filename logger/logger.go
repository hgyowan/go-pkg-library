package logger

import (
	"fmt"
	"github.com/hgyowan/go-pkg-library/envs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

type logger struct {
	Logger     *zap.Logger
	FileLogger *zap.Logger
	GormLogger *gormLogger
}

var ZapLogger *logger

func MustInitZapLogger() {
	// Define custom encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "",
		MessageKey:     "message",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // Capitalize the log level names
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC timestamp format
		EncodeDuration: zapcore.SecondsDurationEncoder, // Duration in seconds
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Short caller (file and line)
	}

	loglevel := strings.ToUpper(envs.LogLevel)
	var zapLogLevel zapcore.Level
	switch loglevel {
	case "DEBUG":
		zapLogLevel = zapcore.DebugLevel
	case "INFO":
		zapLogLevel = zapcore.InfoLevel
	case "WARN":
		zapLogLevel = zapcore.WarnLevel
	case "ERROR":
		zapLogLevel = zapcore.ErrorLevel
	case "FATAL":
		zapLogLevel = zapcore.FatalLevel
	default:
		zapLogLevel = zapcore.DebugLevel
	}

	stderrLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapLogLevel && (level == zapcore.ErrorLevel || level == zapcore.FatalLevel || level == zapcore.WarnLevel)
	})
	stdoutLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapLogLevel && (level == zapcore.DebugLevel || level == zapcore.InfoLevel)
	})

	stderrSyncer := zapcore.Lock(os.Stderr)
	stdoutSyncer := zapcore.Lock(os.Stdout)

	fileLogger := &lumberjack.Logger{
		Filename: fmt.Sprintf("/var/log/containers/%s-%s-error.log", envs.ServerName, envs.ServiceType), // Or any other path
		MaxSize:  50,                                                                                    // MB; after this size, a new log file is created
		MaxAge:   2,                                                                                     // Days
		Compress: false,                                                                                 // Compress the backups using gzip
	}

	fileSyncer := zapcore.AddSync(fileLogger)

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			stderrSyncer,
			stderrLevel,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			stdoutSyncer,
			stdoutLevel,
		),
	)

	fileCore := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			fileSyncer,
			stderrLevel,
		),
	)

	l := zap.New(core)
	fileL := zap.New(fileCore)
	ZapLogger = &logger{Logger: l, FileLogger: fileL, GormLogger: &gormLogger{
		Logger:                l,
		SkipErrRecordNotFound: true,
	}}
}
