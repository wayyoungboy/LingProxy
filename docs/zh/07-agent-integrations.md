# Claude Code 与 Codex 集成

LingProxy 在 `/llm/v1` 暴露 OpenAI 兼容 API。适合接入允许自定义 OpenAI-compatible base URL 的智能体和开发工具。

## Codex

Codex 可以通过用户级 Codex 配置把 OpenAI provider 请求路由到 LingProxy。

1. 在 LingProxy 管理后台创建请求端 API Key。
2. 确认该 Key 绑定的策略可以路由你要使用的模型。
3. 在 `~/.codex/config.toml` 中加入：

```toml
openai_base_url = "http://localhost:8080/llm/v1"
```

4. 启动 Codex 前导出 LingProxy Key：

```bash
export OPENAI_API_KEY="ling-xxxxxxxxxxxxx"
codex
```

不要把这个 base URL 放在 Codex 项目级配置里。网关地址和 API Key 应放在用户级 Codex 配置/环境变量中，项目级规则只用于仓库行为。

## Claude Code

Claude Code 默认不是 OpenAI 兼容模型客户端。因此在 LingProxy 当前仅提供 OpenAI driver 的情况下，不应把 LingProxy 写成 Claude Code 的原生模型网关。

当前可用的 Claude Code 工作流：

- 使用 Claude Code 开发、维护、审查 LingProxy 本身。
- Claude Code 在仓库中工作时，让 OpenAI 兼容 SDK、测试和本地工具通过 LingProxy。
- 如果希望 Claude 模型流量也由 LingProxy 托管，需要先为 LingProxy 增加 Anthropic/Claude driver。

## 快速验证

配置 Codex 或其他 OpenAI 兼容客户端后，先验证网关：

```bash
curl http://localhost:8080/llm/v1/models \
  -H "Authorization: Bearer $OPENAI_API_KEY"
```

再通过相同 base URL 发送聊天补全：

```bash
curl http://localhost:8080/llm/v1/chat/completions \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role": "user", "content": "Say hello from LingProxy."}]
  }'
```
