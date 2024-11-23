package logger

import (
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 定义日志接口
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	WithFields(fields map[string]interface{}) Logger
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `yaml:"level"`      // 日志级别
	Filename   string `yaml:"filename"`   // 日志文件名
	MaxSize    int    `yaml:"maxsize"`    // 单个文件最大尺寸，MB
	MaxAge     int    `yaml:"maxage"`     // 保留天数
	MaxBackups int    `yaml:"maxbackups"` // 保留文件个数
	Compress   bool   `yaml:"compress"`   // 是否压缩
	Console    bool   `yaml:"console"`    // 是否输出到控制台
}

type zapLogger struct {
	logger *zap.SugaredLogger
}

var (
	defaultLogger Logger
	once          sync.Once
	initErr       error
)

// InitializeLogger 初始化全局日志实例
func InitializeLogger(cfg LoggerConfig) error {
	once.Do(func() {
		defaultLogger, initErr = NewLogger(cfg)
	})
	return initErr
}

// NewLogger 创建新的日志实例
func NewLogger(cfg LoggerConfig) (Logger, error) {
	if err := ensureLogDir(cfg.Filename); err != nil {
		return nil, err
	}

	cores := buildZapCores(cfg)
	logger := zap.New(zapcore.NewTee(cores...)).Sugar()

	return &zapLogger{logger: logger}, nil
}

// 实现 Logger 接口
func (l *zapLogger) Info(msg string, args ...interface{})  { l.logger.Infow(msg, args...) }
func (l *zapLogger) Error(msg string, args ...interface{}) { l.logger.Errorw(msg, args...) }
func (l *zapLogger) Debug(msg string, args ...interface{}) { l.logger.Debugw(msg, args...) }
func (l *zapLogger) Warn(msg string, args ...interface{})  { l.logger.Warnw(msg, args...) }

func (l *zapLogger) WithFields(fields map[string]interface{}) Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &zapLogger{logger: l.logger.With(args...)}
}

// 全局方法
func Info(msg string, args ...interface{})            { defaultLogger.Info(msg, args...) }
func Error(msg string, args ...interface{})           { defaultLogger.Error(msg, args...) }
func Debug(msg string, args ...interface{})           { defaultLogger.Debug(msg, args...) }
func Warn(msg string, args ...interface{})            { defaultLogger.Warn(msg, args...) }
func WithFields(fields map[string]interface{}) Logger { return defaultLogger.WithFields(fields) }

// 内部辅助函数
func ensureLogDir(filename string) error {
	if filename == "" {
		return nil
	}
	dir := filepath.Dir(filename)
	return os.MkdirAll(dir, 0755)
}

func buildZapCores(cfg LoggerConfig) []zapcore.Core {
	var cores []zapcore.Core
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	level := getLogLevel(cfg.Level)

	if cfg.Filename != "" {
		cores = append(cores, createFileCore(cfg, encoderConfig, level))
	}

	if cfg.Console {
		cores = append(cores, createConsoleCore(encoderConfig, level))
	}

	return cores
}

func getLogLevel(levelStr string) zapcore.Level {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if levelStr != "" {
		level.UnmarshalText([]byte(levelStr))
	}
	return level.Level()
}

func createFileCore(cfg LoggerConfig, encoderConfig zapcore.EncoderConfig, level zapcore.Level) zapcore.Core {
	writer := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		level,
	)
}

func createConsoleCore(encoderConfig zapcore.EncoderConfig, level zapcore.Level) zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)
}
