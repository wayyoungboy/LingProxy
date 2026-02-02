# 安全指南

## 敏感信息处理

### 配置文件安全

1. **config.yaml 已添加到 .gitignore**
   - 实际的配置文件 `configs/config.yaml` 不会被提交到 Git
   - 使用 `configs/config.yaml.example` 作为模板

2. **默认密码**
   - 默认密码已从硬编码的 `admin123` 改为 `CHANGE_ME`
   - 首次启动前必须修改配置文件中的密码
   - 代码中不再硬编码默认密码

### 环境变量

使用环境变量存储敏感信息（推荐）：

```bash
# 设置管理员密码
export LINGPROXY_ADMIN_PASSWORD=your_strong_password

# 设置 API Key
export LINGPROXY_API_KEY=your_api_key
```

### 代码中的敏感信息检查

✅ **已修复的问题**：
- [x] 配置文件中的硬编码密码
- [x] 代码中的默认密码
- [x] 前端代码中的默认密码
- [x] 示例代码中的硬编码 API Key
- [x] 文档中的示例密码

### 安全最佳实践

1. **生产环境部署**
   - 使用强密码（至少12位，包含大小写字母、数字和特殊字符）
   - 启用认证（`security.auth.enabled: true`）
   - 使用 HTTPS
   - 定期更新密码和 API Key

2. **配置文件管理**
   - 不要将包含真实密码的 `config.yaml` 提交到版本控制
   - 使用环境变量或密钥管理服务
   - 限制配置文件的访问权限（chmod 600）

3. **API Key 管理**
   - Token 以 `ling-` 开头，便于识别和管理
   - 定期轮换 API Key
   - 不要在前端代码中硬编码 API Key

4. **日志安全**
   - 不要在日志中输出密码或 API Key
   - 敏感字段已使用 `json:"-"` 标记，不会序列化

5. **数据库安全**
   - SQLite 数据库文件已添加到 .gitignore
   - 生产环境建议使用 MySQL 或 PostgreSQL
   - 数据库连接字符串不要硬编码在代码中

### 检查清单

部署前请确认：
- [ ] 已修改 `config.yaml` 中的默认密码
- [ ] 已设置强密码（至少12位）
- [ ] 已启用认证（`security.auth.enabled: true`）
- [ ] 配置文件权限已限制（chmod 600）
- [ ] 数据库文件不在版本控制中
- [ ] 日志中不包含敏感信息
- [ ] 使用 HTTPS（生产环境）

### 报告安全问题

如果发现安全问题，请通过以下方式报告：
- GitHub Issues（私有仓库）
- 或直接联系维护者
