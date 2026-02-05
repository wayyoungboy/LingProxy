# Introduction

<div align="center">

<img src="../../assets/lingproxy-logo-full.svg" alt="LingProxy Logo" width="300">

</div>

## What is LingProxy?

LingProxy is a high-performance AI API gateway designed for managing and proxying API calls to various AI service providers. It provides OpenAI-compatible interfaces, intelligent load balancing, request routing, and comprehensive management features.

## Key Features

### 🚀 Core Capabilities
- **OpenAI-Compatible API**: Seamless integration with OpenAI SDK and compatible clients
- **Streaming Support**: Full support for Server-Sent Events (SSE) streaming responses for chat completions
- **Intelligent Load Balancing**: Multiple routing strategies (round-robin, random, weighted, etc.)
- **Request Routing**: Policy-based request routing with model matching and failover
- **Request Logging**: Complete request chain tracing and logging

### 🔐 Security & Authentication
- **Flexible Authentication**: Global authentication toggle, configurable per endpoint
- **Admin Management**: Username/password authentication with secure password hashing
- **Token Management**: Request-side token management with policy association
- **API Key Authentication**: Secure API key-based authentication
- **CORS Support**: Configurable cross-origin resource sharing

### 📊 Management Features
- **Admin Dashboard**: Modern web-based management interface built with Vue 3 + Element Plus
- **Internationalization (i18n)**: Full support for Chinese and English language switching in the frontend interface
- **LLM Resource Management**: Driver-based architecture (currently supports OpenAI driver)
- **Model Management**: Flexible model configuration with pricing and usage limits
- **Policy Management**: Built-in routing policy templates and custom policy instances, random selection policy supports LLM resource pool configuration
- **Request Management**: Complete request logging and tracking functionality
- **Usage Statistics**: Detailed usage statistics grouped by LLM resources, including token usage, request count, success rate, and more
- **System Settings**: Dynamic configuration management
- **System Monitoring**: Real-time system information and statistics
- **Log Management**: View and manage system logs with filtering and search capabilities

## Use Cases

- **Multi-Provider AI Services**: Manage multiple AI service providers through a unified interface
- **Load Balancing**: Distribute requests across multiple resources for better performance
- **Request Routing**: Route requests to different providers based on policies
- **API Gateway**: Act as a gateway for AI service APIs with authentication and logging
- **Development & Testing**: Use as a proxy for development and testing environments

## Architecture

LingProxy follows a modern microservices architecture:

- **Frontend**: Vue 3 + Element Plus + Vite
- **Backend**: Go + Gin + GORM
- **Storage**: SQLite (default) / MySQL / PostgreSQL
- **API**: RESTful API + OpenAI-compatible API

## License

[Add your license information here]

## Support

For issues, questions, or contributions, please visit:
- GitHub Issues: [Add your GitHub repository URL]
- Documentation: [Add documentation URL]
