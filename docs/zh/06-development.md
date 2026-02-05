# 开发指南

## 开发环境设置

### 前置要求

- Go 1.21 或更高版本
- Node.js 18+ 和 npm/yarn
- Git
- SQLite（用于本地开发）

### 设置步骤

1. **克隆仓库**
```bash
git clone https://github.com/your-org/lingproxy.git
cd lingproxy
```

2. **后端设置**
```bash
cd backend
go mod download
```

3. **前端设置**
```bash
cd frontend
npm install
```

4. **配置**
```bash
cp backend/configs/config.yaml.example backend/configs/config.yaml
# 根据需要编辑 config.yaml
```

## 代码结构

### 后端结构

遵循现有结构：
- `cmd/`: 应用入口
- `internal/handler/`: HTTP 处理器
- `internal/service/`: 业务逻辑
- `internal/storage/`: 数据持久化
- `internal/middleware/`: HTTP 中间件
- `internal/config/`: 配置管理

### 前端结构

- `src/views/`: 页面组件
- `src/components/`: 可复用组件
- `src/api/`: API 客户端
- `src/router/`: 路由配置
- `src/locales/`: 国际化语言包
  - `zh/`: 中文语言包
  - `en/`: 英文语言包
  - `index.js`: i18n 配置

## 编码规范

### Go 代码风格

- 遵循 [Effective Go](https://go.dev/doc/effective_go) 指南
- 使用 `gofmt` 格式化
- 遵循命名约定：
  - 导出函数：PascalCase
  - 未导出函数：camelCase
  - 常量：UPPER_SNAKE_CASE

### 错误处理

- 始终检查和处理错误
- 从函数返回错误，不要忽略它们
- 使用描述性错误消息

### 测试

- 为新功能编写单元测试
- 测试文件应命名为 `*_test.go`
- 在适当的地方使用表驱动测试

示例：
```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name string
        input string
        expected string
    }{
        {"test1", "input1", "expected1"},
        {"test2", "input2", "expected2"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Function(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## 添加新功能

### 1. 添加新处理器

1. 创建处理器文件：`internal/handler/feature.go`
2. 实现处理器方法
3. 在 `internal/router/router.go` 中注册路由

示例：
```go
package handler

type FeatureHandler struct {
    storage *storage.StorageFacade
}

func NewFeatureHandler(storage *storage.StorageFacade) *FeatureHandler {
    return &FeatureHandler{storage: storage}
}

func (h *FeatureHandler) GetFeature(c *gin.Context) {
    // 实现
}
```

### 2. 添加新服务

1. 创建服务文件：`internal/service/feature_service.go`
2. 实现业务逻辑
3. 使用存储层进行数据操作

### 3. 添加新存储方法

1. 在 `internal/storage/storage.go` 中添加接口方法
2. 在 `memory_storage.go` 和 `gorm_storage.go` 中实现
3. 添加到 `storage_facade.go`

### 4. 添加新模型

1. 在 `internal/storage/models.go` 中添加模型定义
2. 更新存储实现
3. 如需要，添加迁移

## 运行测试

### 后端测试

```bash
cd backend
go test ./...
go test -v ./internal/handler  # 详细输出
go test -cover ./...            # 覆盖率报告
```

### 前端测试

```bash
cd frontend
npm test
```

## 构建

### 后端

```bash
cd backend
go build -o lingproxy cmd/main.go
```

### 前端

```bash
cd frontend
npm run build
```

## 调试

### 后端调试

1. 使用日志：`logger.Debug()`, `logger.Info()` 等
2. 使用 Go 调试器（delve）
3. 检查 `./logs/lingproxy.log` 中的日志

### 前端调试

1. 使用浏览器 DevTools
2. 检查控制台错误
3. 使用 Vue DevTools 扩展

## Git 工作流

### 分支命名

- `main`: 生产就绪代码
- `develop`: 开发分支
- `feature/feature-name`: 新功能
- `bugfix/bug-name`: 错误修复
- `hotfix/hotfix-name`: 热修复

### 提交消息

遵循约定式提交：
- `feat`: 新功能
- `fix`: 错误修复
- `docs`: 文档更改
- `style`: 代码风格更改
- `refactor`: 代码重构
- `test`: 测试更改
- `chore`: 构建/工具更改

示例：
```
feat: 添加 LLM 资源批量导入
fix: 修复认证问题
docs: 更新 API 文档
```

## Pull Request 流程

1. 创建功能分支
2. 进行更改并提交
3. 编写/更新测试
4. 如需要，更新文档
5. 创建 pull request
6. 处理审查评论
7. 批准后合并

## 代码审查指南

### 审查内容

- 代码正确性
- 代码风格和格式
- 测试覆盖率
- 文档更新
- 性能影响
- 安全考虑

### 审查清单

- [ ] 代码遵循风格指南
- [ ] 包含测试且通过
- [ ] 文档已更新
- [ ] 无安全漏洞
- [ ] 性能可接受
- [ ] 错误处理正确

## 常见任务

### 添加新 API 端点

1. 添加处理器方法
2. 在路由器中添加路由
3. 添加 Swagger 文档
4. 编写测试
5. 更新 API 文档

### 添加新配置选项

1. 添加到 `config.go` 结构体
2. 在 `setDefaults()` 中添加默认值
3. 如需要，添加验证
4. 更新配置文档

### 添加新数据库字段

1. 在 `models.go` 中更新模型
2. 更新存储实现
3. 如需要，添加迁移
4. 更新 API 处理器/服务

### 前端国际化开发

1. **添加新的翻译键**
   - 在 `frontend/src/locales/zh/index.js` 中添加中文翻译
   - 在 `frontend/src/locales/en/index.js` 中添加英文翻译

2. **在组件中使用翻译**
   ```vue
   <script setup>
   import { useI18n } from 'vue-i18n'
   const { t } = useI18n()
   </script>
   
   <template>
     <div>{{ $t('common.home') }}</div>
   </template>
   ```

3. **表单验证规则国际化**
   ```javascript
   const rules = computed(() => ({
     name: [
       { required: true, message: t('form.nameRequired'), trigger: 'blur' }
     ]
   }))
   ```

4. **Element Plus 消息国际化**
   ```javascript
   ElMessage.success(t('common.operationSuccess'))
   ```

5. **语言切换**
   - 语言设置保存在 `localStorage` 中，键为 `lingproxy_locale`
   - 切换语言会自动更新 Element Plus 组件语言

## 故障排除

### 常见问题

**问题**：Go 模块错误
**解决方案**：运行 `go mod tidy`

**问题**：前端构建错误
**解决方案**：删除 `node_modules` 和 `package-lock.json`，然后 `npm install`

**问题**：数据库连接错误
**解决方案**：检查数据库配置并确保数据库正在运行

**问题**：端口已被使用
**解决方案**：在 `config.yaml` 中更改端口或停止使用该端口的进程

## 资源

- [Go 文档](https://go.dev/doc/)
- [Gin 框架](https://gin-gonic.com/docs/)
- [Vue 3 文档](https://vuejs.org/)
- [Element Plus](https://element-plus.org/)
