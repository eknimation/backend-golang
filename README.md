# Backend Go API

A robust Go backend service with MongoDB integration, JWT authentication, and comprehensive API documentation using Swagger.

## 🚀 Features

- **RESTful API** with Echo framework
- **MongoDB** integration with proper indexing
- **JWT Authentication** with middleware protection
- **Swagger Documentation** auto-generated from code annotations
- **Docker & Docker Compose** for easy deployment
- **Health Check** endpoints for monitoring
- **Graceful Shutdown** with proper resource cleanup
- **Request/Response Logging** with structured logging
- **Input Validation** with custom error formatting
- **Password Hashing** with bcrypt
- **Environment-based Configuration**

## 📋 Prerequisites

- Docker & Docker Compose
- Go 1.24+ (for local development)
- Make (optional, for convenience commands)

## 🔧 Quick Start

### Using Docker (Recommended)

1. **Clone and setup:**
   ```bash
   git clone <repository-url>
   cd backend-golang
   ```

2. **Run the setup script:**
   ```bash
   ./setup-dev.sh
   ```

3. **Or manually start services:**
   ```bash
   make up
   # or
   docker-compose up -d
   ```

4. **Check service status:**
   ```bash
   make status
   ```

### Manual Setup

1. **Start MongoDB:**
   ```bash
   docker-compose up -d mongodb
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Generate Swagger docs:**
   ```bash
   swag init -g cmd/api/main.go -o docs
   ```

4. **Run the application:**
   ```bash
   go run cmd/api/main.go
   ```

## 🐳 Docker Commands

### Basic Operations
```bash
make up          # Start all services
make down        # Stop all services
make logs        # View all logs
make logs-api    # View API logs only
make restart     # Restart all services
make status      # Show service status
```

### Development
```bash
make dev         # Start with rebuild
make dev-rebuild # Force rebuild and start
make swagger     # Generate Swagger docs
make health      # Check service health
```

### Database Operations
```bash
make shell-db    # Open MongoDB shell
make db-backup   # Backup database
make db-restore BACKUP_PATH=./backup-folder  # Restore database
```

### Cleanup
```bash
make clean       # Remove all containers, networks, and images
make clean-volumes # Remove all volumes (⚠️ deletes data)
```

## 🌐 API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Authentication
- `POST /v1/users` - Create user account
- `POST /v1/users/login` - User login (returns JWT token)

### User Management (Protected)
- `GET /v1/users` - Get all users (with pagination)
- `GET /v1/users/{id}` - Get user by ID
- `PUT /v1/users/{id}` - Update user
- `DELETE /v1/users/{id}` - Delete user

### Documentation
- `GET /swagger/*` - Swagger UI (development only)

## 📚 API Documentation

When running in development mode, access the interactive Swagger documentation at:
```
http://localhost:5555/swagger/index.html
```

## 🔐 Authentication

The API uses JWT Bearer tokens for authentication. Include the token in the Authorization header:

```bash
Authorization: Bearer <your-jwt-token>
```

### Example Login Request
```bash
curl -X POST http://localhost:5555/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `API_PORT` | Server port | `5555` |
| `APP_ENV` | Environment (development/production) | `development` |
| `BASE_PATH` | API base path | `` |
| `MONGODB_HOST` | MongoDB host | `localhost` |
| `MONGODB_PORT` | MongoDB port | `27017` |
| `MONGODB_USERNAME` | MongoDB username | `root` |
| `MONGODB_PASSWORD` | MongoDB password | `password` |
| `MONGODB_DATABASE_NAME` | Database name | `backend_challenge` |
| `MONGODB_AUTH_SOURCE` | Authentication database | `admin` |
| `JWT_SECRET` | JWT signing secret | (required) |

### Environment Files

- `.env.docker` - Docker environment configuration
- `envs/.env.example` - Example environment file

## 🏗️ Project Structure

```
backend-golang/
├── cmd/api/                 # Application entry points
├── internal/
│   ├── api/controllers/     # HTTP handlers
│   ├── application/usecase/ # Business logic
│   ├── domain/              # Domain models and interfaces
│   └── infrastructure/      # External dependencies
├── pkg/utilities/           # Shared utilities
├── docs/                    # Swagger documentation
├── envs/                    # Environment files
├── docker-compose.yml       # Docker services definition
├── Dockerfile              # Container build instructions
└── Makefile                # Development commands
```

## 🧪 Testing

```bash
# Run tests locally
go test ./...

# Run tests in Docker
make test
```

## 📊 Monitoring & Health Checks

### Health Check Endpoint
```bash
curl http://localhost:5555/health
```

### Container Health Status
```bash
docker-compose ps
```

### Service Logs
```bash
# All services
make logs

# Specific service
make logs-api
make logs-db
```

## 🔧 Development

### Adding New Endpoints

1. Add Swagger annotations to your handler functions
2. Regenerate documentation: `make swagger`
3. Test the endpoint with the Swagger UI

### Database Migrations

Database indexes are automatically created on application startup. See `internal/infrastructure/database/models/` for index definitions.


## 🆘 Troubleshooting

### Common Issues

1. **MongoDB connection failed:**
   ```bash
   make logs-db  # Check MongoDB logs
   make shell-db # Test MongoDB connection
   ```

2. **API not responding:**
   ```bash
   make logs-api  # Check API logs
   make health    # Test health endpoints
   ```

3. **Docker issues:**
   ```bash
   make clean     # Clean up all containers
   make up        # Restart services
   ```

### Reset Everything
```bash
make clean-volumes  # ⚠️ This will delete all data
make up
```