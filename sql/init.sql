-- SRE AI Platform - Optimized Database Schema
-- 精简设计原则: 删除冗余表,统一命名,减少JSON字段,完善审计

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- ============================================
-- 1. 核心模块 (租户、用户、配置)
-- ============================================

-- 租户表
CREATE TABLE IF NOT EXISTS tenants (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL COMMENT '租户名称',
    code VARCHAR(100) NOT NULL UNIQUE COMMENT '租户编码',
    description TEXT COMMENT '描述',
    status TINYINT DEFAULT 1 COMMENT '0:禁用 1:启用 2:暂停',
    settings JSON COMMENT '租户配置',
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_code (code),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户表';

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 1 COMMENT '租户ID',
    username VARCHAR(100) NOT NULL COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT 'bcrypt密码',
    email VARCHAR(255) COMMENT '邮箱',
    phone VARCHAR(50) COMMENT '手机号',
    role ENUM('admin', 'operator', 'viewer') DEFAULT 'viewer' COMMENT '角色',
    status TINYINT DEFAULT 1 COMMENT '0:禁用 1:启用',
    last_login_at TIMESTAMP NULL,
    last_login_ip VARCHAR(50),
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_tenant_user (tenant_id, username),
    INDEX idx_tenant (tenant_id),
    INDEX idx_role (role),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 系统配置表 (统一配置中心,替代config_items+platform_configs)
CREATE TABLE IF NOT EXISTS system_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    category VARCHAR(50) NOT NULL COMMENT '分类:platform/ai/monitoring',
    sub_category VARCHAR(50) COMMENT '子分类:basic/models/strategy',
    config_key VARCHAR(255) NOT NULL UNIQUE COMMENT '配置键',
    config_value TEXT COMMENT '配置值',
    value_type ENUM('string', 'number', 'boolean', 'json', 'password') DEFAULT 'string',
    is_sensitive BOOLEAN DEFAULT FALSE COMMENT '是否敏感',
    description VARCHAR(500),
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category),
    INDEX idx_sub_category (sub_category),
    INDEX idx_key (config_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- ============================================
-- 2. AI模块
-- ============================================

-- AI模型配置表
CREATE TABLE IF NOT EXISTS ai_models (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL COMMENT '配置名称',
    provider ENUM('openai', 'claude', 'azure', 'custom') NOT NULL,
    model VARCHAR(100) NOT NULL COMMENT '模型名',
    api_key VARCHAR(500) NOT NULL COMMENT 'API密钥',
    base_url VARCHAR(500) COMMENT '自定义地址',
    max_tokens INT DEFAULT 4000,
    temperature DECIMAL(3,2) DEFAULT 0.70,
    timeout INT DEFAULT 60,
    is_default BOOLEAN DEFAULT FALSE,
    is_enabled BOOLEAN DEFAULT TRUE,
    description TEXT,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_provider (provider),
    INDEX idx_default (is_default),
    INDEX idx_enabled (is_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI模型配置表';

-- AI分析记录表
CREATE TABLE IF NOT EXISTS ai_analysis (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    alert_id BIGINT COMMENT '关联告警ID',
    analysis_type ENUM('root_cause', 'trend', 'anomaly', 'capacity') DEFAULT 'root_cause',
    input_data JSON COMMENT '输入数据',
    result TEXT COMMENT '分析结果',
    confidence DECIMAL(5,2),
    status TINYINT DEFAULT 1 COMMENT '0:删除 1:正常 2:归档',
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_alert (alert_id),
    INDEX idx_type (analysis_type),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI分析记录表';

-- ============================================
-- 3. 监控告警模块 (简化设计)
-- ============================================

-- Prometheus集群表 (保留但简化,只支持单集群)
CREATE TABLE IF NOT EXISTS prometheus_clusters (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(500) NOT NULL,
    status TINYINT DEFAULT 1 COMMENT '0:禁用 1:启用',
    is_default BOOLEAN DEFAULT FALSE,
    config_json JSON COMMENT '配置JSON',
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_default (is_default)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Prometheus集群表';

-- 告警规则表
CREATE TABLE IF NOT EXISTS alert_rules (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL COMMENT '规则名称',
    description TEXT,
    expr VARCHAR(1000) NOT NULL COMMENT 'PromQL表达式',
    duration VARCHAR(50) DEFAULT '5m',
    severity ENUM('critical', 'warning', 'info') DEFAULT 'warning',
    labels JSON COMMENT '标签',
    annotations JSON COMMENT '注解',
    status TINYINT DEFAULT 1 COMMENT '0:禁用 1:启用',
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name),
    INDEX idx_severity (severity),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='告警规则表';

-- 告警记录表
CREATE TABLE IF NOT EXISTS alerts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    rule_id BIGINT COMMENT '关联规则ID',
    fingerprint VARCHAR(255) NOT NULL UNIQUE COMMENT '告警指纹',
    status ENUM('firing', 'resolved', 'acknowledged') DEFAULT 'firing',
    severity ENUM('critical', 'warning', 'info') DEFAULT 'warning',
    summary VARCHAR(500),
    description TEXT,
    labels JSON,
    starts_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ends_at TIMESTAMP NULL,
    acknowledged_by VARCHAR(100),
    acknowledged_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_rule (rule_id),
    INDEX idx_fingerprint (fingerprint),
    INDEX idx_status (status),
    INDEX idx_severity (severity),
    INDEX idx_starts (starts_at),
    FOREIGN KEY (rule_id) REFERENCES alert_rules(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='告警记录表';

-- ============================================
-- 4. 运维手册模块
-- ============================================

CREATE TABLE IF NOT EXISTS runbooks (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL COMMENT '标题',
    alert_name VARCHAR(255) COMMENT '关联告警名',
    severity ENUM('critical', 'warning', 'info') DEFAULT 'warning',
    content TEXT NOT NULL COMMENT '内容Markdown',
    steps JSON COMMENT '处理步骤',
    view_count INT DEFAULT 0,
    status TINYINT DEFAULT 1 COMMENT '0:禁用 1:启用',
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_alert_name (alert_name),
    INDEX idx_severity (severity),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='运维手册表';

-- ============================================
-- 5. 初始化数据
-- ============================================

-- 默认租户
INSERT INTO tenants (id, name, code, description, status, created_by) VALUES
(1, 'Default Tenant', 'default', 'Default system tenant', 1, 'system')
ON DUPLICATE KEY UPDATE id=id;

-- 管理员账号 (密码: sreAdmin550c)
INSERT INTO users (id, tenant_id, username, password, email, phone, role, status, created_by) VALUES
(1, 1, 'admin', '$2a$10$3EYxcEI2rQAE7oSWBwmVOej.K/VmxfZJpjlXrGWDt5.IyhJegQbPS', 'admin@sre.ai', '13800138000', 'admin', 1, 'system')
ON DUPLICATE KEY UPDATE id=id;

-- AI模型配置示例
INSERT INTO ai_models (name, provider, model, api_key, base_url, is_default, is_enabled, description, created_by) VALUES
('OpenAI GPT-4', 'openai', 'gpt-4', 'sk-your-key', 'https://api.openai.com/v1', TRUE, FALSE, 'OpenAI GPT-4', 'admin'),
('Claude 3 Opus', 'claude', 'claude-3-opus', 'sk-ant-key', 'https://api.anthropic.com/v1', FALSE, FALSE, 'Claude 3 Opus', 'admin'),
('minimax-m2.7', 'custom', 'minimax-m2.7', 'sk-cp-nmzeyymmpslsq7i2', 'https://cloud.infini-ai.com/maas/coding/v1', FALSE, TRUE, 'Minimax M2.7 模型', 'admin')
ON DUPLICATE KEY UPDATE name=name;

-- Prometheus集群
INSERT INTO prometheus_clusters (id, name, url, status, is_default, config_json, created_by) VALUES
(1, 'default', 'http://prometheus:9090', 1, TRUE, '{"scrape_interval":"15s"}', 'system')
ON DUPLICATE KEY UPDATE id=id;

-- 系统配置
INSERT INTO system_configs (category, config_key, config_value, value_type, description) VALUES
('platform', 'app.name', 'SRE AI Platform', 'string', '系统名称'),
('platform', 'app.logo', '', 'string', '系统Logo'),
('platform', 'app.version', '1.0.0', 'string', '系统版本'),
('ai', 'ai.default_model', '1', 'string', '默认AI模型ID'),
('ai', 'ai.max_tokens', '4000', 'number', '最大Token数'),
('monitoring', 'prometheus.timeout', '30', 'number', 'Prometheus查询超时'),
('monitoring', 'alert.default_severity', 'warning', 'string', '默认告警级别')
ON DUPLICATE KEY UPDATE config_key=config_key;

-- 示例告警规则
INSERT INTO alert_rules (name, description, expr, duration, severity, labels, annotations, status, created_by) VALUES
('HighCPUUsage', 'CPU > 80%', '100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80', '5m', 'warning', '{"team":"sre"}', '{"summary":"High CPU"}', 1, 'system'),
('HighMemoryUsage', 'Memory > 85%', '(node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100 > 85', '5m', 'critical', '{"team":"sre"}', '{"summary":"High Memory"}', 1, 'system'),
('DiskFull', 'Disk > 90%', '(node_filesystem_size_bytes - node_filesystem_free_bytes) / node_filesystem_size_bytes * 100 > 90', '10m', 'critical', '{"team":"sre"}', '{"summary":"Disk Full"}', 1, 'system')
ON DUPLICATE KEY UPDATE name=name;

-- 示例运维手册
INSERT INTO runbooks (title, alert_name, severity, content, steps, status, created_by) VALUES
('High CPU Usage处理', 'HighCPUUsage', 'warning', '# High CPU处理\n\n1. 使用top定位进程\n2. 检查应用日志\n3. 考虑扩容', '["定位进程","检查日志","扩容"]', 1, 'system'),
('High Memory处理', 'HighMemoryUsage', 'critical', '# High Memory处理\n\n1. 检查内存使用\n2. 分析堆转储\n3. 重启服务', '["检查内存","分析堆转储","重启"]', 1, 'system')
ON DUPLICATE KEY UPDATE title=title;
