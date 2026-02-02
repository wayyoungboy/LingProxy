package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/lingproxy/lingproxy/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 日志接口
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	WithFields(fields ...Field) Logger
	SetLevel(level LogLevel)
	GetLevel() LogLevel
}

// Field 日志字段
type Field struct {
	Key   string
	Value interface{}
}

// F 创建日志字段
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
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

// ParseLogLevel 解析日志级别字符串
func ParseLogLevel(level string) LogLevel {
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

// LogEntry 日志条目
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// defaultLogger 默认日志实现
type defaultLogger struct {
	output io.Writer
	level  LogLevel
	format string // "json" or "text"
	fields map[string]interface{}
	mu     sync.RWMutex
}

var (
	// 全局日志实例
	globalLogger Logger
)

// multiWriter 多写入器，支持同时写入多个目标
type multiWriter struct {
	writers []io.Writer
}

func (mw *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range mw.writers {
		n, err = w.Write(p)
		if err != nil {
			return n, err
		}
		if n != len(p) {
			return n, io.ErrShortWrite
		}
	}
	return len(p), nil
}

// Init 初始化日志
func Init(cfg config.LogConfig) {
	var writers []io.Writer

	// 设置输出
	switch cfg.Output {
	case "stdout":
		writers = []io.Writer{os.Stdout}
	case "file":
		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0755); err != nil {
			fmt.Printf("创建日志目录失败: %v\n", err)
			writers = []io.Writer{os.Stdout}
		} else {
			writers = []io.Writer{&lumberjack.Logger{
				Filename:   cfg.FilePath,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
			}}
		}
	case "both":
		// 同时输出到 stdout 和文件
		writers = []io.Writer{os.Stdout}
		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0755); err != nil {
			fmt.Printf("创建日志目录失败: %v，仅输出到stdout\n", err)
		} else {
			writers = append(writers, &lumberjack.Logger{
				Filename:   cfg.FilePath,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
			})
		}
	default:
		// 默认同时输出到 stdout 和文件（持久化）
		writers = []io.Writer{os.Stdout}
		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0755); err != nil {
			fmt.Printf("创建日志目录失败: %v，仅输出到stdout\n", err)
		} else {
			writers = append(writers, &lumberjack.Logger{
				Filename:   cfg.FilePath,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
			})
		}
	}

	// 创建多写入器
	var output io.Writer
	if len(writers) == 1 {
		output = writers[0]
	} else {
		output = &multiWriter{writers: writers}
	}

	// 设置日志级别
	level := ParseLogLevel(cfg.Level)

	// 设置格式
	format := cfg.Format
	if format != "json" && format != "text" {
		format = "text"
	}

	globalLogger = &defaultLogger{
		output: output,
		level:  level,
		format: format,
		fields: make(map[string]interface{}),
	}
}

// WithFields 创建带字段的日志实例
func (l *defaultLogger) WithFields(fields ...Field) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	for _, field := range fields {
		newFields[field.Key] = field.Value
	}

	return &defaultLogger{
		output: l.output,
		level:  l.level,
		format: l.format,
		fields: newFields,
	}
}

// SetLevel 设置日志级别
func (l *defaultLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel 获取日志级别
func (l *defaultLogger) GetLevel() LogLevel {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

// Debug 输出调试日志
func (l *defaultLogger) Debug(msg string, fields ...Field) {
	l.log(LevelDebug, msg, fields...)
}

// Info 输出信息日志
func (l *defaultLogger) Info(msg string, fields ...Field) {
	l.log(LevelInfo, msg, fields...)
}

// Warn 输出警告日志
func (l *defaultLogger) Warn(msg string, fields ...Field) {
	l.log(LevelWarn, msg, fields...)
}

// Error 输出错误日志
func (l *defaultLogger) Error(msg string, fields ...Field) {
	l.log(LevelError, msg, fields...)
}

// Fatal 输出致命日志
func (l *defaultLogger) Fatal(msg string, fields ...Field) {
	l.log(LevelFatal, msg, fields...)
	os.Exit(1)
}

// log 输出日志
func (l *defaultLogger) log(level LogLevel, msg string, fields ...Field) {
	l.mu.RLock()
	currentLevel := l.level
	l.mu.RUnlock()

	if level < currentLevel {
		return
	}

	// 合并字段
	allFields := make(map[string]interface{})
	l.mu.RLock()
	for k, v := range l.fields {
		allFields[k] = v
	}
	l.mu.RUnlock()

	for _, field := range fields {
		allFields[field.Key] = field.Value
	}

	// 格式化日志
	var logLine string
	if l.format == "json" {
		entry := LogEntry{
			Timestamp: time.Now().Format(time.RFC3339),
			Level:     level.String(),
			Message:   msg,
			Fields:    allFields,
		}
		data, err := json.Marshal(entry)
		if err != nil {
			logLine = fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s","error":"marshal failed: %v"}`,
				time.Now().Format(time.RFC3339), level.String(), msg, err)
		} else {
			logLine = string(data)
		}
	} else {
		// 文本格式
		timestamp := time.Now().Format("2006-01-02 15:04:05.000")
		logLine = fmt.Sprintf("[%s] %s %s", timestamp, level.String(), msg)
		if len(allFields) > 0 {
			for k, v := range allFields {
				logLine += fmt.Sprintf(" %s=%v", k, v)
			}
		}
	}

	fmt.Fprintln(l.output, logLine)
}

// Debug 输出调试日志（全局函数）
func Debug(msg string, fields ...Field) {
	if globalLogger != nil {
		globalLogger.Debug(msg, fields...)
	}
}

// Info 输出信息日志（全局函数）
func Info(msg string, fields ...Field) {
	if globalLogger != nil {
		globalLogger.Info(msg, fields...)
	}
}

// Warn 输出警告日志（全局函数）
func Warn(msg string, fields ...Field) {
	if globalLogger != nil {
		globalLogger.Warn(msg, fields...)
	}
}

// Error 输出错误日志（全局函数）
func Error(msg string, fields ...Field) {
	if globalLogger != nil {
		globalLogger.Error(msg, fields...)
	}
}

// Fatal 输出致命日志（全局函数）
func Fatal(msg string, fields ...Field) {
	if globalLogger != nil {
		globalLogger.Fatal(msg, fields...)
	}
}

// WithFields 创建带字段的日志实例（全局函数）
func WithFields(fields ...Field) Logger {
	if globalLogger != nil {
		return globalLogger.WithFields(fields...)
	}
	return nil
}

// SetLevel 设置日志级别（全局函数）
func SetLevel(level LogLevel) {
	if globalLogger != nil {
		globalLogger.SetLevel(level)
	}
}

// GetLevel 获取日志级别（全局函数）
func GetLevel() LogLevel {
	if globalLogger != nil {
		return globalLogger.GetLevel()
	}
	return LevelInfo
}
