# SRE AI Platform - 智能运维监控平台

基于 AI 的 SRE 智能监控平台，集成 Prometheus、Grafana、告警管理、AI 分析和配置中心功能。

## 功能特性

- **告警管理**: 告警接收、查看、确认、历史记录
- **告警规则**: 基于 PromQL 的告警规则配置
- **AI 分析**: 智能根因分析、趋势预测、异常检测、AI 对话助手
- **运维手册**: 结构化运维手册管理
- **用户管理**: 基于角色的用户管理 (admin/operator/viewer)，支持 JWT 认证和验证码
- **配置中心**: 分层配置管理（平台/AI/监控/集成），支持热更新
- **多租户**: 租户隔离与集群分配
- **Prometheus 集成**: 单集群配置管理、指标查询、连接健康检查
- **可视化仪表板**: 实时系统状态监控

## 技术栈

### 后端服务
| 服务 | 技术 | 端口 | 职责 |
|------|------|------|------|
| Gateway | Go + Gin | 8080 | API 网关、路由转发 |
| Core | Go + Gin | 8081 | 告警、规则、AI 分析、配置中心、用户认证、Prometheus 代理 |
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
│   │   ├── controller/   # alert, analysis, config, rule, cluster, auth, user, prometheus...
│   │   ├── service/      # 业务逻辑层
│   │   ├── repository/   # 数据访问层
│   │   ├── model/        # 数据模型
│   │   ├── ai/           # AI 服务集成 (OpenAI, Claude)
│   │   └── config/       # 配置管理
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
│   │   ├── pages/        # Dashboard, Alerts, Analysis, Rules, ConfigCenter, Login, UserManagement...
│   │   ├── services/     # API 服务
│   │   └── index.js
│   ├── package.json
│   └── Dockerfile
├── deploy/               # Docker Compose 部署配置
│   ├── docker-compose.yml
│   ├── prometheus-compose.yml
│   └── prometheus/       # Prometheus 配置
├── sql/
│   └── init.sql          # 数据库初始化脚本
├── build.sh              # 构建脚本
└── README.md
```

## 快速开始

### 环境要求
- Docker 20.10+
- Docker Compose 2.0+
- Go 1.21+ (本地开发)
- Node.js 18+ (本地开发)

### 一键部署

```bash
# 完整部署（构建、启动）
./build.sh rebuild
```

### 分步部署

```bash
# 1. 构建所有服务镜像
./build.sh build

# 2. 启动服务
./build.sh up
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
| MySQL | localhost:3306 | 数据库 (root/root123) |
| Redis | localhost:6379 | 缓存 |

### 默认登录账号

- **用户名**: `admin`
- **密码**: `sreAdmin550c`

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

# 重新构建并启动
./build.sh rebuild
```

## 系统架构详解

### 1. 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                          前端层 (Frontend)                        │
│                     React 18 + Ant Design 5                      │
│                         Port: 3000                               │
└───────────────────────────┬─────────────────────────────────────┘
                            │ HTTP
┌───────────────────────────▼─────────────────────────────────────┐
│                        网关层 (Gateway)                          │
│                      Go + Gin (Port: 8080)                       │
│              路由转发 / 负载均衡 / 统一入口                         │
└───────────────┬───────────────────────────────┬─────────────────┘
                │                               │
    ┌───────────▼──────────┐      ┌─────────────▼────────────┐
    │   Core Service       │      │    Runbook Service       │
    │   Port: 8081         │      │    Port: 8082            │
    │                      │      │                          │
    │  • 告警管理           │      │  • 运维手册CRUD           │
    │  • 告警规则           │      │  • 手册搜索               │
    │  • AI 分析            │      │                          │
    │  • 配置中心           │      └──────────────────────────┘
    │  • 用户认证           │
    │  • Prometheus 代理    │      ┌──────────────────────────┐
    │                      │      │    Tenant Service        │
    └───────────┬──────────┘      │    Port: 8083            │
                │                 │                          │
                │                 │  • 租户管理               │
    ┌───────────▼──────────┐      │  • 租户隔离               │
    │    数据层             │      │                          │
    │  ┌────────────────┐  │      └──────────────────────────┘
    │  │   MySQL 8.0    │  │
    │  │   Port: 3306   │  │
    │  └────────────────┘  │
    │  ┌────────────────┐  │
    │  │   Redis 7      │  │
    │  │   Port: 6379   │  │
    │  └────────────────┘  │
    └──────────────────────┘
```

### 2. 核心服务架构 (Core Service)

Core 服务是系统的核心业务服务，采用分层架构设计：

```
┌─────────────────────────────────────────────────────────────┐
│                     Controller 层                            │
│  处理 HTTP 请求，参数校验，调用 Service 层，返回响应            │
├─────────────────────────────────────────────────────────────┤
│  • AlertController      - 告警管理接口                        │
│  • RuleController       - 告警规则接口                        │
│  • ClusterController    - Prometheus 集群接口                 │
│  • AnalysisController   - AI 分析接口                         │
│  • ConfigController     - 配置中心接口                        │
│  • AuthController       - 认证接口 (登录/验证码)               │
│  • UserController       - 用户管理接口                        │
│  • AIModelController    - AI 模型配置接口                     │
│  • DashboardController  - 仪表板数据接口                      │
│  • PrometheusController - Prometheus 代理接口                 │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                      Service 层                              │
│              业务逻辑处理，协调 Repository 和外部服务          │
├─────────────────────────────────────────────────────────────┤
│  • AlertService         - 告警业务逻辑                        │
│  • RuleService          - 规则业务逻辑                        │
│  • ClusterService       - 集群管理逻辑                        │
│  • AnalysisService      - AI 分析业务逻辑                     │
│  • ConfigService        - 配置管理逻辑                        │
│  • AuthService          - 认证授权逻辑 (JWT)                  │
│  • UserService          - 用户管理逻辑                        │
│  • AIModelService       - AI 模型管理逻辑                     │
│  • DashboardService     - 仪表板数据聚合                      │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                    Repository 层                             │
│                  数据访问层，执行 SQL 操作                      │
├─────────────────────────────────────────────────────────────┤
│  • AlertRepository      - 告警数据访问                        │
│  • RuleRepository       - 规则数据访问                        │
│  • ClusterRepository    - 集群数据访问                        │
│  • AnalysisRepository   - 分析记录数据访问                     │
│  • ConfigRepository     - 配置数据访问                        │
│  • UserRepository       - 用户数据访问                        │
│  • AIModelRepository    - AI 模型数据访问                     │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                      Model 层                                │
│                  数据模型定义                                 │
├─────────────────────────────────────────────────────────────┤
│  • User, Tenant, SystemConfig                                │
│  • Alert, AlertRule, PrometheusCluster                       │
│  • AIAnalysis, AIModel                                       │
│  • Runbook                                                   │
└─────────────────────────────────────────────────────────────┘
```

### 3. AI 模块架构

```
┌─────────────────────────────────────────────────────────────┐
│                    AI 分析服务层                               │
│              AnalysisService - 协调 AI 分析流程               │
└───────────────────────────┬─────────────────────────────────┘
                            │
            ┌───────────────┼───────────────┐
            │               │               │
┌───────────▼────┐ ┌────────▼──────┐ ┌──────▼──────┐
│  AI Client     │ │  Prompt       │ │  Analysis   │
│  Interface     │ │  Templates    │ │  Service    │
├────────────────┤ ├───────────────┤ ├─────────────┤
│                │ │               │ │             │
│ • OpenAI Client│ │ • Root Cause  │ │ • Analysis  │
│ • Claude Client│ │ • Trend       │ │   Creation  │
│ • Azure Client │ │ • Anomaly     │ │ • Chat      │
│ • Custom Client│ │ • Capacity    │ │ • Health    │
│                │ │               │ │   Check     │
└────────────────┘ └───────────────┘ └─────────────┘
```

### 4. 认证授权流程

```
┌─────────┐     ┌──────────────┐     ┌─────────────┐     ┌──────────┐
│  Client │────▶│  /auth/login │────▶│  Verify     │────▶│ Generate │
│         │     │              │     │  Password   │     │ JWT      │
└─────────┘     └──────────────┘     └─────────────┘     └────┬─────┘
                                                               │
┌─────────┐     ┌──────────────┐     ┌─────────────┐          │
│  Client │◀────│  Return Token│◀────│  Store User │◀─────────┘
│         │     │              │     │  Info       │
└────┬────┘     └──────────────┘     └─────────────┘
     │
     │ Bearer Token
     ▼
┌──────────────┐     ┌─────────────┐     ┌─────────────┐
│  JWTAuth     │────▶│  Parse      │────▶│  Set User   │
│  Middleware  │     │  Token      │     │  Context    │
└──────────────┘     └─────────────┘     └─────────────┘
```

### 5. 配置中心架构

配置中心采用两级分层结构：

```
┌─────────────────────────────────────────────────────────────┐
│                    配置中心 (Config Center)                   │
├─────────────────────────────────────────────────────────────┤
│  一级分类 (Category)                                          │
│  ├── platform    - 平台基础配置                               │
│  ├── ai          - AI 相关配置                                │
│  ├── monitoring  - 监控配置                                   │
│  └── integration - 集成配置                                   │
├─────────────────────────────────────────────────────────────┤
│  二级分类 (SubCategory)                                       │
│  ├── platform/basic        - 基础设置（系统名称等）            │
│  ├── platform/system       - 系统设置                        │
│  ├── platform/notification - 通知设置                        │
│  ├── ai/models             - 模型配置                        │
│  ├── ai/strategy           - 策略配置                        │
│  └── monitoring/prometheus - Prometheus 配置                 │
└─────────────────────────────────────────────────────────────┘
```

### 6. 数据流图

#### 告警处理流程
```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ Prometheus  │────▶│ AlertManager│────▶│  Webhook    │────▶│   Core      │
│   Alert     │     │             │     │  /alerts/   │     │  Service    │
└─────────────┘     └─────────────┘     └─────────────┘     └──────┬──────┘
                                                                   │
                    ┌─────────────┐     ┌─────────────┐           │
                    │   Notify    │◀────│  Save to    │◀──────────┘
                    │   User      │     │  Database   │
                    └─────────────┘     └─────────────┘
```

#### AI 分析流程
```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   User      │────▶│  Create     │────▶│  Fetch      │────▶│   Call      │
│  Request    │     │  Analysis   │     │  Metrics    │     │   AI API    │
└─────────────┘     └─────────────┘     └─────────────┘     └──────┬──────┘
                                                                   │
┌─────────────┐     ┌─────────────┐     ┌─────────────┐           │
│   Display   │◀────│  Parse      │◀────│  AI         │◀──────────┘
│   Result    │     │  Response   │     │  Response   │
└─────────────┘     └─────────────┘     └─────────────┘
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

### 认证授权 API

| 方法 | 路径 | 描述 | 权限 |
|------|------|------|------|
| GET | `/api/v1/auth/captcha` | 获取验证码 | 公开 |
| POST | `/api/v1/auth/login` | 用户登录 | 公开 |
| POST | `/api/v1/auth/logout` | 用户登出 | 已登录 |
| GET | `/api/v1/auth/me` | 获取当前用户信息 | 已登录 |
| POST | `/api/v1/auth/change-password` | 修改密码 | 已登录 |
| POST | `/api/v1/auth/refresh` | 刷新 Token | 已登录 |

### 用户管理 API

| 方法 | 路径 | 描述 | 权限 |
|------|------|------|------|
| GET | `/api/v1/users` | 获取用户列表 | admin |
| GET | `/api/v1/users/:id` | 获取用户详情 | admin |
| POST | `/api/v1/users` | 创建用户 | admin |
| PUT | `/api/v1/users/:id` | 更新用户 | admin |
| DELETE | `/api/v1/users/:id` | 删除用户 | admin |
| POST | `/api/v1/users/:id/reset-password` | 重置密码 | admin |

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

#### Prometheus 集群（单集群配置）
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/clusters` | 获取集群列表 |
| GET | `/api/v1/clusters/:id` | 获取集群详情 |
| GET | `/api/v1/clusters/default` | 获取默认集群（带健康检查） |
| POST | `/api/v1/clusters` | 创建集群 |
| PUT | `/api/v1/clusters/:id` | 更新集群 |
| DELETE | `/api/v1/clusters/:id` | 删除集群 |
| POST | `/api/v1/clusters/:id/test` | 测试集群连接 |
| PUT | `/api/v1/clusters/:id/default` | 设为默认集群 |

#### AI 分析
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/analysis` | 获取分析列表 |
| POST | `/api/v1/analysis` | 创建分析 |
| GET | `/api/v1/analysis/:id` | 获取分析详情 |
| DELETE | `/api/v1/analysis/:id` | 删除分析 |
| PUT | `/api/v1/analysis/:id/archive` | 归档分析 |
| GET | `/api/v1/analysis/stats` | 获取统计信息 |
| POST | `/api/v1/ai/chat` | AI 对话 |
| GET | `/api/v1/ai/health` | AI 服务健康检查 |
| GET | `/api/v1/ai/model` | 获取 AI 模型信息 |

#### AI 模型配置
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/ai/configs` | 获取配置列表 |
| GET | `/api/v1/ai/configs/:id` | 获取配置详情 |
| POST | `/api/v1/ai/configs` | 创建配置 |
| PUT | `/api/v1/ai/configs/:id` | 更新配置 |
| DELETE | `/api/v1/ai/configs/:id` | 删除配置 |
| POST | `/api/v1/ai/configs/:id/test` | 测试配置 |
| PUT | `/api/v1/ai/configs/:id/default` | 设为默认 |
| GET | `/api/v1/ai/configs/active` | 获取激活的配置 |

#### 配置中心
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/v1/config/tree` | 获取配置树 |
| GET | `/api/v1/config/items` | 获取配置项列表 |
| GET | `/api/v1/config/items/:key` | 获取配置项 |
| PUT | `/api/v1/config/items/:key` | 更新配置值 |
| POST | `/api/v1/config/batch` | 批量更新 |
| POST | `/api/v1/config/items/:key/reset` | 重置为默认值 |
| GET | `/api/v1/config/app/name` | 获取系统名称（公开接口） |

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
| `users` | 系统用户表 |
| `system_configs` | 系统配置表（统一配置中心） |
| `ai_models` | AI 模型配置 |
| `prometheus_clusters` | Prometheus 集群配置 |
| `alert_rules` | 告警规则 |
| `alerts` | 告警记录 |
| `runbooks` | 运维手册 |
| `tenants` | 租户信息 |
| `ai_analysis` | AI 分析结果 |

### 表关系图

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│    tenants      │◀──────│     users       │       │  ai_models      │
├─────────────────┤       ├─────────────────┤       ├─────────────────┤
│ id (PK)         │       │ id (PK)         │       │ id (PK)         │
│ name            │       │ tenant_id (FK)  │       │ name            │
│ code            │       │ username        │       │ provider        │
│ status          │       │ password        │       │ model           │
└─────────────────┘       │ role            │       │ api_key         │
                          │ status          │       │ is_default      │
                          └─────────────────┘       └─────────────────┘

┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│prometheus_      │       │  alert_rules    │◀──────│     alerts      │
│clusters         │       ├─────────────────┤       ├─────────────────┤
├─────────────────┤       │ id (PK)         │       │ id (PK)         │
│ id (PK)         │       │ name            │       │ rule_id (FK)    │
│ name            │       │ expr            │       │ fingerprint     │
│ url             │       │ duration        │       │ status          │
│ is_default      │       │ severity        │       │ severity        │
└─────────────────┘       │ status          │       │ summary         │
                          └─────────────────┘       └─────────────────┘

┌─────────────────┐       ┌─────────────────┐
│  ai_analysis    │       │   runbooks      │
├─────────────────┤       ├─────────────────┤
│ id (PK)         │       │ id (PK)         │
│ alert_id (FK)   │       │ title           │
│ analysis_type   │       │ alert_name      │
│ input_data      │       │ severity        │
│ result          │       │ content         │
│ confidence      │       │ status          │
└─────────────────┘       └─────────────────┘

┌─────────────────┐
│ system_configs  │
├─────────────────┤
│ id (PK)         │
│ category        │
│ config_key      │
│ config_value    │
│ value_type      │
└─────────────────┘
```

## 配置中心

配置中心采用两级分层结构：

### 一级分类
- **platform**: 平台基础配置
- **ai**: AI 相关配置
- **monitoring**: 监控配置
- **integration**: 集成配置

### 二级分类示例
- `platform/basic`: 基础设置（系统名称等）
- `platform/system`: 系统设置
- `platform/notification`: 通知设置
- `ai/models`: 模型配置
- `ai/strategy`: 策略配置
- `monitoring/prometheus`: Prometheus 配置
- `monitoring/alert`: 告警配置

### 配置热更新

平台配置支持热更新：
1. 通过 API 或数据库修改配置
2. 配置将在 10 秒内自动生效
3. 无需重启服务

## AI 集成

支持多种 AI Provider：
- **OpenAI**: GPT-4, GPT-3.5-turbo
- **Claude**: Claude 3 Opus, Sonnet, Haiku
- **Azure OpenAI**
- **自定义**: 支持自定义 API 端点

配置方式：
1. 通过 Web UI 配置 AI 模型（配置中心 -> AI模型配置）
2. 或通过环境变量配置
3. 支持多模型配置和默认模型切换

### AI 功能
- **根因分析**: 分析告警根因
- **趋势分析**: 预测指标趋势
- **异常检测**: 识别异常模式
- **容量规划**: 资源使用预测
- **AI 对话**: 实时与 AI 助手对话

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
6. **安全**: 修改默认管理员密码，使用 HTTPS

## 注意事项

- **认证**: 已实现 JWT 认证 + 验证码，支持角色权限控制 (admin/operator/viewer)
- **单集群**: Prometheus 配置已简化为单集群模式，更易于管理
- **健康检查**: Prometheus 连接支持实时健康检查
- **开发阶段**: 适合演示和开发环境使用

## 许可证

MIT
