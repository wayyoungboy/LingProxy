package handler

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/config"
	"github.com/lingproxy/lingproxy/pkg/logger"
)

var (
	startTime = time.Now()
)

// SystemHandler 系统信息处理器
type SystemHandler struct{}

// NewSystemHandler 创建新的系统信息处理器
func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// SystemInfo 系统信息
type SystemInfo struct {
	SystemName    string  `json:"system_name"`
	SystemVersion string  `json:"system_version"`
	Environment   string  `json:"environment"`
	StartTime     string  `json:"start_time"`
	CurrentTime   string  `json:"current_time"`
	Uptime        string  `json:"uptime"`
	CPUCores      int     `json:"cpu_cores"`
	CPUUsage      float64 `json:"cpu_usage"`
	MemoryTotal   int64   `json:"memory_total"`
	MemoryUsed    int64   `json:"memory_used"`
	MemoryUsage   float64 `json:"memory_usage"`
	DiskTotal     int64   `json:"disk_total"`
	DiskUsed      int64   `json:"disk_used"`
	DiskUsage     float64 `json:"disk_usage"`
	GoVersion     string  `json:"go_version"`
	OS            string  `json:"os"`
	Arch          string  `json:"arch"`
}

// GetSystemInfo 获取系统信息
// @Summary Get system information
// @Description Get system runtime information and resource usage
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "System information"
// @Router /api/v1/system/info [get]
func (h *SystemHandler) GetSystemInfo(c *gin.Context) {
	cfg := config.C

	// 计算运行时长
	uptime := time.Since(startTime)
	uptimeStr := formatDuration(uptime)

	// 获取内存信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memoryUsed := int64(m.Alloc)
	memoryTotal := int64(m.Sys)
	memoryUsage := 0.0
	if memoryTotal > 0 {
		memoryUsage = float64(memoryUsed) / float64(memoryTotal) * 100
	}

	// 获取CPU核心数
	cpuCores := runtime.NumCPU()

	// 获取Go版本和系统信息
	goVersion := runtime.Version()
	osName := runtime.GOOS
	arch := runtime.GOARCH

	info := SystemInfo{
		SystemName:    cfg.App.Name,
		SystemVersion: cfg.App.Version,
		Environment:   cfg.App.Environment,
		StartTime:     startTime.Format(time.RFC3339),
		CurrentTime:   time.Now().Format(time.RFC3339),
		Uptime:        uptimeStr,
		CPUCores:      cpuCores,
		CPUUsage:      0.0, // 需要gopsutil库才能获取准确的CPU使用率
		MemoryTotal:   memoryTotal,
		MemoryUsed:    memoryUsed,
		MemoryUsage:   memoryUsage,
		DiskTotal:     0, // 需要gopsutil库才能获取磁盘信息
		DiskUsed:      0,
		DiskUsage:     0.0,
		GoVersion:     goVersion,
		OS:            osName,
		Arch:          arch,
	}

	logger.Info("获取系统信息成功")
	c.JSON(http.StatusOK, gin.H{"data": info})
}

// formatDuration 格式化时长
func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd%dh%dm%ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
