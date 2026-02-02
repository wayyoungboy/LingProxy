# LingProxy 开发 Agent 设计

## 基本信息
- **名称**：LingProxy Dev Assistant
- **版本**：2.0.0
- **项目类型**：AI API 网关（OpenAI 兼容）
- **技术栈**：Go + Gin + GORM (后端)、Vue 3 + Element Plus + Vite (前端)、SQLite/MySQL/PostgreSQL (数据库)

## 项目架构理解

### 核心模块
1. **路由策略系统**：支持随机、轮询、加权、模型匹配、正则匹配、优先级、故障转移等策略
2. **Token 管理**：请求端 Token 管理，支持策略绑定（"ling-" 前缀 Token）
3. **LLM 资源管理**：多提供商支持（OpenAI、Zai、Anthropic、Google、Azure 等），多模型类别（chat、image、embedding、rerank、audio、video）
4. **认证系统**：全局认证开关，管理员登录（用户名/密码），Token 认证
5. **系统设置**：动态配置管理（基础、缓存、限流、安全、日志、负载均衡、熔断器）
6. **OpenAI 兼容 API**：`/llm/v1/*` 路径前缀

### 代码组织
```
internal/
├── handler/        # HTTP 处理器（RESTful API）
├── service/        # 业务逻辑层
├── storage/        # 数据存储层（接口 + 内存/GORM 实现）
├── middleware/     # 中间件（认证、CORS、日志、限流、代理）
├── config/         # 配置管理（Viper）
├── router/         # 路由配置
└── client/         # AI 服务客户端（OpenAI、Embedding）
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
2. **数据模型设计**：在 `internal/storage/models.go` 中定义/更新模型
3. **存储层实现**：在 `internal/storage/storage.go` 中定义接口，在 `memory_storage.go` 和 `gorm_storage.go` 中实现
4. **业务逻辑层**：在 `internal/service/` 中实现业务逻辑
5. **HTTP 处理器**：在 `internal/handler/` 中实现 RESTful API
6. **路由注册**：在 `internal/router/router.go` 中注册路由（考虑认证开关）
7. **前端集成**：在 `frontend/src/` 中实现 UI 和 API 调用
8. **Swagger 文档**：更新 API 注释和文档

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

**系统设置功能**：
- 配置读取：使用 Viper 读取 `configs/config.yaml`
- 配置更新：动态更新配置并持久化到文件
- 配置验证：验证配置项的有效性
- 重启提示：识别需要重启的配置项

### 3. 代码规范与最佳实践

#### 3.1 Go 后端规范
- **错误处理**：使用 `storage.ErrNotFound` 等标准错误
- **日志记录**：使用 `pkg/logger` 包，记录关键操作和错误
- **认证检查**：根据 `cfg.Security.Auth.Enabled` 条件应用认证中间件
- **ID 生成**：使用 `generateID()` 生成唯一标识
- **时间戳**：统一使用 `time.Now()` 设置 `CreatedAt` 和 `UpdatedAt`
- **JSON 序列化**：敏感字段使用 `json:"-"` 隐藏

#### 3.2 前端规范
- **API 调用**：统一使用 `frontend/src/api/index.js` 中的 API 客户端
- **组件组织**：页面组件放在 `views/`，通用组件放在 `components/`
- **状态管理**：使用 Vue 3 Composition API
- **错误处理**：使用 Element Plus 的 `ElMessage` 显示错误
- **表单验证**：使用 Element Plus 的表单验证规则

#### 3.3 数据库规范
- **迁移**：在 `gorm_storage.go` 的 `Init()` 中使用 `AutoMigrate`
- **查询优化**：使用 GORM 的预加载和索引
- **事务处理**：需要时使用 GORM 事务

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

### 场景2：添加新的 LLM 提供商支持
1. 在 `LLMResource` 的 `Provider` 字段中添加新提供商
2. 更新前端提供商选择列表和 BaseURL 映射
3. 如需要特殊处理，在 `handler/provider.go` 中添加逻辑
4. 更新文档说明新提供商

### 场景3：修复认证相关问题
1. 检查 `middleware/auth.go` 中的认证逻辑
2. 检查 `router/router.go` 中的认证中间件应用
3. 检查 `config.yaml` 中的认证开关配置
4. 验证 Token 验证逻辑和错误处理
5. 测试认证开启和关闭两种场景

### 场景4：更新系统设置功能
1. 在 `config/config.go` 中添加新配置项
2. 在 `service/settings_service.go` 中添加更新逻辑
3. 在 `handler/settings.go` 中添加 API 端点
4. 在前端 `Settings.vue` 中添加 UI
5. 更新配置文档和默认值

## 注意事项

### 关键约束
1. **认证开关**：所有资源类 API 需要根据 `cfg.Security.Auth.Enabled` 条件应用认证
2. **Token 前缀**：请求端 Token 必须以 "ling-" 开头
3. **策略绑定**：只有 "ling-" 开头的 Token 才能绑定策略
4. **配置持久化**：系统设置更新需要保存到 `configs/config.yaml`
5. **数据库迁移**：新增模型需要在 `gorm_storage.go` 的 `AutoMigrate` 中添加

### 常见陷阱
1. **循环依赖**：避免 `service` 和 `handler` 之间的循环导入
2. **认证遗漏**：新增 API 时忘记应用认证中间件
3. **配置默认值**：新增配置项时忘记设置默认值
4. **Swagger 注释**：忘记更新 API 注释导致文档不准确
5. **前端 API 路径**：前后端 API 路径不一致

## 工具与命令

### 开发命令
```bash
# 后端开发
go run cmd/main.go

# 前端开发
cd frontend && npm run dev

# 构建
make build

# 测试
go test ./...

# Swagger 文档生成
swag init
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

## 总结

LingProxy Dev Assistant 是一个专门为 LingProxy 项目定制的智能开发助手，深度理解项目架构和业务逻辑，能够：
- 快速理解项目结构和代码组织
- 按照项目规范生成代码
- 识别和修复常见问题
- 维护文档和 API 文档
- 支持功能开发和系统优化

通过遵循项目的最佳实践和规范，Agent 能够高效地协助开发工作，确保代码质量和项目一致性。
