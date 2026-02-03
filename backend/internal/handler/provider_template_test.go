package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lingproxy/lingproxy/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

// TestDownloadImportTemplate_GenerateFile 测试生成Excel模板文件并保存到项目根目录
func TestDownloadImportTemplate_GenerateFile(t *testing.T) {
	// 设置gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试用的handler（需要storage，这里创建一个mock）
	// 注意：实际测试中可能需要mock storage
	storage := setupTestStorage(t)
	handler := NewLLMResourceHandler(storage)

	// 创建HTTP请求
	req, _ := http.NewRequest("GET", "/api/v1/llm-resources/import/template", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建gin上下文
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 调用处理器
	handler.DownloadImportTemplate(c)

	// 检查响应状态码
	assert.Equal(t, http.StatusOK, w.Code, "响应状态码应该是200")
	assert.NotEmpty(t, w.Body.Bytes(), "响应体不应该为空")

	// 获取项目根目录（从backend/internal/handler目录向上三级到项目根目录）
	projectRoot := filepath.Join("..", "..", "..")
	absPath, _ := filepath.Abs(projectRoot)
	outputPath := filepath.Join(absPath, "test_llm_resources_import_template.xlsx")

	// 保存文件到项目根目录
	err := os.WriteFile(outputPath, w.Body.Bytes(), 0644)
	assert.NoError(t, err, "保存文件应该成功")

	t.Logf("Excel文件已保存到: %s", outputPath)
	t.Logf("文件大小: %d bytes", len(w.Body.Bytes()))

	// 验证文件格式
	validateExcelFile(t, outputPath, w.Body.Bytes())
}

// TestDownloadImportTemplate_ExcelizeDirect 直接使用excelize生成文件进行对比测试
func TestDownloadImportTemplate_ExcelizeDirect(t *testing.T) {
	// 创建Excel文件
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "LLM资源导入模板"
	index, err := f.NewSheet(sheetName)
	assert.NoError(t, err, "创建工作表应该成功")

	f.SetActiveSheet(index)

	// 删除默认Sheet1
	err = f.DeleteSheet("Sheet1")
	if err != nil {
		t.Logf("删除默认Sheet失败（可忽略）: %v", err)
	}

	// 设置表头
	headers := []string{"资源名称", "模型类别", "驱动", "模型标识", "基础URL", "API密钥", "状态"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		err := f.SetCellValue(sheetName, cell, header)
		assert.NoError(t, err, "设置表头应该成功")
	}

	// 设置示例数据
	examples := [][]interface{}{
		{"OpenAI GPT-4", "chat", "openai", "gpt-4", "https://api.openai.com/v1", "sk-xxxxxxxxxxxxx", "active"},
		{"OpenAI GPT-3.5", "chat", "openai", "gpt-3.5-turbo", "https://api.openai.com/v1", "sk-yyyyyyyyyyyyy", "active"},
		{"OpenAI GPT-4o", "chat", "openai", "gpt-4o", "https://api.openai.com/v1", "sk-zzzzzzzzzzzzz", "active"},
	}

	for rowIdx, example := range examples {
		for colIdx, value := range example {
			cell := string(rune('A'+colIdx)) + string(rune('0'+rowIdx+2))
			err := f.SetCellValue(sheetName, cell, value)
			assert.NoError(t, err, "设置示例数据应该成功")
		}
	}

	// 设置列宽
	columnWidths := map[string]float64{
		"A": 20, "B": 12, "C": 10, "D": 20, "E": 30, "F": 25, "G": 10,
	}
	for col, width := range columnWidths {
		err := f.SetColWidth(sheetName, col, col, width)
		if err != nil {
			t.Logf("设置列宽失败（可忽略）: %v", err)
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

	// 将文件写入buffer
	var buf bytes.Buffer
	err = f.Write(&buf)
	assert.NoError(t, err, "写入Excel文件应该成功")

	// 保存到项目根目录（从backend/internal/handler目录向上三级到项目根目录）
	projectRoot := filepath.Join("..", "..", "..")
	absPath, _ := filepath.Abs(projectRoot)
	outputPath := filepath.Join(absPath, "test_direct_excelize_template.xlsx")

	err = os.WriteFile(outputPath, buf.Bytes(), 0644)
	assert.NoError(t, err, "保存文件应该成功")

	t.Logf("直接使用excelize生成的文件已保存到: %s", outputPath)
	t.Logf("文件大小: %d bytes", buf.Len())

	// 验证文件格式
	validateExcelFile(t, outputPath, buf.Bytes())
}

// validateExcelFile 验证Excel文件格式
func validateExcelFile(t *testing.T, filePath string, fileData []byte) {
	// 1. 检查文件大小
	assert.Greater(t, len(fileData), 0, "文件大小应该大于0")
	t.Logf("文件大小: %d bytes", len(fileData))

	// 2. 检查ZIP文件头（Excel文件本质是ZIP）
	// ZIP文件头: 50 4B 03 04 (PK..)
	if len(fileData) >= 4 {
		assert.Equal(t, byte(0x50), fileData[0], "文件头第1字节应该是0x50 (P)")
		assert.Equal(t, byte(0x4B), fileData[1], "文件头第2字节应该是0x4B (K)")
		assert.Equal(t, byte(0x03), fileData[2], "文件头第3字节应该是0x03")
		assert.Equal(t, byte(0x04), fileData[3], "文件头第4字节应该是0x04")
		t.Logf("文件头验证通过: ZIP格式")
	}

	// 3. 尝试用excelize打开文件
	f, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		t.Fatalf("无法用excelize打开文件: %v", err)
	}
	defer f.Close()

	// 4. 检查工作表
	sheetList := f.GetSheetList()
	assert.NotEmpty(t, sheetList, "应该至少有一个工作表")
	t.Logf("工作表列表: %v", sheetList)

	// 5. 检查第一个工作表的内容
	if len(sheetList) > 0 {
		sheetName := sheetList[0]

		// 读取表头
		headers := make([]string, 0)
		for col := 'A'; col <= 'G'; col++ {
			cell := string(col) + "1"
			value, err := f.GetCellValue(sheetName, cell)
			if err == nil && value != "" {
				headers = append(headers, value)
			}
		}
		t.Logf("表头: %v", headers)
		assert.Equal(t, 7, len(headers), "应该有7个表头字段")
		assert.Contains(t, headers, "驱动", "表头应该包含'驱动'字段")
		assert.NotContains(t, headers, "服务提供商", "表头不应该包含'服务提供商'字段")

		// 读取示例数据
		for row := 2; row <= 4; row++ {
			rowData := make([]string, 0)
			for col := 'A'; col <= 'G'; col++ {
				cell := string(col) + string(rune('0'+row))
				value, _ := f.GetCellValue(sheetName, cell)
				rowData = append(rowData, value)
			}
			t.Logf("第%d行数据: %v", row, rowData)
		}
	}

	t.Logf("文件验证通过: %s", filePath)
}

// setupTestStorage 设置测试用的storage（简化版本）
func setupTestStorage(t *testing.T) *storage.StorageFacade {
	// 创建一个内存storage用于测试
	// 注意：DownloadImportTemplate函数实际上不需要storage，但handler构造函数需要
	memStorage := storage.NewMemoryStorage()
	return storage.NewStorageFacade(memStorage)
}
