package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/config"
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
)

// LogHandler 日志处理器
type LogHandler struct {
	logDir string
}

// NewLogHandler 创建新的日志处理器
func NewLogHandler() *LogHandler {
	cfg := config.C
	logDir := "./logs"
	// 支持 file 和 both 模式
	if (cfg.Log.Output == "file" || cfg.Log.Output == "both") && cfg.Log.FilePath != "" {
		logDir = filepath.Dir(cfg.Log.FilePath)
	}
	return &LogHandler{
		logDir: logDir,
	}
}

// ListLogFiles 获取日志文件列表
// @Summary List log files
// @Description Get a list of all log files
// @Tags logs
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of log files"
// @Router /api/v1/logs/files [get]
func (h *LogHandler) ListLogFiles(c *gin.Context) {
	files, err := h.getLogFiles()
	if err != nil {
		logger.Error("获取日志文件列表失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": files})
}

// GetLogs 获取日志内容
// @Summary Get log content
// @Description Get log content with filtering options
// @Tags logs
// @Accept json
// @Produce json
// @Param file query string false "Log file name"
// @Param level query string false "Filter by log level"
// @Param keyword query string false "Search keyword"
// @Param lines query int false "Number of lines to retrieve (default: 100)"
// @Param tail query bool false "Get tail of file (default: false)"
// @Success 200 {object} map[string]interface{} "Log content"
// @Router /api/v1/logs [get]
func (h *LogHandler) GetLogs(c *gin.Context) {
	fileName := c.Query("file")
	level := c.Query("level")
	keyword := c.Query("keyword")
	linesStr := c.DefaultQuery("lines", "100")
	tailStr := c.DefaultQuery("tail", "false")

	var lines int
	if _, err := fmt.Sscanf(linesStr, "%d", &lines); err != nil || lines <= 0 {
		lines = 100
	}
	if lines > 1000 {
		lines = 1000 // 限制最大行数
	}

	tail := tailStr == "true"

	// 如果没有指定文件，使用默认日志文件
	if fileName == "" {
		cfg := config.C
		if (cfg.Log.Output == "file" || cfg.Log.Output == "both") && cfg.Log.FilePath != "" {
			fileName = filepath.Base(cfg.Log.FilePath)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "日志输出到stdout，无法查看文件内容"})
			return
		}
	}

	filePath := filepath.Join(h.logDir, fileName)
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(h.logDir)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件名"})
		return
	}

	logEntries, err := h.readLogFile(filePath, lines, tail, level, keyword)
	if err != nil {
		logger.Error("读取日志文件失败", logger.F("error", err.Error()), logger.F("file", fileName))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logEntries,
		"total": len(logEntries),
		"file":  fileName,
	})
}

// DownloadLogFile 下载日志文件
// @Summary Download log file
// @Description Download a specific log file
// @Tags logs
// @Accept json
// @Produce application/octet-stream
// @Param file path string true "Log file name"
// @Success 200 {file} file "Log file"
// @Failure 404 {object} map[string]string "File not found"
// @Router /api/v1/logs/files/{file}/download [get]
func (h *LogHandler) DownloadLogFile(c *gin.Context) {
	fileName := c.Param("file")
	filePath := filepath.Join(h.logDir, fileName)

	// 安全检查：确保文件在日志目录内
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(h.logDir)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件名"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.File(filePath)
}

// ClearLogs 清空日志文件
// @Summary Clear log file
// @Description Clear content of a log file
// @Tags logs
// @Accept json
// @Produce json
// @Param file query string true "Log file name"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Bad request"
// @Router /api/v1/logs/clear [post]
func (h *LogHandler) ClearLogs(c *gin.Context) {
	fileName := c.Query("file")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件名不能为空"})
		return
	}

	filePath := filepath.Join(h.logDir, fileName)

	// 安全检查
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(h.logDir)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件名"})
		return
	}

	// 清空文件
	if err := os.Truncate(filePath, 0); err != nil {
		logger.Error("清空日志文件失败", logger.F("error", err.Error()), logger.F("file", fileName))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("日志文件已清空", logger.F("file", fileName))
	c.JSON(http.StatusOK, gin.H{"message": "日志文件已清空"})
}

// getLogFiles 获取日志文件列表
func (h *LogHandler) getLogFiles() ([]map[string]interface{}, error) {
	files := []map[string]interface{}{}

	entries, err := os.ReadDir(h.logDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if !strings.HasSuffix(fileName, ".log") && !strings.HasSuffix(fileName, ".gz") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, map[string]interface{}{
			"name":     fileName,
			"size":     info.Size(),
			"mod_time": info.ModTime().Format(time.RFC3339),
		})
	}

	// 按修改时间排序（最新的在前）
	sort.Slice(files, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, files[i]["mod_time"].(string))
		timeJ, _ := time.Parse(time.RFC3339, files[j]["mod_time"].(string))
		return timeI.After(timeJ)
	})

	return files, nil
}

// LogEntry 日志条目（用于返回）
type LogEntry struct {
	Line      int    `json:"line"`
	Timestamp string `json:"timestamp,omitempty"`
	Level     string `json:"level,omitempty"`
	Message   string `json:"message"`
	Raw       string `json:"raw"`
}

// readLogFile 读取日志文件
func (h *LogHandler) readLogFile(filePath string, lines int, tail bool, levelFilter, keywordFilter string) ([]LogEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var allLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 如果要从尾部读取，反转顺序
	if tail && len(allLines) > lines {
		start := len(allLines) - lines
		allLines = allLines[start:]
	} else if len(allLines) > lines {
		allLines = allLines[:lines]
	}

	// 解析和过滤日志
	entries := []LogEntry{}
	for i, line := range allLines {
		entry := h.parseLogLine(line, i+1)

		// 级别过滤
		if levelFilter != "" && entry.Level != "" && !strings.EqualFold(entry.Level, levelFilter) {
			continue
		}

		// 关键词过滤
		if keywordFilter != "" && !strings.Contains(strings.ToLower(entry.Raw), strings.ToLower(keywordFilter)) {
			continue
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// parseLogLine 解析日志行
func (h *LogHandler) parseLogLine(line string, lineNum int) LogEntry {
	entry := LogEntry{
		Line: lineNum,
		Raw:  line,
	}

	_ = h // 避免未使用接收器的警告

	// 尝试解析JSON格式
	if strings.HasPrefix(line, "{") {
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err == nil {
			if ts, ok := logEntry["timestamp"].(string); ok {
				entry.Timestamp = ts
			}
			if lvl, ok := logEntry["level"].(string); ok {
				entry.Level = lvl
			}
			if msg, ok := logEntry["message"].(string); ok {
				entry.Message = msg
			}
			return entry
		}
	}

	// 解析文本格式: [2024-01-01 12:00:00.000] INFO message
	parts := strings.SplitN(line, "] ", 2)
	if len(parts) == 2 {
		// 提取时间戳
		if strings.HasPrefix(parts[0], "[") {
			timestamp := strings.TrimPrefix(parts[0], "[")
			entry.Timestamp = timestamp
		}

		// 提取级别和消息
		rest := parts[1]
		levelParts := strings.SplitN(rest, " ", 2)
		if len(levelParts) == 2 {
			entry.Level = levelParts[0]
			entry.Message = levelParts[1]
		} else {
			entry.Message = rest
		}
	} else {
		entry.Message = line
	}

	return entry
}
