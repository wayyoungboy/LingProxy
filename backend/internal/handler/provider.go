package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/pkg/logger"
	"github.com/lingproxy/lingproxy/internal/storage"
	"github.com/xuri/excelize/v2"
)

// LLMResourceHandler LLM资源处理器
type LLMResourceHandler struct {
	storage *storage.StorageFacade
}

// NewLLMResourceHandler 创建新的LLM资源处理器
func NewLLMResourceHandler(storage *storage.StorageFacade) *LLMResourceHandler {
	return &LLMResourceHandler{
		storage: storage,
	}
}

// ListLLMResources godoc
// @Summary List all LLM resources
// @Description Get a list of all AI service LLM resources
// @Tags llm-resources
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of LLM resources"
// @Router /api/v1/llm-resources [get]
func (h *LLMResourceHandler) ListLLMResources(c *gin.Context) {
	resources, err := h.storage.ListLLMResources()
	if err != nil {
		logger.Error("获取LLM资源列表失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Info("获取LLM资源列表成功", logger.F("count", len(resources)))
	c.JSON(http.StatusOK, gin.H{"data": resources})
}

// GetLLMResource godoc
// @Summary Get LLM resource by ID
// @Description Get a specific AI service LLM resource by ID
// @Tags llm-resources
// @Accept json
// @Produce json
// @Param id path string true "LLM Resource ID"
// @Success 200 {object} map[string]interface{} "LLM resource details"
// @Failure 404 {object} map[string]string "LLM resource not found"
// @Router /api/v1/llm-resources/{id} [get]
func (h *LLMResourceHandler) GetLLMResource(c *gin.Context) {
	id := c.Param("id")
	resource, err := h.storage.GetLLMResource(id)
	if err != nil {
		logger.Warn("获取LLM资源失败：资源不存在", logger.F("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "LLM resource not found"})
		return
	}
	logger.Info("获取LLM资源成功", logger.F("id", id), logger.F("name", resource.Name))
	c.JSON(http.StatusOK, gin.H{"data": resource})
}

// CreateLLMResource godoc
// @Summary Create a new LLM resource
// @Description Create a new AI service LLM resource configuration
// @Tags llm-resources
// @Accept json
// @Produce json
// @Param resource body storage.LLMResource true "LLM resource configuration"
// @Success 201 {object} map[string]interface{} "Created LLM resource"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/llm-resources [post]
func (h *LLMResourceHandler) CreateLLMResource(c *gin.Context) {
	var resource storage.LLMResource
	if err := c.ShouldBindJSON(&resource); err != nil {
		logger.Error("创建LLM资源失败", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证必填字段
	if resource.Name == "" {
		logger.Warn("创建LLM资源失败：名称为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "名称是必填项"})
		return
	}
	if resource.Type == "" {
		logger.Warn("创建LLM资源失败：模型类别为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型类别是必填项"})
		return
	}
	// 驱动固定为openai
	if resource.Driver == "" {
		resource.Driver = "openai"
	} else if resource.Driver != "openai" {
		logger.Warn("创建LLM资源失败：不支持的驱动", logger.F("driver", resource.Driver))
		c.JSON(http.StatusBadRequest, gin.H{"error": "目前仅支持openai驱动"})
		return
	}
	if resource.Model == "" {
		logger.Warn("创建LLM资源失败：模型标识为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型标识是必填项"})
		return
	}
	if resource.BaseURL == "" {
		logger.Warn("创建LLM资源失败：Base URL为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Base URL是必填项"})
		return
	}
	if resource.APIKey == "" {
		logger.Warn("创建LLM资源失败：API Key为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "API Key是必填项"})
		return
	}

	if err := h.storage.CreateLLMResource(&resource); err != nil {
		logger.Error("创建LLM资源失败", logger.F("error", err.Error()), logger.F("name", resource.Name))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("创建LLM资源成功", logger.F("id", resource.ID), logger.F("name", resource.Name), logger.F("model", resource.Model))
	c.JSON(http.StatusCreated, gin.H{"data": resource})
}

// UpdateLLMResource godoc
// @Summary Update an existing LLM resource
// @Description Update an existing AI service LLM resource configuration
// @Tags llm-resources
// @Accept json
// @Produce json
// @Param id path string true "LLM Resource ID"
// @Param resource body storage.LLMResource true "LLM resource configuration"
// @Success 200 {object} map[string]interface{} "Updated LLM resource"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "LLM resource not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/llm-resources/{id} [put]
func (h *LLMResourceHandler) UpdateLLMResource(c *gin.Context) {
	id := c.Param("id")
	var resource storage.LLMResource
	if err := c.ShouldBindJSON(&resource); err != nil {
		logger.Error("更新LLM资源失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证必填字段
	if resource.Name == "" {
		logger.Warn("更新LLM资源失败：名称为空", logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "名称是必填项"})
		return
	}
	if resource.Type == "" {
		logger.Warn("更新LLM资源失败：模型类别为空", logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型类别是必填项"})
		return
	}
	// 驱动固定为openai
	if resource.Driver == "" {
		resource.Driver = "openai"
	} else if resource.Driver != "openai" {
		logger.Warn("更新LLM资源失败：不支持的驱动", logger.F("id", id), logger.F("driver", resource.Driver))
		c.JSON(http.StatusBadRequest, gin.H{"error": "目前仅支持openai驱动"})
		return
	}
	if resource.Model == "" {
		logger.Warn("更新LLM资源失败：模型标识为空", logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "模型标识是必填项"})
		return
	}
	if resource.BaseURL == "" {
		logger.Warn("更新LLM资源失败：Base URL为空", logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Base URL是必填项"})
		return
	}
	if resource.APIKey == "" {
		logger.Warn("更新LLM资源失败：API Key为空", logger.F("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"error": "API Key是必填项"})
		return
	}

	// 确保资源ID与路径参数一致
	resource.ID = id

	if err := h.storage.UpdateLLMResource(&resource); err != nil {
		if err == storage.ErrNotFound {
			logger.Warn("更新LLM资源失败：资源不存在", logger.F("id", id))
			c.JSON(http.StatusNotFound, gin.H{"error": "LLM resource not found"})
			return
		}
		logger.Error("更新LLM资源失败", logger.F("error", err.Error()), logger.F("id", id), logger.F("name", resource.Name))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("更新LLM资源成功", logger.F("id", id), logger.F("name", resource.Name), logger.F("model", resource.Model))
	c.JSON(http.StatusOK, gin.H{"data": resource})
}

// DeleteLLMResource godoc
// @Summary Delete an LLM resource
// @Description Delete an AI service LLM resource configuration
// @Tags llm-resources
// @Accept json
// @Produce json
// @Param id path string true "LLM Resource ID"
// @Success 204 {object} nil "LLM resource deleted"
// @Failure 404 {object} map[string]string "LLM resource not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/llm-resources/{id} [delete]
func (h *LLMResourceHandler) DeleteLLMResource(c *gin.Context) {
	id := c.Param("id")

	if err := h.storage.DeleteLLMResource(id); err != nil {
		if err == storage.ErrNotFound {
			logger.Warn("删除LLM资源失败：资源不存在", logger.F("id", id))
			c.JSON(http.StatusNotFound, gin.H{"error": "LLM resource not found"})
			return
		}
		logger.Error("删除LLM资源失败", logger.F("error", err.Error()), logger.F("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info("删除LLM资源成功", logger.F("id", id))
	c.Status(http.StatusNoContent)
}

// ImportLLMResources 批量导入LLM资源
// @Summary Import LLM resources from Excel
// @Description Import multiple LLM resources from Excel file
// @Tags llm-resources
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel file"
// @Success 200 {object} map[string]interface{} "Import result"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/llm-resources/import [post]
func (h *LLMResourceHandler) ImportLLMResources(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		logger.Warn("批量导入LLM资源失败：文件获取失败", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件获取失败: " + err.Error()})
		return
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		logger.Error("批量导入LLM资源失败：文件打开失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件打开失败: " + err.Error()})
		return
	}
	defer src.Close()

	// 读取Excel文件
	f, err := excelize.OpenReader(src)
	if err != nil {
		logger.Error("批量导入LLM资源失败：Excel文件解析失败", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Excel文件解析失败: " + err.Error()})
		return
	}
	defer f.Close()

	// 读取第一个工作表
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		logger.Error("批量导入LLM资源失败：读取工作表失败", logger.F("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "读取工作表失败: " + err.Error()})
		return
	}

	if len(rows) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Excel文件至少需要包含表头和一行数据"})
		return
	}

	// 解析表头
	header := rows[0]
	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	// 验证必需字段（每个字段至少需要一个别名存在）
	requiredFieldGroups := map[string][]string{
		"资源名称":  {"资源名称", "name"},
		"模型类别":  {"模型类别", "type"},
		"驱动":    {"驱动", "driver"},
		"模型标识":  {"模型标识", "model"},
		"基础URL": {"基础url", "base_url", "baseurl"},
		"API密钥": {"api密钥", "api_key", "apikey"},
	}
	missingFields := []string{}
	for fieldName, aliases := range requiredFieldGroups {
		found := false
		for _, alias := range aliases {
			if _, exists := headerMap[strings.ToLower(alias)]; exists {
				found = true
				break
			}
		}
		if !found {
			missingFields = append(missingFields, fieldName)
		}
	}
	if len(missingFields) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("缺少必需字段: %v", missingFields)})
		return
	}

	// 获取列索引
	nameCol := getColumnIndex(headerMap, []string{"资源名称", "name"})
	typeCol := getColumnIndex(headerMap, []string{"模型类别", "type"})
	driverCol := getColumnIndex(headerMap, []string{"驱动", "driver"})
	modelCol := getColumnIndex(headerMap, []string{"模型标识", "model"})
	baseURLCol := getColumnIndex(headerMap, []string{"基础url", "base_url", "baseurl"})
	apiKeyCol := getColumnIndex(headerMap, []string{"api密钥", "api_key", "apikey"})
	statusCol := getColumnIndex(headerMap, []string{"状态", "status"})

	// 解析数据行
	successCount := 0
	failCount := 0
	errors := []string{}

	for i, row := range rows[1:] {
		rowNum := i + 2 // Excel行号（从2开始，因为第1行是表头）

		// 跳过空行
		if len(row) == 0 {
			continue
		}

		// 获取字段值
		name := getCellValue(row, nameCol)
		typeVal := getCellValue(row, typeCol)
		driver := getCellValue(row, driverCol)
		model := getCellValue(row, modelCol)
		baseURL := getCellValue(row, baseURLCol)
		apiKey := getCellValue(row, apiKeyCol)
		status := getCellValue(row, statusCol)
		if status == "" {
			status = "active"
		}
		// 驱动默认为openai，如果为空或不是openai则设置为openai
		if driver == "" || strings.ToLower(driver) != "openai" {
			driver = "openai"
		}

		// 验证必填字段
		if name == "" || typeVal == "" || model == "" || baseURL == "" || apiKey == "" {
			failCount++
			errors = append(errors, fmt.Sprintf("第%d行: 必填字段不能为空", rowNum))
			continue
		}

		// 创建资源
		resource := &storage.LLMResource{
			Name:    name,
			Type:    strings.ToLower(typeVal),
			Driver:  strings.ToLower(driver),
			Model:   model,
			BaseURL: baseURL,
			APIKey:  apiKey,
			Status:  strings.ToLower(status),
		}

		if err := h.storage.CreateLLMResource(resource); err != nil {
			failCount++
			errors = append(errors, fmt.Sprintf("第%d行: %s", rowNum, err.Error()))
			continue
		}

		successCount++
	}

	logger.Info("批量导入LLM资源完成", logger.F("success", successCount), logger.F("fail", failCount))
	c.JSON(http.StatusOK, gin.H{
		"message": "导入完成",
		"success": successCount,
		"fail":    failCount,
		"errors":  errors,
		"total":   successCount + failCount,
	})
}

// DownloadImportTemplate 下载导入模板
// @Summary Download import template
// @Description Download Excel template for importing LLM resources
// @Tags llm-resources
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Success 200 {file} file "Excel template file"
// @Router /api/v1/llm-resources/import/template [get]
func (h *LLMResourceHandler) DownloadImportTemplate(c *gin.Context) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "LLM资源导入模板"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		logger.Error("创建Excel模板失败", logger.F("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建Excel模板失败"})
		return
	}

	// 先设置新sheet为活动sheet，再删除默认Sheet1
	f.SetActiveSheet(index)

	// 删除默认Sheet1
	if err := f.DeleteSheet("Sheet1"); err != nil {
		logger.Warn("删除默认Sheet失败", logger.F("error", err.Error()))
		// 删除失败不影响文件生成，继续执行
	}

	// 设置表头
	headers := []string{"资源名称", "模型类别", "驱动", "模型标识", "基础URL", "API密钥", "状态"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			logger.Error("设置表头失败", logger.F("error", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "设置表头失败"})
			return
		}
	}

	// 设置示例数据（只包含核心字段，不包含模型元数据）
	examples := [][]interface{}{
		{"OpenAI GPT-4", "chat", "openai", "gpt-4", "https://api.openai.com/v1", "sk-xxxxxxxxxxxxx", "active"},
		{"OpenAI GPT-3.5", "chat", "openai", "gpt-3.5-turbo", "https://api.openai.com/v1", "sk-yyyyyyyyyyyyy", "active"},
		{"OpenAI GPT-4o", "chat", "openai", "gpt-4o", "https://api.openai.com/v1", "sk-zzzzzzzzzzzzz", "active"},
	}

	for rowIdx, example := range examples {
		for colIdx, value := range example {
			cell := fmt.Sprintf("%c%d", 'A'+colIdx, rowIdx+2)
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				logger.Error("设置示例数据失败", logger.F("error", err.Error()))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "设置示例数据失败"})
				return
			}
		}
	}

	// 设置列宽（根据字段长度优化）
	columnWidths := map[string]float64{
		"A": 20, // 资源名称
		"B": 12, // 模型类别
		"C": 10, // 驱动
		"D": 20, // 模型标识
		"E": 30, // 基础URL
		"F": 25, // API密钥
		"G": 10, // 状态
	}
	for col, width := range columnWidths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			logger.Warn("设置列宽失败", logger.F("column", col), logger.F("error", err.Error()))
		}
	}

	// 设置表头样式
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", "G1", headerStyle)
	}

	// 设置响应头（必须在写入数据之前设置）
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", `attachment; filename="llm_resources_import_template.xlsx"`)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Status(http.StatusOK)

	// 写入响应
	if err := f.Write(c.Writer); err != nil {
		logger.Error("写入Excel模板失败", logger.F("error", err.Error()))
		// 注意：此时响应头已经设置，不能再用c.JSON，需要直接写入错误信息
		if !c.Writer.Written() {
			c.String(http.StatusInternalServerError, "写入Excel模板失败: %v", err)
		}
		return
	}

	logger.Info("Excel模板下载成功")
}

// getColumnIndex 获取列索引（支持多个可能的列名）
func getColumnIndex(headerMap map[string]int, possibleNames []string) int {
	for _, name := range possibleNames {
		if idx, exists := headerMap[strings.ToLower(name)]; exists {
			return idx
		}
	}
	return -1
}

// getCellValue 安全获取单元格值
func getCellValue(row []string, colIndex int) string {
	if colIndex < 0 || colIndex >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[colIndex])
}
