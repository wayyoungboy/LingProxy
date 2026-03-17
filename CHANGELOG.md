# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security

## [1.0.0] - 2025-03-17

### Added

- **Core Features**
  - Unified OpenAI-compatible API interface for seamless AI service integration
  - Full Server-Sent Events (SSE) streaming support for chat completions
  - Intelligent round-robin load balancing across multiple resources
  - Configurable automatic retry with exponential backoff
  - Circuit breaking to prevent cascading failures
  - Complete request chain tracing and logging

- **Security & Authentication**
  - Flexible global authentication toggle
  - Admin login with username/password (password hash storage)
  - API key management with policy association
  - CORS support with flexible configuration
  - Encrypted storage for API keys and passwords

- **Management Features**
  - Modern Vue 3 + Element Plus admin dashboard
  - Full internationalization (Chinese and English)
  - Policy management with 8 built-in routing templates:
    - Random, Round-robin, Weighted, Model-match
    - Regex-match, Regex-model-match, Priority, Failover
  - LLM resource management with driver-based architecture
  - Model categories: chat, image, embedding, rerank, audio, video
  - Batch import/export via Excel templates or JSON
  - Resource connectivity testing
  - Usage statistics with time range and resource filtering
  - Dynamic system configuration management
  - Real-time system monitoring (CPU, memory, uptime)
  - Log management with filtering and search

- **Architecture**
  - Frontend-backend separation (Vue 3 + Go)
  - Dual storage: memory (dev) and SQLite (production)
  - RESTful API with Swagger documentation
  - Client libraries: Python, JavaScript, Go

- **DevOps**
  - Multi-platform binary builds (Linux, macOS, Windows)
  - Docker support with docker-compose
  - CI/CD pipeline with GitHub Actions
  - Automated releases with GitHub releases

### Security

- Encrypted password storage
- API key encryption at rest
- Configurable authentication requirements

---

## Version History

| Version | Date | Description |
|---------|------|-------------|
| 1.0.0 | 2025-03-17 | Initial release |

[Unreleased]: https://github.com/lingproxy/lingproxy/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/lingproxy/lingproxy/releases/tag/v1.0.0