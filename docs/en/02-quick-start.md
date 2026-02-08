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

### Using Docker Compose (Recommended)

```bash
# From project root
docker-compose -f docker/docker-compose.yml up -d

# View logs
docker-compose -f docker/docker-compose.yml logs -f

# Stop services
docker-compose -f docker/docker-compose.yml down
```

For more details about Docker deployment, see the Docker section above or check the project's `docker/README.md` file.

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
# Using curl
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

## Next Steps

- Read the [Configuration Guide](03-configuration.md) for detailed configuration options
- Check the [API Documentation](04-api-reference.md) for API usage
- Review the [Architecture Guide](05-architecture.md) for system design details
- See the [Development Guide](06-development.md) for contributing
