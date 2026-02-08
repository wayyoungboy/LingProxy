# Configuration Guide

## Configuration File

The main configuration file is located at `backend/configs/config.yaml`. Copy `config.yaml.example` to create your configuration file.

## Configuration Structure

### Application Configuration

```yaml
app:
  environment: "development"  # development, staging, production
  port: 8080                  # Server port
  # name and version have defaults, usually not needed
  # host has default "0.0.0.0", usually not needed
```

### Storage Configuration

```yaml
storage:
  type: "gorm"  # memory or gorm
  gorm:
    driver: "sqlite"  # sqlite, mysql
    dsn: "lingproxy.db"
    # MySQL example:
    # driver: "mysql"
    # dsn: "username:password@tcp(localhost:3306)/lingproxy?charset=utf8mb4&parseTime=True&loc=Local"
```

**Storage Types:**
- `memory`: In-memory storage (for development/testing)
- `gorm`: Database storage (SQLite, MySQL, PostgreSQL)

### Log Configuration

```yaml
log:
  level: "info"      # debug, info, warn, error, fatal
  format: "json"     # text, json
  output: "both"     # stdout, file, both (recommended: both)
  file_path: "./logs/lingproxy.log"
  # max_size, max_backups, max_age, compress have defaults
```

**Log Levels:**
- `debug`: Detailed debugging information
- `info`: General informational messages
- `warn`: Warning messages
- `error`: Error messages
- `fatal`: Fatal errors

**Output Modes:**
- `stdout`: Console output only
- `file`: File output only
- `both`: Both console and file (recommended)

### Security Configuration

```yaml
security:
  auth:
    enabled: true  # Enable/disable authentication globally
  # cors, rate_limit, jwt can be configured via admin interface
```

**Authentication:**
- When `enabled: true`: All APIs (except login) require authentication
- When `enabled: false`: All APIs (except login) are accessible without authentication

### Provider Configuration

```yaml
provider:
  timeout: "30s"        # Request timeout
  max_retries: 3        # Maximum retry count (0 = disabled)
  retry_delay: "1s"     # Base retry delay (actual delay increases exponentially)
  max_idle_conns: 100   # Maximum idle connections
  max_conns_per_host: 100  # Maximum connections per host
  idle_conn_timeout: "90s"  # Idle connection timeout
```

**Retry Configuration:**
- `max_retries`: Maximum number of retry attempts for failed requests. Set to `0` to disable retry.
- `retry_delay`: Base delay between retry attempts. The actual delay increases exponentially: `delay = retry_delay × attempt_number`.
- Retry applies to network errors, timeouts, and 5xx server errors.
- Does not retry 4xx client errors (except 429 rate limit), authentication errors, or context cancellations.
- Configurable via admin interface: Settings → Provider Settings.

## Environment Variables

You can override configuration values using environment variables with the `LINGPROXY_` prefix:

```bash
# Example: Override port
export LINGPROXY_APP_PORT=9000

# Example: Override database DSN
export LINGPROXY_STORAGE_GORM_DSN="mysql://user:pass@localhost/db"
```

## Default Values

All configuration options have sensible defaults. See `backend/internal/config/config.go` for complete default values:

- `app.name`: "LingProxy"
- `app.version`: "1.0.0"
- `app.host`: "0.0.0.0"
- `app.port`: 8080
- `app.environment`: "development"
- `storage.type`: "memory"
- `log.level`: "info"
- `log.format`: "json"
- `log.output`: "both"
- `security.auth.enabled`: true
- `provider.timeout`: "30s"
- `provider.max_retries`: 3
- `provider.retry_delay`: "1s"

## Admin Credentials

Default admin credentials (set on first startup):
- Username: `admin`
- Password: `admin123`

**Important**: Change the default password after first login!

## Configuration via Admin Interface

Many settings can be configured via the admin dashboard:
- **System Settings**: Basic settings, cache, rate limiting, security, logging, load balancing, provider retry
- **LLM Resources**: Add, edit, and manage AI service resources
- **Models**: Configure model details, pricing, and limits
- **Policies**: Create and manage routing policies
- **API Keys**: Create and manage request-side API keys with policy binding

Changes made via the admin interface are saved to the configuration file.

### Provider Retry Settings

You can configure retry behavior for resource requests through the admin interface:

1. Navigate to **Settings** → **Provider Settings**
2. Configure:
   - **Request Timeout**: Maximum time to wait for a response (seconds)
   - **Max Retries**: Maximum number of retry attempts (0-10, 0 = disabled)
   - **Retry Delay**: Base delay between retries (seconds, actual delay increases exponentially)

These settings apply immediately to all new requests without requiring a restart.

## Production Configuration

For production deployments:

1. Set `app.environment` to `"production"`
2. Use `gorm` storage type with a production database (MySQL/PostgreSQL)
3. Enable authentication (`security.auth.enabled: true`)
4. Configure proper log rotation
5. Use strong passwords and API keys
6. Configure CORS appropriately for your domain

## Troubleshooting

### Configuration File Not Found

If the configuration file is not found, the application will use default values. Check the log output for configuration loading messages.

### Invalid Configuration

If configuration validation fails, check the error message and verify:
- YAML syntax is correct
- Required fields are present
- Values are within valid ranges

### Configuration Not Applied

After modifying the configuration file:
1. Restart the application
2. Check logs for configuration loading messages
3. Verify changes via the admin interface
