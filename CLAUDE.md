# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SRE AI Platform - An intelligent monitoring platform based on AI, integrating Prometheus, Grafana, alert management, AI analysis, and configuration center features.

## Architecture

### Backend Services (Go + Gin)

The backend follows a microservices architecture with 4 services:

1. **gateway** (port 8080): API Gateway that routes requests to backend services
   - Routes: `/api/core/*` → Core, `/api/runbook/*` → Runbook, `/api/tenant/*` → Tenant
   - Health check: `/health`, Metrics: `/metrics`

2. **core** (port 8081): Core service for alerts, rules, clusters, dashboard, AI analysis, user auth
   - Key features: Alert management, Prometheus cluster management, Alert rules, AI analysis, Dashboard data, JWT authentication
   - Config hot-reload via `ConfigManager` that polls DB every 10 seconds

3. **runbook** (port 8082): Runbook service for operational manuals
   - CRUD operations for runbooks associated with alert patterns

4. **tenant** (port 8083): Tenant service for multi-tenancy management

Each service has identical structure: `config/`, `controller/`, `service/`, `repository/`, `model/`, `middleware/`

### Frontend (React 18 + Ant Design 5)

- Port 3000 (proxied to gateway at 8080 in dev)
- Pages: Dashboard, Alerts, Analysis, Rules, Config, Login, UserManagement
- Uses Recharts for visualization, Axios for HTTP

### Infrastructure

- MySQL 8.0 (port 3306): Main database (root/root123)
- Redis (port 6379): Cache and session storage
- Prometheus (port 9090): Metrics collection
- Grafana (port 3001): Visualization (admin/admin)

## Common Commands

### Build & Deploy

```bash
# Build all services (Docker images)
./build.sh build

# Start all services
./build.sh up

# Stop all services
./build.sh down

# Restart services
./build.sh restart

# View logs
./build.sh logs

# Clean everything (containers, volumes, images)
./build.sh clean

# Fresh deployment (clean + build + up)
./build.sh rebuild
```

### Database Initialization

```bash
# Initialize database (after MySQL is healthy)
docker exec -i sre-mysql mysql -uroot -proot123 sre_platform < sql/init.sql
```

### Local Development

```bash
# Start infrastructure
docker-compose -f deploy/docker-compose.yml up -d mysql redis

# Run Core service locally (requires local MySQL/Redis)
cd backend/core && go run main.go

# Run frontend locally
cd frontend && npm install && npm start
```

## Service Configuration

Each backend service has a `config.yaml`:

```yaml
server:
  port: "8081"
  read_timeout: 30
  write_timeout: 30

log:
  level: "info"      # debug, info, warn, error
  format: "json"     # json, console
  output_path: "stdout"

database:
  host: "localhost"
  port: "3306"
  user: "root"
  password: "root123"
  name: "sre_platform"
  max_open_conns: 25
  max_idle_conns: 10

redis:
  addr: "localhost:6379"
  password: ""
  db: 0
```

Environment variables override config values: `PORT`, `LOG_LEVEL`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `REDIS_URL`

## Key Code Patterns

### Backend Service Initialization (core/main.go pattern)

1. Load config: `config.LoadConfig("config.yaml")`
2. Init logger: `config.InitLogger(cfg.Log)`
3. Init DB: `sql.Open("mysql", dsn)` with connection pool settings
4. Init Redis: `redis.NewClient()`
5. Init ConfigManager for hot-reload: `config.NewConfigManager(db, redisClient, logger)`
6. Init repositories → services → controllers (dependency injection)
7. Setup Gin router with middleware: `Logger`, `Recovery`, `CORS`, `JWTAuth`
8. Register routes and start server

### API Response Pattern

Controllers use `c.JSON(200, gin.H{"data": result})` for success and `c.JSON(500, gin.H{"error": err.Error()})` for errors.

### Hierarchical Config System

The config center uses a two-level hierarchy:
- **Level 1 Categories**: platform, ai, monitoring, integration
- **Level 2 Subcategories**: basic, system, notification, models, strategy, etc.

Key files:
- Model: `backend/core/model/config_item.go`
- Repository: `backend/core/repository/config_repository.go`
- Service: `backend/core/service/config_service.go`
- Controller: `backend/core/controller/config_controller.go`
- Frontend: `frontend/src/pages/ConfigCenter.js`

## Authentication

The system implements JWT authentication with role-based access control:
- Roles: admin, operator, viewer
- Default login: admin/sreAdmin550c
- Public endpoints: `/api/v1/auth/*`, `/health`, `/metrics`

## Important Notes

- **No automated tests**: The project currently has no unit or integration tests
- **AI Integration**: Supports OpenAI GPT-4 and Claude 3 via configurable providers
- **Multi-tenancy**: Tenant isolation with cluster assignment is supported

## Service URLs

| Service | URL |
|---------|-----|
| Frontend | http://localhost:3000 |
| Gateway API | http://localhost:8080 |
| Core API | http://localhost:8081 |
| Runbook API | http://localhost:8082 |
| Tenant API | http://localhost:8083 |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3001 (admin/admin) |
