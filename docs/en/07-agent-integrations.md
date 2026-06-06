# Claude Code and Codex Integrations

LingProxy exposes an OpenAI-compatible API at `/llm/v1`. Use it with agents and developer tools that allow an OpenAI-compatible base URL.

## Codex

Codex can route OpenAI-provider requests through LingProxy from the user-level Codex config.

1. Create a LingProxy request-side API key in the admin UI.
2. Make sure the key is bound to a policy that can route the models you want to use.
3. Add this to `~/.codex/config.toml`:

```toml
openai_base_url = "http://localhost:8080/llm/v1"
```

4. Export the LingProxy key before starting Codex:

```bash
export OPENAI_API_KEY="ling-xxxxxxxxxxxxx"
codex
```

Codex project config files should not be used for this base URL override. Keep the gateway URL and API key in the user-level Codex config/environment, then use project-level rules only for repository behavior.

## Claude Code

Claude Code is not an OpenAI-compatible model client by default, so LingProxy should not be presented as a direct Claude Code model gateway while LingProxy only ships the OpenAI driver.

Useful Claude Code workflows today:

- Use Claude Code to edit, operate, or review LingProxy itself.
- Route OpenAI-compatible SDKs, tests, and local tools through LingProxy while Claude Code works in the repo.
- Add an Anthropic/Claude driver to LingProxy first if you want Claude model traffic to be managed by the gateway.

## Quick Validation

After configuring Codex or another OpenAI-compatible client, verify the gateway:

```bash
curl http://localhost:8080/llm/v1/models \
  -H "Authorization: Bearer $OPENAI_API_KEY"
```

Then run a chat completion through the same base URL:

```bash
curl http://localhost:8080/llm/v1/chat/completions \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role": "user", "content": "Say hello from LingProxy."}]
  }'
```
