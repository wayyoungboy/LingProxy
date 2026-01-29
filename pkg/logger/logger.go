package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/lingproxy/lingproxy/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 日志接口
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

// defaultLogger 默认日志实现
type defaultLogger struct {
	output io.Writer
	level  LogLevel
}

// LogLevel 日志级别
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

var (
	// 全局日志实例
	globalLogger Logger
)

// Init 初始化日志
func Init(cfg config.LogConfig) {
	var output io.Writer

	// 设置输出
	switch cfg.Output {
	case "stdout":
		output = os.Stdout
	case "file":
		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0755); err != nil {
			fmt.Printf("创建日志目录失败: %v\n", err)
			output = os.Stdout
		} else {
			output = &lumberjack.Logger{
				Filename:   cfg.FilePath,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
			}
		}
	default:
		output = os.Stdout
	}

	// 设置日志级别
	level := parseLogLevel(cfg.Level)

	globalLogger = &defaultLogger{
		output: output,
		level:  level,
	}
}

// parseLogLevel 解析日志级别
func parseLogLevel(level string) LogLevel {
	switch level {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	default:
		return LevelInfo
	}
}

// Debug 输出调试日志
func (l *defaultLogger) Debug(msg string, fields ...interface{}) {
	l.log(LevelDebug, msg, fields...)
}

// Info 输出信息日志
func (l *defaultLogger) Info(msg string, fields ...interface{}) {
	l.log(LevelInfo, msg, fields...)
}

// Warn 输出警告日志
func (l *defaultLogger) Warn(msg string, fields ...interface{}) {
	l.log(LevelWarn, msg, fields...)
}

// Error 输出错误日志
func (l *defaultLogger) Error(msg string, fields ...interface{}) {
	l.log(LevelError, msg, fields...)
}

// Fatal 输出致命日志
func (l *defaultLogger) Fatal(msg string, fields ...interface{}) {
	l.log(LevelFatal, msg, fields...)
	os.Exit(1)
}

// log 输出日志
func (l *defaultLogger) log(level LogLevel, msg string, fields ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("[%s] %s %s", timestamp, level.String(), msg)

	if len(fields) > 0 {
		logMsg += " "
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				logMsg += fmt.Sprintf("%s=%v ", fields[i], fields[i+1])
			}
		}
	}

	fmt.Fprintln(l.output, logMsg)
}

// Debug 输出调试日志
func Debug(msg string, fields ...interface{}) {
	if globalLogger != nil {
		globalLogger.Debug(msg, fields...)
	}
}

// Info 输出信息日志
func Info(msg string, fields ...interface{}) {
	if globalLogger != nil {
		globalLogger.Info(msg, fields...)
	}
}

// Warn 输出警告日志
func Warn(msg string, fields ...interface{}) {
	if globalLogger != nil {
		globalLogger.Warn(msg, fields...)
	}
}

// Error 输出错误日志
func Error(msg string, fields ...interface{}) {
	if globalLogger != nil {
		globalLogger.Error(msg, fields...)
	}
}

// Fatal 输出致命日志
func Fatal(msg string, fields ...interface{}) {
	if globalLogger != nil {
		globalLogger.Fatal(msg, fields...)
	}
}