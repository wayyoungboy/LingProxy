# LingProxy 开发 Agent 设计

## 基本信息
- **名称**：LingProxy Dev Assistant
- **版本**：2.1.0
- **项目类型**：AI API 网关（OpenAI 兼容）
- **技术栈**：
  - **后端**：Go 1.24+ + Gin + GORM + Viper
  - **前端**：Vue 3 + Element Plus + Vite + Vue Router + Vue I18n
  - **数据库**：SQLite/MySQL/PostgreSQL
  - **其他**：Swagger、Excelize（批量导入）

## 项目架构理解

### 核心模块
1. **路由策略系统**：支持随机、轮询、加权、模型匹配、正则匹配、优先级、故障转移等策略
2. **Token 管理**：请求端 Token 管理，支持策略绑定（"ling-" 前缀 Token），支持 Token 重置和复制
3. **LLM 资源管理**：
   - 基于驱动的架构（Driver-based），目前支持 OpenAI 驱动
   - 多模型类别：chat（对话）、image（生图）、embedding（嵌入）、rerank（重排序）、audio（语音）、video（视频）
   - 批量导入/导出：支持 Excel 模板和 JSON 格式批量导入
   - 资源测试：支持测试资源连通性，验证配置正确性
   - 搜索功能：支持按资源名称、BaseURL、模型标识的模糊搜索
4. **认证系统**：全局认证开关，管理员登录（用户名/密码），Token 认证，API Key 认证
5. **系统设置**：动态配置管理（基础、缓存、限流、安全、日志、负载均衡、熔断器）
6. **日志管理**：支持日志文件查看、下载、清理
7. **统计功能**：系统统计、LLM 资源使用统计、用户统计
8. **国际化支持**：前端支持中英文切换（Vue I18n）
9. **OpenAI 兼容 API**：`/llm/v1/*` 路径前缀

### 代码组织
```
backend/
├── cmd/            # 应用入口
├── configs/        # 配置文件（config.yaml）
├── internal/
│   ├── handler/   # HTTP 处理器（RESTful API）
│   ├── service/   # 业务逻辑层
│   ├── storage/   # 数据存储层（接口 + 内存/GORM 实现）
│   ├── middleware/# 中间件（认证、CORS、日志、限流、代理）
│   ├── config/    # 配置管理（Viper）
│   ├── router/    # 路由配置
│   ├── client/    # AI 服务客户端（OpenAI、Embedding）
│   └── pkg/       # 内部工具包（logger、balancer、monitor、password）
├── swagger/        # Swagger 文档
└── examples/       # 示例代码

frontend/
├── src/
│   ├── api/       # API 客户端封装（统一使用 axios）
│   ├── components/# Vue 组件（Layout、Sidebar 等）
│   ├── views/     # 页面视图
│   ├── router/    # Vue Router 配置
│   ├── locales/   # 国际化资源（zh、en）
│   ├── utils/     # 工具函数（constants、helpers）
│   ├── composables/# Vue Composables（useLoading、usePagination）
│   └── config/    # 配置文件（menu.js）
├── public/         # 静态资源
└── package.json    # 依赖管理
```

## Agent 核心能力

### 1. 项目分析与理解
- **架构分析**：理解分层架构（Handler → Service → Storage）
- **依赖关系**：分析模块间依赖，识别循环依赖
- **API 映射**：理解 RESTful API 结构和 OpenAI 兼容 API
- **数据模型**：理解存储模型（User、Token、Policy、LLMResource、Model 等）
- **路由分析**：理解路由配置和认证中间件应用逻辑

### 2. 功能开发支持

#### 2.1 新功能开发流程
1. **需求分析**：理解功能需求，确定影响范围
2. **数据模型设计**：在 `backend/internal/storage/models.go` 中定义/更新模型
3. **存储层实现**：
   - 在 `backend/internal/storage/storage.go` 中定义接口
   - 在 `memory_storage.go` 和 `gorm_storage.go` 中实现
   - 在 `gorm_storage.go` 的 `Init()` 中添加 `AutoMigrate` 迁移
4. **业务逻辑层**：在 `backend/internal/service/` 中实现业务逻辑
5. **HTTP 处理器**：在 `backend/internal/handler/` 中实现 RESTful API
6. **路由注册**：在 `backend/internal/router/router.go` 中注册路由（考虑认证开关）
7. **前端集成**：
   - 在 `frontend/src/api/index.js` 中添加 API 方法
   - 在 `frontend/src/views/` 中创建或更新页面组件
   - 在 `frontend/src/router/index.js` 中注册路由
   - 在 `frontend/src/config/menu.js` 中添加菜单项（如需要）
   - 添加国际化文本（`frontend/src/locales/zh/index.js` 和 `en/index.js`）
8. **Swagger 文档**：更新 API 注释，运行 `swag init` 生成文档

#### 2.2 特定功能模式

**策略管理功能**：
- 策略模板（PolicyTemplate）：内置模板，JSON Schema 定义参数
- 策略实例（Policy）：基于模板创建，支持参数配置
- 策略执行器（PolicyExecutor）：实现不同策略类型的执行逻辑
- Token 策略绑定：Token 可关联策略，影响请求路由

**Token 管理功能**：
- Token 生成：自动生成 "ling-" 前缀的 Token
- Token 策略绑定：支持为 Token 绑定路由策略
- Token 验证：在认证中间件中验证 Token

**LLM 资源管理功能**：
- 驱动架构：基于 Driver 的架构设计，目前仅支持 "openai" 驱动
- 批量导入：支持 Excel 模板（通过 `excelize` 库）和 JSON 格式批量导入
- 资源测试：支持测试 chat 和 embedding 类型资源的连通性
- 搜索功能：前端支持按名称、BaseURL、模型标识的模糊搜索
- 测试状态：资源包含 `test_status` 字段（untested、passed、failed）

**系统设置功能**：
- 配置读取：使用 Viper 读取 `backend/configs/config.yaml`
- 配置更新：动态更新配置并持久化到文件
- 配置验证：验证配置项的有效性
- 重启提示：识别需要重启的配置项

**日志管理功能**：
- 日志文件列表：获取系统日志文件列表
- 日志查看：支持按文件查看日志内容
- 日志下载：支持下载日志文件
- 日志清理：支持清理指定日志文件

**统计功能**：
- 系统统计：总请求数、总用户数、总资源数、成功率、平均响应时间
- LLM 资源使用统计：按资源分组，包含 Token 使用量、请求数、成功率等
- 用户统计：单个用户的详细统计信息

### 3. 代码规范与最佳实践

#### 3.1 Go 后端规范
- **错误处理**：使用 `storage.ErrNotFound` 等标准错误
- **日志记录**：使用 `pkg/logger` 包，记录关键操作和错误
- **认证检查**：根据 `cfg.Security.Auth.Enabled` 条件应用认证中间件
- **ID 生成**：使用 `generateID()` 生成唯一标识
- **时间戳**：统一使用 `time.Now()` 设置 `CreatedAt` 和 `UpdatedAt`
- **JSON 序列化**：敏感字段使用 `json:"-"` 隐藏

#### 3.2 前端规范
- **API 调用**：
  - 统一使用 `frontend/src/api/index.js` 中的 API 客户端
  - API 客户端已配置请求/响应拦截器，自动处理认证 Token
  - 错误处理已在拦截器中统一处理，使用 i18n 显示错误消息
  - 支持 blob 响应类型（用于文件下载）
  - API 方法示例：
    ```javascript
    // 在组件中使用
    import api from '@/api'
    
    // 调用 API
    try {
      const response = await api.getLLMResources({ search: 'gpt' })
      // response 已经是 response.data，不需要再取 .data
      console.log(response.data) // 实际数据
    } catch (error) {
      // 错误已在拦截器中处理，这里可以处理业务逻辑
      console.error('Failed to fetch resources:', error)
    }
    
    // 文件下载示例
    const blob = await api.downloadLLMResourcesTemplate()
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'template.xlsx'
    link.click()
    ```
- **组件组织**：
  - 页面组件放在 `views/` 目录
  - 通用组件放在 `components/` 目录（如 Layout、Sidebar）
  - 使用 Vue 3 Composition API 和 `<script setup>` 语法
- **国际化**：
  - 所有用户可见文本必须使用 i18n：`i18n.global.t('key')`
  - 文本资源定义在 `frontend/src/locales/zh/index.js` 和 `en/index.js`
  - Element Plus 组件语言包在 `main.js` 中根据当前语言设置
  - 语言设置保存在 `localStorage` 的 `lingproxy_locale` 键中
  - 使用示例：
    ```javascript
    import { useI18n } from 'vue-i18n'
    
    const { t } = useI18n()
    // 在模板中使用
    // <el-button>{{ t('common.save') }}</el-button>
    
    // 在脚本中使用
    ElMessage.success(t('message.saveSuccess'))
    ```
- **状态管理**：
  - 使用 Vue 3 Composition API（`ref`、`reactive`、`computed`）
  - 使用 Composables（`useLoading`、`usePagination`）复用逻辑
  - 示例：
    ```javascript
    import { ref, computed } from 'vue'
    import { useLoading } from '@/composables/useLoading'
    import { usePagination } from '@/composables/usePagination'
    
    const { loading, withLoading } = useLoading()
    const { page, pageSize, total, handlePageChange } = usePagination()
    
    const resources = ref([])
    const fetchResources = async () => {
      await withLoading(async () => {
        const response = await api.getLLMResources({
          page: page.value,
          page_size: pageSize.value
        })
        resources.value = response.data
        total.value = response.total
      })
    }
    ```
- **错误处理**：
  - API 错误已在拦截器中统一处理（401 会自动跳转登录页）
  - 业务错误使用 Element Plus 的 `ElMessage` 显示
  - 使用 i18n 显示错误消息，确保多语言支持
  - 示例：
    ```javascript
    import { ElMessage } from 'element-plus'
    import { useI18n } from 'vue-i18n'
    
    const { t } = useI18n()
    
    try {
      await api.createLLMResource(resourceData)
      ElMessage.success(t('message.createSuccess'))
    } catch (error) {
      // 网络错误和 HTTP 错误已在拦截器中处理
      // 这里处理业务逻辑错误
      if (error.response?.status === 400) {
        ElMessage.warning(t('message.validationError'))
      }
    }
    ```
- **表单验证**：使用 Element Plus 的表单验证规则
- **常量管理**：统一使用 `frontend/src/utils/constants.js` 中的常量

#### 3.3 数据库规范
- **迁移**：在 `backend/internal/storage/gorm_storage.go` 的 `Init()` 中使用 `AutoMigrate`
- **查询优化**：使用 GORM 的预加载（`Preload`）和索引
- **事务处理**：需要时使用 GORM 事务（`db.Transaction`）
- **模型定义**：所有模型定义在 `backend/internal/storage/models.go` 中
- **敏感字段**：密码、API Key 等敏感字段使用 `json:"-"` 隐藏，不序列化到 JSON

### 4. Bug 修复与问题排查

#### 4.1 常见问题模式
- **循环依赖**：检查 `service` 和 `handler` 之间的导入关系
- **认证问题**：检查认证中间件是否正确应用，Token 验证逻辑
- **配置问题**：检查 `config.yaml` 格式和默认值设置
- **存储问题**：检查数据库连接和迁移状态
- **路由问题**：检查路由注册顺序和路径匹配

#### 4.2 修复流程
1. **问题定位**：分析错误日志和堆栈跟踪
2. **影响分析**：确定问题影响范围
3. **方案设计**：设计修复方案，考虑向后兼容
4. **实施修复**：应用修复，更新相关代码
5. **测试验证**：运行测试，验证修复效果
6. **文档更新**：更新相关文档和注释

### 5. 测试支持

#### 5.1 测试类型
- **单元测试**：测试 Service 层和 Storage 层逻辑
- **集成测试**：测试 API 端点和数据库交互
- **前端测试**：测试 Vue 组件和 API 调用

#### 5.2 测试文件位置
- 存储层测试：`internal/storage/*_test.go`
- 服务层测试：`internal/service/*_test.go`（如需要）
- 处理器测试：`internal/handler/*_test.go`（如需要）

### 6. 文档维护

#### 6.1 文档类型
- **README.md / README_zh.md**：项目说明、快速开始、API 参考
- **Swagger 文档**：API 文档（自动生成）
- **代码注释**：Go 代码中的 godoc 注释

#### 6.2 文档更新原则
- **同步更新**：代码变更时同步更新文档
- **中英文一致**：保持中英文文档内容一致
- **示例准确**：确保代码示例和 API 示例准确
- **配置说明**：更新配置项说明和默认值

### 7. 部署与运维

#### 7.1 构建流程
- **后端构建**：使用 `go build`，支持交叉编译
- **前端构建**：使用 `npm run build`，生成静态文件
- **Docker 构建**：使用 `Dockerfile` 构建镜像

#### 7.2 配置管理
- **环境配置**：支持 development、staging、production 环境
- **配置验证**：启动时验证配置有效性
- **配置热更新**：系统设置支持动态更新（部分配置需重启）

## 工作流程

### 标准开发流程
1. **需求理解**：分析用户需求，确定功能范围
2. **架构设计**：设计数据模型、API 接口、业务逻辑
3. **代码实现**：按照分层架构实现代码
4. **测试验证**：编写和执行测试
5. **文档更新**：更新 README、Swagger、代码注释
6. **代码审查**：检查代码质量和规范
7. **集成验证**：验证前后端集成和功能完整性

### 快速修复流程
1. **问题定位**：快速定位问题根源
2. **影响评估**：评估修复的影响范围
3. **快速修复**：应用最小化修复
4. **验证测试**：快速验证修复效果
5. **文档同步**：更新相关文档

## 使用场景示例

### 场景1：添加新的路由策略类型
1. 在 `PolicyTemplate` 中添加新模板定义（`service/policy_template_service.go`）
2. 实现新的 `PolicyExecutor`（`service/policy_executor.go`）
3. 在 `PolicyService` 中注册新执行器
4. 更新前端策略模板列表和参数表单
5. 更新文档说明新策略类型

### 场景2：添加新的 LLM 驱动支持
1. 在 `LLMResource` 的 `Driver` 字段中添加新驱动类型（如 "anthropic"、"google"）
2. 在 `backend/internal/client/` 中实现新驱动的客户端
3. 在 `handler/provider.go` 中更新驱动验证逻辑
4. 更新前端驱动选择列表（如需要）
5. 更新文档说明新驱动

### 场景2.1：实现批量导入功能
1. 后端：在 `handler/provider.go` 中实现导入接口
   - 支持 Excel 文件解析（使用 `excelize` 库）
   - 支持 JSON 数组解析
   - 实现重复检测逻辑（基于 type、model、base_url、api_key）
   - 返回导入结果（成功、失败、重复数量及详情）
2. 前端：在 `views/LLMResources.vue` 中添加导入功能
   - 文件上传组件（支持 Excel 和 JSON）
   - 导入结果展示（成功/失败/重复列表）
   - 模板下载功能

### 场景2.2：实现资源测试功能
1. 后端：在 `handler/provider.go` 中实现测试接口
   - 根据资源类型（chat、embedding）调用相应测试方法
   - 设置测试超时（30秒）
   - 返回测试结果（响应时间、模型信息、Token 使用等）
   - 更新资源的 `test_status` 字段
2. 前端：在 `views/LLMResources.vue` 中添加测试按钮
   - 仅对 `active` 状态的资源显示测试按钮
   - 显示测试结果弹窗
   - 更新资源列表中的测试状态

### 场景3：修复认证相关问题
1. 检查 `middleware/auth.go` 中的认证逻辑
2. 检查 `router/router.go` 中的认证中间件应用
3. 检查 `config.yaml` 中的认证开关配置
4. 验证 Token 验证逻辑和错误处理
5. 测试认证开启和关闭两种场景

### 场景4：更新系统设置功能
1. 在 `backend/internal/config/config.go` 中添加新配置项
2. 在 `backend/internal/service/settings_service.go` 中添加更新逻辑
3. 在 `backend/internal/handler/settings.go` 中添加 API 端点
4. 在前端 `views/Settings.vue` 中添加 UI
5. 添加国际化文本（中英文）
6. 更新配置文档和默认值

### 场景5：添加国际化文本
1. 在 `frontend/src/locales/zh/index.js` 中添加中文文本
2. 在 `frontend/src/locales/en/index.js` 中添加英文文本
3. 在组件中使用 `i18n.global.t('key')` 引用文本
4. 确保 Element Plus 组件语言包正确设置（在 `main.js` 中）

### 场景6：添加新的 API 端点
1. 后端：
   - 在 `backend/internal/service/` 中添加业务逻辑
   - 在 `backend/internal/handler/` 中添加处理器
   - 在 `backend/internal/router/router.go` 中注册路由（考虑认证开关）
   - 添加 Swagger 注释
2. 前端：
   - 在 `frontend/src/api/index.js` 中添加 API 方法
   - 在相应的 Vue 组件中调用 API
   - 添加错误处理和加载状态
   - 添加国际化文本

## 注意事项

### 关键约束
1. **认证开关**：所有资源类 API 需要根据 `cfg.Security.Auth.Enabled` 条件应用认证
2. **Token 前缀**：请求端 Token 必须以 "ling-" 开头
3. **策略绑定**：只有 "ling-" 开头的 Token 才能绑定策略
4. **配置持久化**：系统设置更新需要保存到 `backend/configs/config.yaml`
5. **数据库迁移**：新增模型需要在 `backend/internal/storage/gorm_storage.go` 的 `AutoMigrate` 中添加
6. **驱动架构**：LLM 资源使用 Driver 字段，目前仅支持 "openai"，添加新驱动需要实现相应客户端
7. **国际化要求**：所有前端用户可见文本必须使用 i18n，不能硬编码中文或英文
8. **API 路径一致性**：前后端 API 路径必须一致，前端使用 `frontend/src/utils/constants.js` 中的 `API_BASE_URL`
9. **批量导入格式**：Excel 导入模板必须包含：Name、Type、Driver、Model、BaseURL、APIKey、Status
10. **资源测试限制**：仅支持测试 `active` 状态的资源，仅支持 `chat` 和 `embedding` 类型

### 常见陷阱
1. **循环依赖**：避免 `service` 和 `handler` 之间的循环导入
2. **认证遗漏**：新增 API 时忘记应用认证中间件
3. **配置默认值**：新增配置项时忘记设置默认值
4. **Swagger 注释**：忘记更新 API 注释导致文档不准确
5. **前端 API 路径**：前后端 API 路径不一致
6. **国际化遗漏**：前端添加新功能时忘记添加国际化文本
7. **数据库迁移遗漏**：新增模型字段时忘记更新 `AutoMigrate`
8. **敏感信息泄露**：忘记使用 `json:"-"` 隐藏敏感字段
9. **批量导入验证**：忘记验证导入数据的格式和必填字段
10. **资源测试超时**：忘记设置测试超时时间，导致长时间等待
11. **Excel 导入格式**：Excel 文件格式不正确或字段名不匹配
12. **Token 前缀验证**：忘记验证 Token 前缀，导致非 "ling-" 前缀的 Token 被创建

## API 路径与常量配置

### 后端 API 路径
- **管理 API**：`/api/v1/*`（需要认证，除非认证开关关闭）
- **OpenAI 兼容 API**：`/llm/v1/*`（使用 Token 认证）
- **健康检查**：`/health`（无需认证）

### 前端常量配置
所有前端常量定义在 `frontend/src/utils/constants.js` 中：
- `API_BASE_URL`：API 基础 URL（默认：`http://localhost:8080/api/v1`）
- `API_TIMEOUT`：API 请求超时时间（默认：30000ms）
- `STORAGE_KEYS`：LocalStorage 键名（TOKEN、USER_INFO、LOCALE 等）
- `MESSAGE_DURATION`：消息提示持续时间

### API 路径映射规则
- 前端 API 方法名对应后端路径：
  - `api.getLLMResources()` → `GET /api/v1/llm-resources`
  - `api.createLLMResource()` → `POST /api/v1/llm-resources`
  - `api.updateLLMResource(id)` → `PUT /api/v1/llm-resources/:id`
  - `api.deleteLLMResource(id)` → `DELETE /api/v1/llm-resources/:id`
- OpenAI 兼容 API（不需要通过前端 API 客户端）：
  - `POST /llm/v1/chat/completions`
  - `GET /llm/v1/models`

## 工具与命令

### 开发命令
```bash
# 后端开发
cd backend
go run cmd/main.go

# 前端开发
cd frontend
npm install          # 首次运行需要安装依赖
npm run dev          # 启动开发服务器（默认 http://localhost:3000）

# 构建
make build           # 构建后端和前端
make build-backend   # 仅构建后端
make build-frontend  # 仅构建前端

# 测试
cd backend
go test ./...        # 运行所有测试
go test -cover ./... # 运行测试并显示覆盖率

# Swagger 文档生成
cd backend
swag init            # 生成 Swagger 文档（需要安装 swag）
```

### 服务管理
```bash
# 启动服务
make start

# 停止服务
make stop

# 重启服务
make restart
```

## 开发最佳实践

### 代码质量
- **错误处理**：所有可能出错的地方都要有错误处理，使用统一的错误返回格式
- **日志记录**：关键操作和错误都要记录日志，使用适当的日志级别
- **代码注释**：公共函数和复杂逻辑要有清晰的注释
- **类型安全**：Go 代码充分利用类型系统，前端使用 TypeScript（如迁移）或 JSDoc

### 性能优化
- **数据库查询**：避免 N+1 查询，使用预加载和批量查询
- **前端优化**：使用 Vue 3 的响应式系统，避免不必要的重渲染
- **API 调用**：合理使用缓存，避免重复请求
- **批量操作**：支持批量操作的功能优先使用批量接口

### 安全性
- **敏感信息**：密码、API Key 等敏感信息使用 `json:"-"` 隐藏
- **输入验证**：所有用户输入都要验证，防止注入攻击
- **认证授权**：正确应用认证中间件，验证用户权限
- **HTTPS**：生产环境必须使用 HTTPS

### 可维护性
- **模块化**：功能模块化，保持单一职责
- **配置化**：可配置项放在配置文件中，避免硬编码
- **文档化**：及时更新 README、API 文档和代码注释
- **测试**：关键功能要有单元测试和集成测试

## 总结

LingProxy Dev Assistant 是一个专门为 LingProxy 项目定制的智能开发助手，深度理解项目架构和业务逻辑，能够：
- 快速理解项目结构和代码组织
- 按照项目规范生成代码（包括国际化支持）
- 识别和修复常见问题
- 维护文档和 API 文档
- 支持功能开发和系统优化
- 协助实现批量导入、资源测试等高级功能
- 确保前后端 API 一致性和国际化完整性

通过遵循项目的最佳实践和规范，Agent 能够高效地协助开发工作，确保代码质量、安全性和项目一致性。在开发新功能时，务必考虑国际化、错误处理、日志记录和文档更新等各个方面。
