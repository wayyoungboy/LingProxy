# Quick Start Guide

## Prerequisites

- **Backend**: Go 1.21 or higher
- **Frontend**: Node.js 18+ and npm/yarn
- **Database**: SQLite (included) or MySQL/PostgreSQL (optional)

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/your-org/lingproxy.git
cd lingproxy
```

### 2. Backend Setup

#### Install Dependencies

```bash
cd backend
go mod tidy
```

#### Configure the Application

```bash
# Copy the example configuration file
cp configs/config.yaml.example configs/config.yaml

# Edit the configuration file
# The default admin credentials are:
# Username: admin
# Password: admin123
```

#### Run the Backend

```bash
# Development mode
go run cmd/main.go

# Or build and run
go build -o lingproxy cmd/main.go
./lingproxy
```

The backend will start at `http://localhost:8080`

### 3. Frontend Setup

#### Install Dependencies

```bash
cd frontend
npm install
```

#### Run the Frontend

```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

### 4. Access the Dashboard

1. Open your browser and navigate to `http://localhost:3000`
2. Login with default credentials:
   - Username: `admin`
   - Password: `admin123`
3. **Important**: Change the default password after first login!

## Docker Deployment

### Prerequisites

- **Docker**: Docker 20.10+ and Docker Compose 2.0+
- **Configuration**: Ensure `backend/configs/config.yaml` exists and is properly configured

### Quick Start with Makefile (Recommended)

The easiest way to start all services:

```bash
# From project root - one command to start everything
make docker-compose-up
```

This command will:
1. ✅ Check if configuration file exists (create from example if needed)
2. ✅ Create necessary directories (logs, run)
3. ✅ Start SeekDB and LingProxy services (build if needed)
4. ✅ Wait for SeekDB to be ready
5. ✅ Create the database automatically
6. ✅ Display access URLs

**Access URLs** (after startup):
- **Backend API**: http://localhost:8080/api/v1
- **Health Check**: http://localhost:8080/api/v1/health

**Note**: The Docker deployment only includes the backend service. For frontend development, run it separately:
```bash
cd frontend
npm run dev
# Frontend will be available at http://localhost:3000
```

### Manual Docker Compose Commands

If you prefer to use Docker Compose directly:

```bash
# 1. Prepare configuration file
cp backend/configs/config.yaml.example backend/configs/config.yaml
# Edit backend/configs/config.yaml and configure SeekDB connection:
# storage:
#   type: "gorm"
#   gorm:
#     driver: "mysql"
#     dsn: "root:@tcp(seekdb:2881)/lingproxy?charset=utf8mb4&parseTime=True&loc=Local"

# 2. Start services
docker-compose -f docker/docker-compose.yml up -d --build

# 3. Database will be created automatically on backend startup
# No manual database creation needed

# 4. View logs
docker-compose -f docker/docker-compose.yml logs -f

# 5. Stop services
docker-compose -f docker/docker-compose.yml down
```

### Other Useful Makefile Commands

```bash
# Check service status
make docker-compose-ps

# View logs
make docker-compose-logs

# Stop services
make docker-compose-down

# Restart services
make docker-compose-restart

# Initialize database only
make docker-compose-init-db
```

### Service Architecture

The Docker deployment uses a **frontend-backend separation** architecture:

- **SeekDB**: MySQL-compatible database (ports 2881, 2886)
  - Data stored in Docker volume `seekdb-data`
  - Health check ensures service is ready before backend starts

- **LingProxy Backend**: Backend API service (port 8080)
  - Pure backend service, no frontend included
  - Automatically creates database on startup if not exists
  - Uses Docker-specific configuration (`config.yaml.docker`)

- **Frontend**: Run separately for development
  - Use `npm run dev` in the `frontend` directory
  - Runs on port 3000 with Vite dev server
  - API requests are proxied to `http://localhost:8080` via Vite proxy

### Troubleshooting

**Service won't start:**
```bash
# Check service status
make docker-compose-ps

# View logs
make docker-compose-logs

# Check if ports are in use
lsof -i :8080
lsof -i :2881
```

**Database connection issues:**
```bash
# Verify SeekDB is running
docker exec seekdb mysql -h127.0.0.1 -uroot -P2881 -e "SHOW DATABASES;"

# Recreate database
make docker-compose-init-db
```

**Configuration issues:**
```bash
# Check configuration file
cat backend/configs/config.yaml

# Verify SeekDB connection string
grep -A 3 "storage:" backend/configs/config.yaml
```

For more details about Docker deployment, see the [Configuration Guide](03-configuration.md) and [Development Guide](06-development.md).

## First Steps

### 1. Configure LLM Resources

1. Navigate to **LLM Resources** in the dashboard
2. Click **Add Resource**
3. Fill in the resource information:
   - Name: A descriptive name
   - Type: Model category (chat, image, embedding, etc.)
   - Driver: Currently only "openai" is supported
   - Model: Model identifier (e.g., gpt-4, gpt-3.5-turbo)
   - Base URL: API endpoint URL
   - API Key: Your API key

### 2. Create an API Key

1. Navigate to **API Key Management** in the dashboard
2. Click **Create API Key**
3. Fill in the API key information:
   - Name: API key name/description
   - Policy: Select a routing policy (optional)

### 3. Test the API

```bash
# Using curl (backend runs on port 8080)
curl -X POST http://localhost:8080/llm/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  }'
```

**Note**: If you're running the frontend locally (`npm run dev`), API requests from the frontend will be automatically proxied to the backend.

## Next Steps

- Read the [Configuration Guide](03-configuration.md) for detailed configuration options
- Check the [API Documentation](04-api-reference.md) for API usage
- Review the [Architecture Guide](05-architecture.md) for system design details
- See the [Development Guide](06-development.md) for contributing
