package logger

import (
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var logger *zap.SugaredLogger

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	logger.Panicw(msg, keysAndValues...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	logger.Fatalw(msg, keysAndValues...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

func InitLog(style string, path string, level string) error {
	encoder := getZapEncoder(style)
	writeSyncer, err := getZapWriterSync(path)
	if err != nil {
		return err
	}
	zapLevel := getZapLevel(level)
	core := zapcore.NewCore(encoder, writeSyncer, zapLevel)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	return nil
}

func Sync() {
	logger.Sync()
}

// 设定Zap编码格式
func getZapEncoder(style string) zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "file",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	switch style {
	case "JSON", "Json", "json":
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
		return zapcore.NewJSONEncoder(cfg)
	case "CONSOLE", "Console", "console":
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(cfg)
	default:
		return zapcore.NewJSONEncoder(cfg)
	}
}

// 设定日志输出按天分割
func getZapWriterSync(path string) (zapcore.WriteSyncer, error) {
	hook, err := rotateLogs.New(
		path+".%Y%m%d",
		rotateLogs.WithLinkName(path),
		rotateLogs.WithMaxAge(time.Hour*24*7),
		rotateLogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(hook), nil
}

// 设定日志级别
func getZapLevel(level string) zapcore.Level {
	switch level {
	case "DEBUG", "Debug", "debug":
		return zapcore.DebugLevel
	case "INFO", "Info", "info":
		return zapcore.InfoLevel
	case "WARN", "Warn", "warn":
		return zapcore.WarnLevel
	case "ERROR", "Error", "error":
		return zapcore.ErrorLevel
	case "PANIC", "Panic", "panic":
		return zapcore.PanicLevel
	case "FATAL", "Fatal", "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
