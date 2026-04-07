#!/bin/bash

# SRE AI Platform Build Script
# Optimized by Staff+ SRE Architect

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Directories
BACKEND_DIR="backend"
FRONTEND_DIR="frontend"
DEPLOY_DIR="deploy"
SQL_DIR="sql"

# Compose files
COMPOSE_FILE="$DEPLOY_DIR/docker-compose.yml"
PROMETHEUS_COMPOSE="$DEPLOY_DIR/prometheus-compose.yml"

# Helper Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# Check if docker and docker-compose are available
check_prerequisites() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed or not in PATH"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        log_error "docker-compose is not installed or not in PATH"
        exit 1
    fi
}

# Build backend services
build_backend() {
    log_info "Building backend services with Docker..."

    for service in gateway core runbook tenant; do
        log_info "Building $service service image..."
        docker build -t sre-$service:latest -f "$BACKEND_DIR/$service/Dockerfile" "$BACKEND_DIR/$service"
    done

    log_info "Backend build completed!"
}

# Build frontend
build_frontend() {
    log_info "Building frontend with Docker (no cache)..."
    docker build --no-cache -t sre-frontend:latest -f "$FRONTEND_DIR/Dockerfile" "$FRONTEND_DIR"
    log_info "Frontend build completed!"
}

# Full build
build() {
    check_prerequisites
    log_info "Starting full build..."
    build_backend
    build_frontend
    log_info "Build completed successfully!"
}

# Initialize database
init_db() {
    check_prerequisites
    log_info "Initializing database..."

    # Check if MySQL container is running
    if ! docker ps | grep -q sre-mysql; then
        log_warn "MySQL container is not running. Starting MySQL and Redis..."
        docker-compose -f "$COMPOSE_FILE" up -d mysql redis

        # Wait for MySQL to be healthy
        log_info "Waiting for MySQL to be healthy..."
        retries=0
        max_retries=30
        until docker exec sre-mysql mysqladmin ping -h localhost -proot123 --silent 2>/dev/null || [ $retries -eq $max_retries ]; do
            sleep 2
            retries=$((retries + 1))
            log_debug "Waiting for MySQL... ($retries/$max_retries)"
        done

        if [ $retries -eq $max_retries ]; then
            log_error "MySQL failed to become healthy within timeout"
            exit 1
        fi
    fi

    # Check if database is already initialized (has tables)
    table_count=$(docker exec sre-mysql mysql -uroot -proot123 -N -s -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='sre_platform'" 2>/dev/null || echo "0")

    if [ "$table_count" -gt "0" ]; then
        log_warn "Database already has $table_count tables. Skipping initialization."
        log_warn "Use './build.sh clean' first if you want to reset the database."
        return 0
    fi

    # Execute init SQL
    log_info "Executing init.sql..."
    if docker exec -i sre-mysql mysql -uroot -proot123 sre_platform < "$SQL_DIR/init.sql" 2>/dev/null; then
        log_info "Database initialized successfully!"
    else
        log_error "Failed to initialize database"
        exit 1
    fi
}

# Start all services
up() {
    check_prerequisites
    log_info "Creating network if not exists..."
    docker network create sre-network 2>/dev/null || true

    log_info "Starting core services..."
    docker-compose -f "$COMPOSE_FILE" up -d

    log_info "Starting monitoring services..."
    docker-compose -f "$PROMETHEUS_COMPOSE" up -d

    log_info "All services started!"
    log_info "Frontend: http://localhost:3000"
    log_info "Gateway: http://localhost:8080"
    log_info "Prometheus: http://localhost:9090"
    log_info "Grafana: http://localhost:3001"
}

# Stop all services
down() {
    check_prerequisites
    log_info "Stopping all services..."

    # Stop prometheus services first
    if [ -f "$PROMETHEUS_COMPOSE" ]; then
        docker-compose -f "$PROMETHEUS_COMPOSE" down 2>/dev/null || true
    fi

    # Stop main services
    if [ -f "$COMPOSE_FILE" ]; then
        docker-compose -f "$COMPOSE_FILE" down 2>/dev/null || true
    fi

    log_info "All services stopped!"
}

# Stop services and remove volumes (keeps images)
down_volumes() {
    check_prerequisites
    log_info "Stopping all services and removing volumes..."

    if [ -f "$PROMETHEUS_COMPOSE" ]; then
        docker-compose -f "$PROMETHEUS_COMPOSE" down -v 2>/dev/null || true
    fi

    if [ -f "$COMPOSE_FILE" ]; then
        docker-compose -f "$COMPOSE_FILE" down -v 2>/dev/null || true
    fi

    log_info "Services stopped and volumes removed!"
}

# Restart services (keep images and volumes)
restart() {
    log_info "Restarting services..."
    down
    sleep 2
    up
}

# Rebuild images and restart
rebuild() {
    log_info "Rebuilding images and restarting services..."
    down
    build
    up
}

# Clean everything (containers, volumes, images)
clean() {
    check_prerequisites
    log_info "Cleaning all containers, volumes and images..."

    # Stop all services and remove volumes properly
    down_volumes

    # Remove custom images
    log_info "Removing SRE images..."
    images=$(docker images -q sre-frontend sre-gateway sre-core sre-runbook sre-tenant 2>/dev/null || true)
    if [ -n "$images" ]; then
        docker rmi $images 2>/dev/null || true
    fi

    # Clean up dangling images
    log_info "Cleaning up dangling images..."
    docker image prune -f 2>/dev/null || true

    log_info "Clean completed!"
}

# Full fresh deployment
deploy_fresh() {
    log_warn "Starting fresh deployment (will DELETE all existing data)..."
    clean
    log_info "Building all services..."
    build
    log_info "Starting services with fresh database..."
    up
    log_info "Waiting for services to be ready..."
    sleep 15
    log_info "Initializing database..."
    init_db
    log_info "Fresh deployment completed!"
}

# View logs
logs() {
    docker-compose -f "$COMPOSE_FILE" logs -f
}

# Health check
health() {
    log_info "Checking service health..."

    services="sre-mysql sre-redis sre-core sre-gateway sre-frontend sre-runbook sre-tenant"
    all_healthy=true

    for service in $services; do
        if docker ps | grep -q "$service"; then
            status=$(docker inspect --format='{{.State.Status}}' "$service" 2>/dev/null || echo "unknown")
            health=$(docker inspect --format='{{.State.Health.Status}}' "$service" 2>/dev/null || echo "N/A")
            if [ "$health" != "N/A" ]; then
                log_info "$service: $status (health: $health)"
            else
                log_info "$service: $status"
            fi
        else
            log_error "$service: NOT RUNNING"
            all_healthy=false
        fi
    done

    if [ "$all_healthy" = true ]; then
        log_info "All services are running!"
    else
        log_error "Some services are not running!"
        exit 1
    fi
}

# Status check
status() {
    echo "========================================"
    echo "SRE AI Platform - Service Status"
    echo "========================================"
    docker-compose -f "$COMPOSE_FILE" ps 2>/dev/null || echo "No services running"
    echo ""
    echo "Images:"
    docker images | grep sre- || echo "No SRE images found"
}

# Main
case "${1:-}" in
    build)
        build
        ;;
    init-db)
        init_db
        ;;
    up)
        up
        ;;
    down)
        down
        ;;
    restart)
        restart
        ;;
    rebuild)
        rebuild
        ;;
    clean)
        clean
        ;;
    deploy-fresh)
        deploy_fresh
        ;;
    logs)
        logs
        ;;
    health)
        health
        ;;
    status)
        status
        ;;
    *)
        echo "SRE AI Platform Build Script"
        echo ""
        echo "Usage: $0 {build|init-db|up|down|restart|rebuild|clean|deploy-fresh|logs|health|status}"
        echo ""
        echo "Commands:"
        echo "  build         - Build all Docker images"
        echo "  init-db       - Initialize database (safe to run multiple times)"
        echo "  up            - Start all services"
        echo "  down          - Stop all services (keep volumes)"
        echo "  restart       - Restart all services (keep images and volumes)"
        echo "  rebuild       - Rebuild images and restart services"
        echo "  clean         - Remove containers, volumes, images (DELETES ALL DATA)"
        echo "  deploy-fresh  - Full clean deployment with fresh database"
        echo "  logs          - View service logs"
        echo "  health        - Check service health status"
        echo "  status        - Show container and image status"
        echo ""
        echo "Examples:"
        echo "  $0 deploy-fresh    # First time setup"
        echo "  $0 rebuild         # After code changes"
        echo "  $0 restart         # Quick restart"
        exit 1
        ;;
esac
