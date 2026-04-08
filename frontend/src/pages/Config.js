import React, { useState, useEffect } from 'react';
import {
  Card,
  Button,
  message,
  Form,
  Input,
  Space,
  Tag,
  Typography,
  Divider,
  Alert,
  Spin,
} from 'antd';
import {
  SettingOutlined,
  ReloadOutlined,
  LinkOutlined,
  CheckCircleOutlined,
  DisconnectOutlined,
  GlobalOutlined,
  DatabaseOutlined,
} from '@ant-design/icons';
import {
  getConfigs,
  reloadConfigs,
  getDefaultCluster,
  updateCluster,
  testCluster,
} from '../services/api';

const { Title, Text } = Typography;

const Config = ({ darkMode }) => {
  const [configs, setConfigs] = useState({});
  const [cluster, setCluster] = useState(null);
  const [loading, setLoading] = useState(true);
  const [testing, setTesting] = useState(false);
  const [editing, setEditing] = useState(false);
  const [testResult, setTestResult] = useState(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [configsRes, clusterRes] = await Promise.all([
        getConfigs(),
        getDefaultCluster(),
      ]);
      setConfigs(configsRes.data?.data || {});
      const clusterData = clusterRes.data?.data;
      setCluster(clusterData);
      if (clusterData) {
        form.setFieldsValue({
          name: clusterData.name,
          url: clusterData.url,
        });
      }
    } catch (error) {
      console.error('Failed to fetch config:', error);
      message.error('获取配置失败');
    } finally {
      setLoading(false);
    }
  };

  const handleReload = async () => {
    try {
      await reloadConfigs();
      message.success('配置已重新加载');
      fetchData();
    } catch (error) {
      message.error('重新加载配置失败');
    }
  };

  const handleTestConnection = async () => {
    if (!cluster?.id) {
      message.warning('请先配置 Prometheus 地址');
      return;
    }
    try {
      setTesting(true);
      setTestResult(null);
      const res = await testCluster(cluster.id);
      setTestResult({
        success: res.data?.success,
        message: res.data?.message,
      });
      if (res.data?.success) {
        message.success('连接测试成功');
      } else {
        message.error('连接测试失败: ' + res.data?.message);
      }
    } catch (error) {
      message.error('测试连接失败');
      setTestResult({ success: false, message: '网络错误' });
    } finally {
      setTesting(false);
    }
  };

  const handleSave = async (values) => {
    try {
      if (cluster?.id) {
        await updateCluster(cluster.id, { ...values, status: 1 });
      }
      message.success('配置已保存');
      setEditing(false);
      fetchData();
    } catch (error) {
      message.error('保存配置失败');
    }
  };

  const cardStyle = {
    borderRadius: '16px',
    boxShadow: darkMode
      ? '0 4px 20px rgba(0, 0, 0, 0.3)'
      : '0 4px 20px rgba(0, 0, 0, 0.08)',
    background: darkMode ? '#1f1f1f' : '#fff',
    marginBottom: '24px',
  };

  const statusColor = cluster?.status === 1 ? 'success' : 'error';
  const statusText = cluster?.status === 1 ? '已连接' : '未连接';
  const StatusIcon = cluster?.status === 1 ? CheckCircleOutlined : DisconnectOutlined;

  return (
    <div>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '24px',
        }}
      >
        <Title
          level={4}
          style={{ margin: 0, color: darkMode ? '#fff' : '#1d1d1f' }}
        >
          <SettingOutlined style={{ marginRight: '8px' }} />
          系统配置
        </Title>
        <Button icon={<ReloadOutlined />} onClick={handleReload}>
          重新加载配置
        </Button>
      </div>

      {/* Prometheus 配置卡片 */}
      <Card
        title={
          <Space>
            <GlobalOutlined style={{ color: '#1677ff' }} />
            <span style={{ color: darkMode ? '#fff' : '#1d1d1f' }}>
              Prometheus 监控配置
            </span>
          </Space>
        }
        style={cardStyle}
        extra={
          !editing && (
            <Button type="primary" onClick={() => setEditing(true)}>
              编辑配置
            </Button>
          )
        }
      >
        {loading ? (
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <Spin />
          </div>
        ) : editing ? (
          <Form
            form={form}
            layout="vertical"
            onFinish={handleSave}
            style={{ maxWidth: '600px' }}
          >
            <Form.Item
              name="name"
              label="配置名称"
              rules={[{ required: true, message: '请输入配置名称' }]}
            >
              <Input
                prefix={<DatabaseOutlined />}
                placeholder="例如：生产环境 Prometheus"
              />
            </Form.Item>
            <Form.Item
              name="url"
              label="Prometheus 地址"
              rules={[
                { required: true, message: '请输入 Prometheus 地址' },
                {
                  pattern: /^https?:\/\/.+/,
                  message: '请输入正确的 URL 格式，例如：http://localhost:9090',
                },
              ]}
            >
              <Input
                prefix={<LinkOutlined />}
                placeholder="例如：http://localhost:9090"
              />
            </Form.Item>
            <Form.Item>
              <Space>
                <Button type="primary" htmlType="submit">
                  保存配置
                </Button>
                <Button onClick={() => setEditing(false)}>取消</Button>
              </Space>
            </Form.Item>
          </Form>
        ) : cluster ? (
          <div>
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                gap: '16px',
                marginBottom: '24px',
              }}
            >
              <div
                style={{
                  width: '64px',
                  height: '64px',
                  borderRadius: '16px',
                  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                }}
              >
                <GlobalOutlined style={{ fontSize: '32px', color: '#fff' }} />
              </div>
              <div>
                <Title level={5} style={{ margin: 0, marginBottom: '4px' }}>
                  {cluster.name}
                </Title>
                <Space>
                  <Tag color={statusColor} icon={<StatusIcon />}>
                    {statusText}
                  </Tag>
                  <Text type="secondary">ID: {cluster.id}</Text>
                </Space>
              </div>
            </div>

            <Divider style={{ margin: '16px 0' }} />

            <div style={{ marginBottom: '16px' }}>
              <Text type="secondary">Prometheus 地址</Text>
              <div style={{ marginTop: '4px' }}>
                <Text
                  copyable
                  style={{
                    fontSize: '16px',
                    fontFamily: 'monospace',
                    background: darkMode ? '#141414' : '#f5f5f5',
                    padding: '8px 12px',
                    borderRadius: '8px',
                    display: 'inline-block',
                  }}
                >
                  {cluster.url}
                </Text>
              </div>
            </div>

            {testResult && (
              <Alert
                message={testResult.success ? '连接测试成功' : '连接测试失败'}
                description={testResult.message}
                type={testResult.success ? 'success' : 'error'}
                showIcon
                closable
                onClose={() => setTestResult(null)}
                style={{ marginBottom: '16px' }}
              />
            )}

            <Space>
              <Button
                type="primary"
                icon={<CheckCircleOutlined />}
                loading={testing}
                onClick={handleTestConnection}
              >
                测试连接
              </Button>
              <Button icon={<ReloadOutlined />} onClick={fetchData}>
                刷新状态
              </Button>
            </Space>
          </div>
        ) : (
          <Alert
            message="未配置 Prometheus"
            description="请配置 Prometheus 地址以启用监控功能"
            type="warning"
            showIcon
            action={
              <Button type="primary" onClick={() => setEditing(true)}>
                立即配置
              </Button>
            }
          />
        )}
      </Card>

      {/* 平台配置卡片 */}
      <Card
        title={
          <Space>
            <SettingOutlined style={{ color: '#52c41a' }} />
            <span style={{ color: darkMode ? '#fff' : '#1d1d1f' }}>
              平台基础配置
            </span>
          </Space>
        }
        style={cardStyle}
      >
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
            gap: '16px',
          }}
        >
          {Object.entries(configs).map(([key, value]) => (
            <div
              key={key}
              style={{
                padding: '16px',
                background: darkMode ? '#141414' : '#f5f5f5',
                borderRadius: '8px',
              }}
            >
              <Text type="secondary" style={{ fontSize: '12px' }}>
                {key}
              </Text>
              <div style={{ marginTop: '4px' }}>
                <Text
                  copyable
                  style={{
                    fontSize: '14px',
                    fontFamily: 'monospace',
                  }}
                >
                  {value}
                </Text>
              </div>
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
};

export default Config;
