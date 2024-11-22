package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 定义简化的日志接口
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	WithFields(fields map[string]interface{}) Logger
}

// ZapLogger 实现 Logger 接口
type ZapLogger struct {
	logger *zap.SugaredLogger
}

var (
	defaultLogger Logger
	once          sync.Once
)

// Config 日志配置
type LoggerConfig struct {
	Level      string `yaml:"level"`      // 日志级别
	Filename   string `yaml:"filename"`   // 日志文件名
	MaxSize    int    `yaml:"maxsize"`    // 单个文件最大尺寸，单位 MB
	MaxAge     int    `yaml:"maxage"`     // 保留天数
	MaxBackups int    `yaml:"maxbackups"` // 保留文件个数
	Compress   bool   `yaml:"compress"`   // 是否压缩
	Console    bool   `yaml:"console"`    // 是否同时输出到控制台
}

var initErr error

// Init 初始化全局日志实例
func Init(cfg LoggerConfig) error {
	once.Do(func() {
		defaultLogger, initErr = NewLogger(cfg)
	})
	return initErr
}

// NewLogger 创建新的日志实例
func NewLogger(cfg LoggerConfig) (Logger, error) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if cfg.Level != "" {
		level.UnmarshalText([]byte(cfg.Level))
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core

	// 文件输出
	if cfg.Filename != "" {
		writer := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,    // MB
			MaxAge:     cfg.MaxAge,     // days
			MaxBackups: cfg.MaxBackups, // files
			Compress:   cfg.Compress,   // 是否压缩
			LocalTime:  true,           // 使用本地时间
		}
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(writer),
			level,
		))
	}

	// 控制台输出
	if cfg.Console {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core).Sugar()

	return &ZapLogger{logger: logger}, nil
}

// 实现 Logger 接口
func (l *ZapLogger) Info(msg string, args ...interface{})  { l.logger.Infow(msg, args...) }
func (l *ZapLogger) Error(msg string, args ...interface{}) { l.logger.Errorw(msg, args...) }
func (l *ZapLogger) Debug(msg string, args ...interface{}) { l.logger.Debugw(msg, args...) }
func (l *ZapLogger) Warn(msg string, args ...interface{})  { l.logger.Warnw(msg, args...) }

func (l *ZapLogger) WithFields(fields map[string]interface{}) Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &ZapLogger{logger: l.logger.With(args...)}
}

// 全局方法
func Info(msg string, args ...interface{})            { defaultLogger.Info(msg, args...) }
func Error(msg string, args ...interface{})           { defaultLogger.Error(msg, args...) }
func Debug(msg string, args ...interface{})           { defaultLogger.Debug(msg, args...) }
func Warn(msg string, args ...interface{})            { defaultLogger.Warn(msg, args...) }
func WithFields(fields map[string]interface{}) Logger { return defaultLogger.WithFields(fields) }

// InitFromConfig 从配置初始化日志
func InitFromConfig(cfg LoggerConfig) error {
	return Init(cfg)
}
