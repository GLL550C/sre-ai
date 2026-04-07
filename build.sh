#!/bin/bash

# SRE AI Platform 构建脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 目录配置
DEPLOY_DIR="deploy"
COMPOSE_FILE="$DEPLOY_DIR/docker-compose.yml"
PROMETHEUS_COMPOSE="$DEPLOY_DIR/prometheus-compose.yml"

# 日志函数
info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 构建所有镜像
build() {
    info "构建后端服务镜像..."
    for service in gateway core runbook tenant; do
        docker build -t sre-$service:latest -f "backend/$service/Dockerfile" "backend/$service"
    done

    info "构建前端镜像..."
    docker build --no-cache -t sre-frontend:latest -f "frontend/Dockerfile" "frontend"

    info "构建完成"
}

# 启动服务
up() {
    info "启动服务..."
    docker-compose -f "$COMPOSE_FILE" up -d
    docker-compose -f "$PROMETHEUS_COMPOSE" up -d
    info "服务已启动"
    info "前端: http://localhost:3000"
    info "网关: http://localhost:8080"
}

# 停止服务
down() {
    info "停止服务..."
    docker-compose -f "$PROMETHEUS_COMPOSE" down 2>/dev/null || true
    docker-compose -f "$COMPOSE_FILE" down 2>/dev/null || true
    info "服务已停止"
}

# 停止服务并删除卷
down_volumes() {
    info "停止服务并删除数据卷..."
    docker-compose -f "$PROMETHEUS_COMPOSE" down -v 2>/dev/null || true
    docker-compose -f "$COMPOSE_FILE" down -v 2>/dev/null || true
    info "服务已停止，数据卷已删除"
}

# 重启服务
restart() {
    info "重启服务..."
    down
    sleep 2
    up
}

# 重新构建并启动（删除卷，全新部署）
rebuild() {
    info "重新构建镜像并启动服务..."
    down_volumes
    build
    up
}

# 主命令
command="${1:-}"

case "$command" in
    build)
        build
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
    *)
        echo "SRE AI Platform 构建脚本"
        echo ""
        echo "用法: $0 {build|up|down|restart|rebuild}"
        echo ""
        echo "命令:"
        echo "  build    - 重新构建所有镜像"
        echo "  up       - 启动所有服务"
        echo "  down     - 停止所有服务"
        echo "  restart  - 重启所有服务"
        echo "  rebuild  - 删除卷、重新构建镜像、启动服务"
        echo ""
        echo "示例:"
        echo "  $0 rebuild   # 代码更新后重新部署"
        echo "  $0 restart   # 快速重启服务"
        exit 1
        ;;
esac
