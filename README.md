# SRE AI Platform - 智能监控平台

基于 AI 的 SRE 智能监控平台，集成 Prometheus、Grafana、告警管理、AI 分析和配置中心功能。

## 功能特性

- **告警管理**: 告警接收、查看、确认、历史记录
- **告警规则**: 基于 PromQL 的告警规则配置
- **AI 分析**: 智能根因分析、趋势预测、异常检测
- **运维手册**: 结构化运维手册管理
- **配置中心**: 分层配置管理（平台/AI/监控/集成）
- **多租户**: 租户隔离与集群分配
- **Prometheus 集成**: 多集群管理、指标查询
- **可视化仪表板**: 实时系统状态监控

## 技术栈

### 后端服务
| 服务 | 技术 | 端口 | 职责 |
|------|------|------|------|
| Gateway | Go + Gin | 8080 | API 网关、路由转发 |
| Core | Go + Gin | 8081 | 告警、规则、AI 分析、配置中心 |
| Runbook | Go + Gin | 8082 | 运维手册管理 |
| Tenant | Go + Gin | 8083 | 租户管理 |

### 基础设施
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **监控**: Prometheus + Grafana + Alertmanager
- **日志**: Zap
- **容器**: Docker + Docker Compose

### 前端
- React 18
- Ant Design 5
- Recharts (图表)
- Axios (HTTP 客户端)

## 项目结构

```
sre-ai/
├── backend/
│   ├── gateway/          # API 网关服务
│   │   ├── main.go
│   │   ├── config.yaml
│   │   ├── Dockerfile
│   │   ├── controller/   # HTTP 处理器
│   │   ├── service/      # 业务逻辑
│   │   ├── middleware/   # 中间件
│   │   └── config/       # 配置管理
│   ├── core/             # 核心服务
│   │   ├── main.go
│   │   ├── config.yaml
│   │   ├── Dockerfile
│   │   ├── controller/   # alert, analysis, config, rule, cluster...
│   │   ├── service/
│   │   ├── repository/   # 数据访问层
│   │   ├── model/        # 数据模型
│   │   ├── ai/           # AI 服务集成
│   │   └── config/
│   ├── runbook/          # 运维手册服务
│   │   ├── main.go
│   │   ├── config.yaml
│   │   └── ...
│   └── tenant/           # 租户服务
│       ├── main.go
│       ├── config.yaml
│       └── ...
├── frontend/             # React 前端应用
│   ├── src/
│   │   ├── pages/        # Dashboard, Alerts, Analysis, Rules, Config...
│   │   ├── components/   # 可复用组件
│   │   └── services/     # API 服务
│   ├── package.json
│   └── Dockerfile
├── deploy/               # Docker Compose 部署配置
│   ├── docker-compose.yml
│   ├── prometheus-compose.yml
│   └── prometheus/       # Prometheus 配置
├── sql/
│   └── init.sql          # 数据库初始化脚本
└── build.sh              # 构建脚本
```

## 快速开始

### 环境要求
- Docker 20.10+
- Docker Compose 2.0+
- Go 1.21+ (本地开发)
- Node.js 18+ (本地开发)

### 一键部署

```bash
# 完整部署（清理、构建、启动、初始化数据库）
./build.sh deploy-fresh
```

### 分步部署

```bash
# 1. 构建所有服务镜像
./build.sh build

# 2. 启动服务
./build.sh up

# 3. 初始化数据库（仅首次）
./build.sh init-db
```

### 访问服务

| 服务 | 地址 | 说明 |
|------|------|------|
| 前端 | http://localhost:3000 | React 应用 |
| Gateway API | http://localhost:8080 | API 网关 |
| Core API | http://localhost:8081 | 核心服务 |
| Runbook API | http://localhost:8082 | 运维手册服务 |
| Tenant API | http://localhost:8083 | 租户服务 |
| Prometheus | http://localhost:9090 | 指标收集 |
| Grafana | http://localhost:3001 | 可视化 (admin/admin) |
| Alertmanager | http://localhost:9093 | 告警管理 |
| MySQL | localhost:3306 | 数据库 |
| Redis | localhost:6379 | 缓存 |

### 常用命令

```bash
# 查看日志
./build.sh logs

# 重启服务
./build.sh restart

# 停止服务
./build.sh down

# 清理所有数据（包括镜像和卷）
./build.sh clean
```

## 配置说明

### 服务配置 (config.yaml)

每个后端服务都有独立的 `config.yaml`：

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

### 环境变量覆盖

| 变量 | 说明 |
|------|------|
| `PORT` | 服务端口 |
| `LOG_LEVEL` | 日志级别 |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | 数据库连接 |
| `REDIS_URL` | Redis 地址 |
| `CORE_SERVICE_URL`, `RUNBOOK_SERVICE_URL`, `TENANT_SERVICE_URL` | 后端服务地址 |

## API 文档

### Gateway 路由

| 路径 | 目标服务 |
|------|----------|
| `/api/core/*` | Core Service |
| `/api/runbook/*` | Runbook Service |
| `/api/tenant/*` | Tenant Service |
| `/health` | 健康检查 |
| `/metrics` | Prometheus 指标 |

### Core Service API

#### 告警管理
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/alerts` | 获取告警列表 |
| GET | `/api/v1/alerts/:id` | 获取告警详情 |
| POST | `/api/v1/alerts` | 创建告警 |
| POST | `/api/v1/alerts/webhook` | 接收告警 Webhook |
| PUT | `/api/v1/alerts/:id/ack` | 确认告警 |

#### 告警规则
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/rules` | 获取规则列表 |
| GET | `/api/v1/rules/:id` | 获取规则详情 |
| POST | `/api/v1/rules` | 创建规则 |
| PUT | `/api/v1/rules/:id` | 更新规则 |
| DELETE | `/api/v1/rules/:id` | 删除规则 |

#### Prometheus 集群
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/clusters` | 获取集群列表 |
| POST | `/api/v1/clusters` | 创建集群 |
| PUT | `/api/v1/clusters/:id` | 更新集群 |
| DELETE | `/api/v1/clusters/:id` | 删除集群 |

#### AI 分析
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/analysis` | 获取分析列表 |
| POST | `/api/v1/analysis` | 创建分析 |
| POST | `/api/v1/ai/chat` | AI 对话 |
| POST | `/api/v1/ai/chat/stream` | 流式 AI 对话 |
| GET | `/api/v1/ai/health` | AI 服务健康检查 |

#### AI 模型配置
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/ai/configs` | 获取配置列表 |
| POST | `/api/v1/ai/configs` | 创建配置 |
| PUT | `/api/v1/ai/configs/:id` | 更新配置 |
| POST | `/api/v1/ai/configs/:id/test` | 测试配置 |
| PUT | `/api/v1/ai/configs/:id/default` | 设为默认 |

#### 配置中心
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/config/tree` | 获取配置树 |
| GET | `/api/v1/config/items` | 获取配置项列表 |
| GET | `/api/v1/config/items/:key` | 获取配置项 |
| PUT | `/api/v1/config/items/:key` | 更新配置值 |
| POST | `/api/v1/config/batch` | 批量更新 |
| GET | `/api/v1/config/export` | 导出配置 |
| POST | `/api/v1/config/import` | 导入配置 |

#### Prometheus 代理
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/prometheus/query` | PromQL 查询 |
| GET | `/api/v1/prometheus/query_range` | 范围查询 |

#### 仪表板
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/dashboard` | 获取仪表板数据 |
| GET | `/api/v1/dashboard/metrics` | 获取指标数据 |

### Runbook Service API

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/runbooks` | 获取手册列表 |
| POST | `/api/v1/runbooks` | 创建手册 |
| GET | `/api/v1/runbooks/:id` | 获取手册详情 |
| PUT | `/api/v1/runbooks/:id` | 更新手册 |
| DELETE | `/api/v1/runbooks/:id` | 删除手册 |
| GET | `/api/v1/runbooks/search` | 搜索手册 |

### Tenant Service API

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/tenants` | 获取租户列表 |
| POST | `/api/v1/tenants` | 创建租户 |
| GET | `/api/v1/tenants/:id` | 获取租户详情 |
| PUT | `/api/v1/tenants/:id` | 更新租户 |
| DELETE | `/api/v1/tenants/:id` | 删除租户 |

## 数据库结构

### 核心表

| 表名 | 说明 |
|------|------|
| `config_items` | 配置中心配置项（新） |
| `ai_model_configs` | AI 模型配置 |
| `platform_configs` | 平台配置（兼容旧版） |
| `prometheus_clusters` | Prometheus 集群 |
| `alert_rules` | 告警规则 |
| `alerts` | 告警记录 |
| `runbooks` | 运维手册 |
| `tenants` | 租户信息 |
| `ai_analysis` | AI 分析结果 |
| `dashboard_configs` | 仪表板配置 |
| `metrics_cache` | 指标缓存 |

## 配置中心

配置中心采用两级分层结构：

### 一级分类
- **platform**: 平台基础配置
- **ai**: AI 相关配置
- **monitoring**: 监控配置
- **integration**: 集成配置

### 二级分类示例
- `platform/basic`: 基础设置
- `platform/system`: 系统设置
- `platform/notification`: 通知设置
- `ai/models`: 模型配置
- `ai/strategy`: 策略配置

### 配置热更新

平台配置支持热更新：
1. 通过 API 或数据库修改配置
2. 配置将在 10 秒内自动生效
3. 无需重启服务
4. 手动触发: `POST /api/v1/configs/reload`

## AI 集成

支持多种 AI Provider：
- **OpenAI**: GPT-4, GPT-3.5-turbo
- **Claude**: Claude 3 Opus, Sonnet, Haiku
- **Azure OpenAI**
- **自定义**: 支持自定义 API 端点

配置方式：
1. 通过 Web UI 配置 AI 模型
2. 或通过环境变量配置
3. 支持多模型配置和默认模型切换

## 开发指南

### 本地开发

```bash
# 启动基础设施
docker-compose -f deploy/docker-compose.yml up -d mysql redis

# 运行 Core 服务
cd backend/core
go run main.go

# 运行其他服务（新开终端）
cd backend/gateway && go run main.go
cd backend/runbook && go run main.go
cd backend/tenant && go run main.go
```

### 前端开发

```bash
cd frontend
npm install
npm start
```

前端将在 http://localhost:3000 启动，API 请求自动代理到 Gateway。

### 添加新 API 端点

1. **Core Service 示例**:
```go
// controller/my_controller.go
func (c *MyController) MyHandler(ctx *gin.Context) {
    // 处理逻辑
    ctx.JSON(200, gin.H{"data": result})
}

// main.go
v1.GET("/my-endpoint", myController.MyHandler)
```

2. **Gateway 路由**（如需新服务）:
```go
// 在 gateway/main.go 中添加
api.Any("/myservice/*path", gatewayController.ProxyToMyService)
```

## 生产部署建议

1. **数据库**: 使用外部托管 MySQL（如 RDS）
2. **Redis**: 使用 Redis Cluster 或托管服务
3. **配置**: 使用环境变量覆盖敏感信息
4. **日志**: 配置日志收集（如 ELK、Fluentd）
5. **监控**: 配置 Prometheus 远程存储
6. **安全**: 添加认证层（项目当前无认证）

## 注意事项

- **无认证**: 当前版本未实现认证层
- **无测试**: 项目暂无单元/集成测试
- **开发阶段**: 适合演示和开发环境使用

## 许可证

MIT
