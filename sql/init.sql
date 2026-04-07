-- SRE AI Platform Database Schema

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '加密密码',
    email VARCHAR(255) COMMENT '邮箱',
    phone VARCHAR(50) COMMENT '手机号',
    role ENUM('admin', 'operator', 'viewer') DEFAULT 'viewer' COMMENT '角色: admin管理员, operator运维, viewer访客',
    status INT DEFAULT 1 COMMENT '1:active, 0:inactive',
    last_login_at TIMESTAMP NULL COMMENT '最后登录时间',
    last_login_ip VARCHAR(50) COMMENT '最后登录IP',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_role (role),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 配置项定义表(新配置中心)
CREATE TABLE IF NOT EXISTS config_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    category VARCHAR(50) NOT NULL COMMENT '一级分类: platform/ai/monitoring/integration',
    sub_category VARCHAR(50) NOT NULL COMMENT '二级分类',
    key_name VARCHAR(255) NOT NULL UNIQUE COMMENT '配置键',
    value TEXT COMMENT '配置值',
    type VARCHAR(20) DEFAULT 'string' COMMENT '类型: string/number/boolean/json/password',
    options TEXT COMMENT '可选项(JSON数组)',
    default_val TEXT COMMENT '默认值',
    required BOOLEAN DEFAULT FALSE COMMENT '是否必填',
    `sensitive` BOOLEAN DEFAULT FALSE COMMENT '是否敏感',
    description VARCHAR(500) COMMENT '描述',
    sort_order INT DEFAULT 0 COMMENT '排序',
    icon VARCHAR(100) COMMENT '图标',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category),
    INDEX idx_sub_category (sub_category),
    INDEX idx_key (key_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- AI 模型配置表
CREATE TABLE IF NOT EXISTS ai_model_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL COMMENT '配置名称',
    provider VARCHAR(50) NOT NULL COMMENT 'AI provider: openai, claude, azure, custom',
    model VARCHAR(100) NOT NULL COMMENT 'Model name: gpt-4, claude-3-opus, etc.',
    api_key VARCHAR(500) NOT NULL COMMENT 'API Key',
    base_url VARCHAR(500) COMMENT 'Custom API base URL',
    max_tokens INT DEFAULT 4000 COMMENT 'Max tokens per request',
    temperature DECIMAL(3,2) DEFAULT 0.70 COMMENT 'Temperature 0-2',
    timeout INT DEFAULT 60 COMMENT 'Request timeout in seconds',
    is_default BOOLEAN DEFAULT FALSE COMMENT 'Is default configuration',
    is_enabled BOOLEAN DEFAULT TRUE COMMENT 'Is enabled',
    description TEXT COMMENT 'Description',
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_provider (provider),
    INDEX idx_is_default (is_default),
    INDEX idx_is_enabled (is_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 平台配置表(保留用于兼容性,逐步迁移到config_items)
CREATE TABLE IF NOT EXISTS platform_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    config_key VARCHAR(255) NOT NULL UNIQUE,
    config_value TEXT,
    description VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_config_key (config_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Prometheus 集群配置表
CREATE TABLE IF NOT EXISTS prometheus_clusters (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(500) NOT NULL,
    status INT DEFAULT 1 COMMENT '1:active, 0:inactive',
    is_default BOOLEAN DEFAULT FALSE COMMENT '是否为默认集群',
    config_json TEXT COMMENT 'Prometheus configuration in JSON',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name),
    INDEX idx_status (status),
    INDEX idx_is_default (is_default)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 告警规则表
CREATE TABLE IF NOT EXISTS alert_rules (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    expr VARCHAR(1000) NOT NULL COMMENT 'PromQL expression',
    duration VARCHAR(50) DEFAULT '5m' COMMENT 'Duration like 5m, 1h',
    severity ENUM('critical', 'warning', 'info') DEFAULT 'warning',
    labels JSON COMMENT 'Additional labels',
    annotations JSON COMMENT 'Alert annotations',
    status INT DEFAULT 1 COMMENT '1:active, 0:inactive',
    cluster_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (cluster_id) REFERENCES prometheus_clusters(id) ON DELETE SET NULL,
    INDEX idx_name (name),
    INDEX idx_severity (severity),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 告警记录表
CREATE TABLE IF NOT EXISTS alerts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    rule_id BIGINT,
    fingerprint VARCHAR(255) NOT NULL UNIQUE COMMENT 'Alert fingerprint',
    status ENUM('firing', 'resolved', 'acknowledged') DEFAULT 'firing',
    severity ENUM('critical', 'warning', 'info') DEFAULT 'warning',
    summary VARCHAR(500),
    description TEXT,
    labels JSON,
    starts_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ends_at TIMESTAMP NULL,
    acknowledged_by VARCHAR(255),
    acknowledged_at TIMESTAMP NULL,
    cluster_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (rule_id) REFERENCES alert_rules(id) ON DELETE SET NULL,
    FOREIGN KEY (cluster_id) REFERENCES prometheus_clusters(id) ON DELETE SET NULL,
    INDEX idx_status (status),
    INDEX idx_severity (severity),
    INDEX idx_starts_at (starts_at),
    INDEX idx_fingerprint (fingerprint)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 运维手册表
CREATE TABLE IF NOT EXISTS runbooks (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    alert_name VARCHAR(255) COMMENT 'Associated alert name pattern',
    severity ENUM('critical', 'warning', 'info') DEFAULT 'warning',
    content TEXT NOT NULL COMMENT 'Runbook content in Markdown',
    steps JSON COMMENT 'Step-by-step resolution steps',
    related_alerts JSON COMMENT 'Related alert patterns',
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    view_count INT DEFAULT 0,
    status INT DEFAULT 1 COMMENT '1:active, 0:inactive',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_title (title),
    INDEX idx_alert_name (alert_name),
    INDEX idx_severity (severity),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 租户表
CREATE TABLE IF NOT EXISTS tenants (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE COMMENT 'Tenant unique code',
    description TEXT,
    config JSON COMMENT 'Tenant specific configuration',
    clusters JSON COMMENT 'Assigned cluster IDs',
    status INT DEFAULT 1 COMMENT '1:active, 0:inactive, 2:suspended',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name),
    INDEX idx_code (code),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- AI 分析结果表
CREATE TABLE IF NOT EXISTS ai_analysis (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    alert_id BIGINT,
    cluster_id BIGINT,
    alert_fingerprint VARCHAR(255),
    analysis_type ENUM('root_cause', 'trend', 'anomaly', 'capacity', 'correlation', 'recommendation') DEFAULT 'root_cause',
    analysis_mode VARCHAR(50) DEFAULT 'manual' COMMENT 'realtime, historical, predictive, manual',
    input_data JSON COMMENT 'Input data for analysis',
    result TEXT COMMENT 'AI analysis result',
    root_cause TEXT COMMENT 'Identified root cause',
    suggestions JSON COMMENT 'List of suggestions',
    related_alerts JSON COMMENT 'Related alert IDs',
    confidence DECIMAL(5,2) COMMENT 'Confidence score 0-100',
    model_version VARCHAR(100),
    status INT DEFAULT 1 COMMENT '0:deleted, 1:active, 2:archived',
    created_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (alert_id) REFERENCES alerts(id) ON DELETE SET NULL,
    FOREIGN KEY (cluster_id) REFERENCES prometheus_clusters(id) ON DELETE SET NULL,
    INDEX idx_alert_fingerprint (alert_fingerprint),
    INDEX idx_cluster_id (cluster_id),
    INDEX idx_analysis_type (analysis_type),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 仪表板配置表
CREATE TABLE IF NOT EXISTS dashboard_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    tenant_id BIGINT,
    layout JSON COMMENT 'Dashboard layout configuration',
    widgets JSON COMMENT 'Widget configurations',
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    INDEX idx_name (name),
    INDEX idx_tenant (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 指标缓存表
CREATE TABLE IF NOT EXISTS metrics_cache (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    metric_key VARCHAR(500) NOT NULL,
    metric_value TEXT,
    cluster_id BIGINT,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_metric_key (metric_key),
    INDEX idx_expires_at (expires_at),
    INDEX idx_cluster (cluster_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认AI配置示例
INSERT INTO ai_model_configs (name, provider, model, api_key, base_url, max_tokens, temperature, timeout, is_default, is_enabled, description, created_by) VALUES
('OpenAI GPT-4', 'openai', 'gpt-4', 'sk-your-api-key-here', 'https://api.openai.com/v1', 4000, 0.7, 60, TRUE, FALSE, 'OpenAI GPT-4 model for analysis', 'admin'),
('Claude 3 Opus', 'claude', 'claude-3-opus-20240229', 'sk-ant-api-your-key-here', 'https://api.anthropic.com/v1', 4000, 0.7, 60, FALSE, FALSE, 'Anthropic Claude 3 Opus model', 'admin');

-- 插入默认 Prometheus 集群
INSERT INTO prometheus_clusters (name, url, status, is_default, config_json) VALUES
('default', 'http://prometheus:9090', 1, TRUE, '{"scrape_interval": "15s", "evaluation_interval": "15s"}');

-- 插入示例告警规则
INSERT INTO alert_rules (name, description, expr, duration, severity, labels, annotations, status) VALUES
('HighCPUUsage', 'CPU usage is above 80%', '100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80', '5m', 'warning', '{"team": "sre"}', '{"summary": "High CPU usage detected", "description": "CPU usage is above 80% on {{ $labels.instance }}"}', 1),
('HighMemoryUsage', 'Memory usage is above 85%', '(node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100 > 85', '5m', 'critical', '{"team": "sre"}', '{"summary": "High memory usage detected", "description": "Memory usage is above 85% on {{ $labels.instance }}"}', 1),
('DiskFull', 'Disk usage is above 90%', '(node_filesystem_size_bytes - node_filesystem_free_bytes) / node_filesystem_size_bytes * 100 > 90', '10m', 'critical', '{"team": "sre"}', '{"summary": "Disk is almost full", "description": "Disk usage is above 90% on {{ $labels.instance }}"}', 1);

-- 插入示例运维手册
INSERT INTO runbooks (title, alert_name, severity, content, steps, related_alerts, created_by, status) VALUES
('High CPU Usage Resolution', 'HighCPUUsage', 'warning', '# High CPU Usage Resolution\n\n## Overview\nThis runbook helps resolve high CPU usage alerts.\n\n## Common Causes\n1. Application memory leak\n2. Increased traffic\n3. Background jobs\n4. Resource contention', '["1. Identify the process consuming high CPU using top/htop", "2. Check application logs for errors", "3. Review recent deployments", "4. Consider scaling if traffic increased", "5. Restart service if necessary"]', '["HighCPUUsage", "ProcessHighCPU"]', 'admin', 1),
('High Memory Usage Resolution', 'HighMemoryUsage', 'critical', '# High Memory Usage Resolution\n\n## Overview\nThis runbook helps resolve high memory usage alerts.\n\n## Common Causes\n1. Memory leak in application\n2. Large data processing\n3. Cache not expiring\n4. Connection pool exhaustion', '["1. Check memory usage by process", "2. Analyze heap dumps if available", "3. Review recent code changes", "4. Check for memory leaks", "5. Restart service if OOM risk"]', '["HighMemoryUsage", "OOMKilled"]', 'admin', 1),
('Disk Full Resolution', 'DiskFull', 'critical', '# Disk Full Resolution\n\n## Overview\nThis runbook helps resolve disk full alerts.\n\n## Common Causes\n1. Log files growing too large\n2. Temporary files not cleaned\n3. Database growing\n4. Large files uploaded', '["1. Check disk usage: df -h", "2. Find large files: du -sh /*", "3. Clean old logs if safe", "4. Check for core dumps", "5. Expand disk if necessary"]', '["DiskFull", "DiskHighUsage"]', 'admin', 1);

-- 插入示例租户
INSERT INTO tenants (name, code, description, config, clusters, status) VALUES
('Default Tenant', 'default', 'Default system tenant', '{"timezone": "UTC", "language": "en"}', '[1]', 1),
('Production', 'production', 'Production environment', '{"timezone": "UTC", "language": "en"}', '[1]', 1);

-- 插入默认管理员账号 (密码: sreAdmin550c, bcrypt加密)
INSERT INTO users (username, password, email, phone, role, status, created_at, updated_at) VALUES
('admin', '$2a$10$3EYxcEI2rQAE7oSWBwmVOej.K/VmxfZJpjlXrGWDt5.IyhJegQbPS', 'admin@sre.ai', '13800138000', 'admin', 1, NOW(), NOW());
