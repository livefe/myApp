package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"myApp/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局日志实例
var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

// 日志级别
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

// InitLogger 初始化日志，直接从config包获取配置
func InitLogger() {
	// 从全局配置中获取日志配置
	logConfig := config.Conf.Logger

	// 获取配置参数
	level := logConfig.Level
	filePath := logConfig.FilePath
	maxSize := logConfig.MaxSize
	maxBackups := logConfig.MaxBackups
	maxAge := logConfig.MaxAge
	compress := logConfig.Compress
	console := logConfig.Console
	// 设置默认值
	if level == "" {
		level = InfoLevel
	}
	if filePath == "" {
		filePath = "./logs/app.log"
	}
	if maxSize == 0 {
		maxSize = 100
	}
	if maxBackups == 0 {
		maxBackups = 10
	}
	if maxAge == 0 {
		maxAge = 30
	}

	// 创建日志目录
	logDir := filepath.Dir(filePath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			fmt.Printf("创建日志目录失败: %v\n", err)
		}
	}

	// 设置日志级别
	var zapLevel zapcore.Level
	switch level {
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	case FatalLevel:
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置日志输出
	var cores []zapcore.Core

	// 文件输出
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	})
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		fileWriter,
		zapLevel,
	)
	cores = append(cores, fileCore)

	// 控制台输出
	if console {
		consoleEncoderConfig := encoderConfig
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(consoleEncoderConfig),
			zapcore.AddSync(os.Stdout),
			zapLevel,
		)
		cores = append(cores, consoleCore)
	}

	// 创建日志实例
	core := zapcore.NewTee(cores...)
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Sugar = Logger.Sugar()
}

// 自定义时间格式编码器
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// Debug 输出Debug级别日志
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Info 输出Info级别日志
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Warn 输出Warn级别日志
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error 输出Error级别日志
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Fatal 输出Fatal级别日志
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Debugf 输出Debug级别格式化日志
func Debugf(format string, args ...interface{}) {
	Sugar.Debugf(format, args...)
}

// Infof 输出Info级别格式化日志
func Infof(format string, args ...interface{}) {
	Sugar.Infof(format, args...)
}

// Warnf 输出Warn级别格式化日志
func Warnf(format string, args ...interface{}) {
	Sugar.Warnf(format, args...)
}

// Errorf 输出Error级别格式化日志
func Errorf(format string, args ...interface{}) {
	Sugar.Errorf(format, args...)
}

// Fatalf 输出Fatal级别格式化日志
func Fatalf(format string, args ...interface{}) {
	Sugar.Fatalf(format, args...)
}

// WithField 添加单个字段
func WithField(key string, value interface{}) *zap.Logger {
	return Logger.With(zap.Any(key, value))
}

// WithFields 添加多个字段
func WithFields(fields map[string]interface{}) *zap.Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return Logger.With(zapFields...)
}

// WithError 添加错误信息
func WithError(err error) *zap.Logger {
	return Logger.With(zap.Error(err))
}
