import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Card,
  Form,
  Input,
  InputNumber,
  Select,
  Switch,
  Button,
  message,
  Typography,
  Space,
  Tag,
  Tooltip,
  Tabs,
  Table,
  Popconfirm,
  Badge,
  Row,
  Col,
  Empty,
  Modal,
  Spin,
  Divider,
  Alert,
} from 'antd';
import {
  SettingOutlined,
  RobotOutlined,
  DashboardOutlined,
  ApiOutlined,
  BellOutlined,
  ToolOutlined,
  InfoCircleOutlined,
  ThunderboltOutlined,
  MessageOutlined,
  SafetyOutlined,
  GlobalOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  EditFilled,
  CheckOutlined,
  CloseOutlined,
  AppstoreOutlined,
  CheckCircleFilled,
  DisconnectOutlined,
  DatabaseOutlined,
  SaveOutlined,
} from '@ant-design/icons';
import {
  getConfigItems,
  updateConfigValue,
  getAIModelConfigs,
  createAIModelConfig,
  updateAIModelConfig,
  deleteAIModelConfig,
  testAIModelConfig,
  setDefaultAIModelConfig,
  getDefaultCluster,
  updateCluster,
  testCluster,
} from '../services/api';

const { Text, Title } = Typography;
const { Option } = Select;
const { TextArea } = Input;

// 分类配置
const CATEGORY_CONFIG = {
  platform: {
    label: '平台设置',
    icon: <SettingOutlined />,
    color: '#1890ff',
    subs: {
      basic: { label: '基础配置', icon: <InfoCircleOutlined /> },
      system: { label: '系统配置', icon: <ToolOutlined /> },
      notification: { label: '通知配置', icon: <BellOutlined /> },
    },
  },
  ai: {
    label: 'AI模型配置',
    icon: <RobotOutlined />,
    color: '#722ed1',
    subs: {
      models: { label: '模型配置', icon: <ThunderboltOutlined /> },
      strategy: { label: '分析策略', icon: <DashboardOutlined /> },
      chat: { label: '对话设置', icon: <MessageOutlined /> },
    },
  },
  monitoring: {
    label: '监控告警',
    icon: <DashboardOutlined />,
    color: '#52c41a',
    subs: {
      prometheus: { label: 'Prometheus', icon: <GlobalOutlined /> },
      alert: { label: '告警配置', icon: <BellOutlined /> },
    },
  },
  integration: {
    label: '集成',
    icon: <ApiOutlined />,
    color: '#fa8c16',
    subs: {
      sso: { label: '单点登录', icon: <SafetyOutlined /> },
      api: { label: 'API设置', icon: <ApiOutlined /> },
      webhook: { label: 'Webhook', icon: <ThunderboltOutlined /> },
    },
  },
};

const ConfigCenter = ({ darkMode }) => {
  const { category } = useParams();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState(['platform', 'basic']);
  const [configItems, setConfigItems] = useState([]);
  const [form] = Form.useForm();
  const [aiConfigs, setAiConfigs] = useState([]);
  // 系统名称编辑状态
  const [systemName, setSystemName] = useState('SRE AI Platform');
  const [editingSystemName, setEditingSystemName] = useState(false);
  const [systemNameInput, setSystemNameInput] = useState('');

  // 组件挂载时从 localStorage 读取缓存（避免闪烁）
  useEffect(() => {
    const cachedName = localStorage.getItem('systemName');
    if (cachedName) {
      setSystemName(cachedName);
    }
  }, []);

  // 根据 URL 参数初始化分类
  useEffect(() => {
    if (category && CATEGORY_CONFIG[category]) {
      const firstSubKey = Object.keys(CATEGORY_CONFIG[category].subs)[0];
      setSelectedCategory([category, firstSubKey]);
    } else {
      navigate('/config/platform');
    }
  }, [category, navigate]);

  useEffect(() => {
    if (selectedCategory[0]) {
      fetchConfigItems();
      fetchAiConfigs();
    }
  }, [selectedCategory]);

  const fetchConfigItems = async () => {
    try {
      setLoading(true);
      const [cat, subCategory] = selectedCategory;
      const res = await getConfigItems(cat, subCategory);
      // 转换后端字段名到前端使用的字段名
      const items = (res.data?.data || []).map(item => ({
        ...item,
        key: item.config_key,
        value: item.config_value,
        type: item.value_type,
      }));
      setConfigItems(items);
      const initialValues = {};
      items.forEach((item) => {
        initialValues[item.key] = item.value;
      });
      form.setFieldsValue(initialValues);
      // 如果是基础配置页面，加载系统名称
      if (cat === 'platform' && subCategory === 'basic') {
        const nameItem = items.find(item => item.key === 'app.name');
        if (nameItem && nameItem.value) {
          setSystemName(nameItem.value);
          setSystemNameInput(nameItem.value);
          // 同步更新 localStorage 和通知其他组件
          localStorage.setItem('systemName', nameItem.value);
          window.dispatchEvent(new CustomEvent('systemNameChanged', {
            detail: { name: nameItem.value }
          }));
        }
      }
    } catch (error) {
      message.error('获取配置失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchAiConfigs = async () => {
    try {
      const res = await getAIModelConfigs();
      setAiConfigs(res.data?.data || []);
    } catch (error) {
      console.error('获取AI配置失败', error);
    }
  };

  // 自动保存 - 当值改变时自动保存
  const handleValuesChange = async (changed) => {
    const key = Object.keys(changed)[0];
    const value = changed[key];

    try {
      await updateConfigValue(key, value);
      message.success('配置已自动保存');
    } catch (error) {
      message.error('保存失败');
    }
  };

  const handleSubTabChange = (key) => {
    setSelectedCategory([category, key]);
    // 清除编辑状态
    setEditingSystemName(false);
  };

  // 处理系统名称编辑
  const handleEditSystemName = () => {
    setEditingSystemName(true);
    setSystemNameInput(systemName);
  };

  // 保存系统名称
  const handleSaveSystemName = async () => {
    if (!systemNameInput.trim()) {
      message.error('系统名称不能为空');
      return;
    }
    try {
      await updateConfigValue('app.name', systemNameInput.trim());
      setSystemName(systemNameInput.trim());
      setEditingSystemName(false);
      message.success('系统名称已更新');
      // 更新本地存储，通知其他组件
      localStorage.setItem('systemName', systemNameInput.trim());
      // 同步更新当前页面显示
      setSystemName(systemNameInput.trim());
      // 触发自定义事件通知其他组件刷新
      window.dispatchEvent(new CustomEvent('systemNameChanged', {
        detail: { name: systemNameInput.trim() }
      }));
    } catch (error) {
      message.error('保存失败');
    }
  };

  // 取消编辑系统名称
  const handleCancelEditSystemName = () => {
    setEditingSystemName(false);
    setSystemNameInput(systemName);
  };

  const renderConfigItem = (item) => {
    const commonProps = {
      key: item.key,
      name: item.key,
      label: (
        <Space>
          <Text strong>{item.description}</Text>
          {item.required && <Tag color="red">必填</Tag>}
          {item.sensitive && <Tag color="orange">敏感</Tag>}
        </Space>
      ),
    };

    switch (item.type) {
      case 'string':
        return (
          <Form.Item {...commonProps}>
            <Input placeholder={`请输入${item.description}`} />
          </Form.Item>
        );
      case 'number':
        return (
          <Form.Item {...commonProps}>
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
        );
      case 'boolean':
        return (
          <Form.Item {...commonProps} valuePropName="checked">
            <Switch checkedChildren="开启" unCheckedChildren="关闭" />
          </Form.Item>
        );
      case 'password':
        return (
          <Form.Item {...commonProps}>
            <Input.Password placeholder={`请输入${item.description}`} />
          </Form.Item>
        );
      case 'json':
        return (
          <Form.Item {...commonProps}>
            <TextArea rows={4} placeholder={`请输入JSON格式的${item.description}`} />
          </Form.Item>
        );
      default:
        if (item.options) {
          const options = JSON.parse(item.options);
          return (
            <Form.Item {...commonProps}>
              <Select placeholder={`请选择${item.description}`}>
                {options.map((opt) => (
                  <Option key={opt} value={opt}>
                    {opt}
                  </Option>
                ))}
              </Select>
            </Form.Item>
          );
        }
        return (
          <Form.Item {...commonProps}>
            <Input placeholder={`请输入${item.description}`} />
          </Form.Item>
        );
    }
  };

  const cardStyle = {
    borderRadius: '16px',
    boxShadow: darkMode ? '0 4px 20px rgba(0, 0, 0, 0.3)' : '0 4px 20px rgba(0, 0, 0, 0.08)',
    background: darkMode ? '#1f1f1f' : '#fff',
  };

  const currentCategory = CATEGORY_CONFIG[selectedCategory[0]];
  const currentSubCategory = currentCategory?.subs[selectedCategory[1]];

  // 构建子 Tab 项
  const subTabItems = currentCategory
    ? Object.entries(currentCategory.subs).map(([subKey, sub]) => ({
        key: subKey,
        label: (
          <span>
            {sub.icon}
            <span style={{ marginLeft: 6 }}>{sub.label}</span>
          </span>
        ),
      }))
    : [];

  return (
    <div>
      {/* 页面标题 */}
      <div style={{ marginBottom: '24px' }}>
        <h2 style={{ margin: 0, color: darkMode ? '#fff' : '#1d1d1f' }}>
          <SettingOutlined style={{ marginRight: '8px' }} />
          配置中心
          {currentCategory && (
            <>
              <span style={{ margin: '0 12px', color: '#999' }}>/</span>
              <span style={{ color: currentCategory.color }}>
                {currentCategory.icon}
                <span style={{ marginLeft: 8 }}>{currentCategory.label}</span>
              </span>
            </>
          )}
        </h2>
      </div>

      {/* 子分类 Tabs */}
      {subTabItems.length > 0 && (
        <Card style={{ ...cardStyle, marginBottom: 16 }} bodyStyle={{ padding: '16px 24px' }}>
          <Tabs
            activeKey={selectedCategory[1]}
            onChange={handleSubTabChange}
            items={subTabItems}
            type="line"
            size="middle"
          />
        </Card>
      )}

      {/* 配置内容区域 */}
      <Spin spinning={loading}>
        <Card
          style={{
            ...cardStyle,
            minHeight: 'calc(100vh - 320px)',
          }}
          title={
            <Space>
              {currentSubCategory?.icon}
              <span>{currentSubCategory?.label}</span>
            </Space>
          }
          extra={
            <Space>
              <Text type="secondary">{configItems.length} 项配置</Text>
            </Space>
          }
        >
          {/* 系统名称特殊编辑区域 - 只在基础配置页面显示 */}
          {selectedCategory[0] === 'platform' && selectedCategory[1] === 'basic' && (
            <>
              <Alert
                message={
                  <div>
                    <div style={{ display: 'flex', alignItems: 'center', marginBottom: 8 }}>
                      <AppstoreOutlined style={{ fontSize: 24, marginRight: 12, color: '#1890ff' }} />
                      <div>
                        <Title level={5} style={{ margin: 0 }}>系统名称</Title>
                        <Text type="secondary">设置平台的显示名称，将显示在登录页和侧边栏</Text>
                      </div>
                    </div>
                    <div style={{ marginTop: 16, padding: '12px 16px', background: darkMode ? '#141414' : '#f5f5f5', borderRadius: 8 }}>
                      {editingSystemName ? (
                        <Space style={{ width: '100%' }}>
                          <Input
                            value={systemNameInput}
                            onChange={(e) => setSystemNameInput(e.target.value)}
                            placeholder="请输入系统名称"
                            style={{ width: 300 }}
                            maxLength={50}
                            showCount
                            autoFocus
                            onPressEnter={handleSaveSystemName}
                          />
                          <Button
                            type="primary"
                            icon={<CheckOutlined />}
                            onClick={handleSaveSystemName}
                          >
                            保存
                          </Button>
                          <Button
                            icon={<CloseOutlined />}
                            onClick={handleCancelEditSystemName}
                          >
                            取消
                          </Button>
                        </Space>
                      ) : (
                        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                          <div>
                            <Text strong style={{ fontSize: 18 }}>{systemName || 'SRE AI Platform'}</Text>
                            <Tag color="blue" style={{ marginLeft: 12 }}>当前名称</Tag>
                          </div>
                          <Button
                            type="primary"
                            icon={<EditFilled />}
                            onClick={handleEditSystemName}
                          >
                            修改系统名称
                          </Button>
                        </div>
                      )}
                    </div>
                  </div>
                }
                type="info"
                showIcon={false}
                style={{ marginBottom: 24, borderRadius: 8 }}
              />
              <Divider style={{ margin: '24px 0' }} />
            </>
          )}

          {selectedCategory[0] === 'ai' && selectedCategory[1] === 'models' ? (
            <AiModelConfigList
              configs={aiConfigs}
              onRefresh={fetchAiConfigs}
              darkMode={darkMode}
            />
          ) : selectedCategory[0] === 'monitoring' && selectedCategory[1] === 'prometheus' ? (
            <PrometheusClusterList darkMode={darkMode} />
          ) : (
            <Form
              form={form}
              layout="vertical"
              onValuesChange={handleValuesChange}
              style={{ maxWidth: 800 }}
            >
              {/* 在基础配置页面，过滤掉app.name，因为已经在上面单独显示 */}
              {configItems
                .filter(item => !(selectedCategory[0] === 'platform' && selectedCategory[1] === 'basic' && item.key === 'app.name'))
                .map((item) => renderConfigItem(item))}

              {configItems.length === 0 && (
                <Empty description="暂无配置项" image={Empty.PRESENTED_IMAGE_SIMPLE} />
              )}
            </Form>
          )}
        </Card>
      </Spin>
    </div>
  );
};

// AI模型配置列表组件
const AiModelConfigList = ({ configs, onRefresh, darkMode }) => {
  const [modalVisible, setModalVisible] = useState(false);
  const [editingConfig, setEditingConfig] = useState(null);
  const [form] = Form.useForm();
  const [testingId, setTestingId] = useState(null);

  const handleCreate = () => {
    setEditingConfig(null);
    form.resetFields();
    form.setFieldsValue({
      provider: 'openai',
      max_tokens: 4000,
      temperature: 0.7,
      timeout: 60,
      is_enabled: true,
      is_default: false,
    });
    setModalVisible(true);
  };

  const handleEdit = (record) => {
    setEditingConfig(record);
    form.setFieldsValue({ ...record });
    setModalVisible(true);
  };

  const handleDelete = async (id) => {
    try {
      await deleteAIModelConfig(id);
      message.success('删除成功');
      onRefresh();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handleTest = async (id) => {
    try {
      setTestingId(id);
      const res = await testAIModelConfig(id);
      if (res.data?.success) {
        message.success(res.data?.message || '连接成功');
      } else {
        message.error(res.data?.message || '连接失败');
      }
    } catch (error) {
      message.error('测试失败');
    } finally {
      setTestingId(null);
    }
  };

  const handleSetDefault = async (id) => {
    try {
      await setDefaultAIModelConfig(id);
      message.success('已设为默认');
      onRefresh();
    } catch (error) {
      message.error('设置失败');
    }
  };

  const handleSubmit = async (values) => {
    try {
      if (editingConfig) {
        await updateAIModelConfig(editingConfig.id, values);
        message.success('更新成功');
      } else {
        await createAIModelConfig(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      onRefresh();
    } catch (error) {
      message.error(editingConfig ? '更新失败' : '创建失败');
    }
  };

  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      render: (text, record) => (
        <Space>
          {text}
          {record.is_default && <Tag color="gold">默认</Tag>}
          {!record.is_enabled && <Tag>已禁用</Tag>}
        </Space>
      ),
    },
    { title: '提供商', dataIndex: 'provider' },
    { title: '模型', dataIndex: 'model' },
    {
      title: '状态',
      render: (_, record) => (
        <Badge status={record.is_enabled ? 'success' : 'error'} text={record.is_enabled ? '启用' : '禁用'} />
      ),
    },
    {
      title: '操作',
      width: 280,
      render: (_, record) => (
        <Space>
          <Button size="small" loading={testingId === record.id} onClick={() => handleTest(record.id)}>
            测试
          </Button>
          {!record.is_default && (
            <Button size="small" onClick={() => handleSetDefault(record.id)}>
              设为默认
            </Button>
          )}
          <Button size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm title="确定删除?" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: '16px' }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          添加模型配置
        </Button>
      </div>
      <Table columns={columns} dataSource={configs} rowKey="id" pagination={false} />

      <Modal
        title={editingConfig ? '编辑模型配置' : '添加模型配置'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={() => form.submit()}
        width={700}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item name="name" label="配置名称" rules={[{ required: true }]}>
            <Input placeholder="例如: OpenAI GPT-4" />
          </Form.Item>
          <Form.Item name="provider" label="提供商" rules={[{ required: true }]}>
            <Select>
              <Option value="openai">OpenAI</Option>
              <Option value="claude">Claude</Option>
              <Option value="azure">Azure OpenAI</Option>
              <Option value="custom">自定义</Option>
            </Select>
          </Form.Item>
          <Form.Item name="model" label="模型" rules={[{ required: true }]}>
            <Input placeholder="例如: gpt-4, claude-3-opus" />
          </Form.Item>
          <Form.Item name="api_key" label="API Key" rules={[{ required: true }]}>
            <Input.Password placeholder="输入API Key" />
          </Form.Item>
          <Form.Item name="base_url" label="Base URL (可选)">
            <Input placeholder="例如: https://api.openai.com/v1" />
          </Form.Item>
          <Row gutter={16}>
            <Col span={8}>
              <Form.Item name="max_tokens" label="Max Tokens">
                <InputNumber min={1} max={32000} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="temperature" label="Temperature">
                <InputNumber min={0} max={2} step={0.1} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="timeout" label="超时(秒)">
                <InputNumber min={1} max={300} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="is_enabled" valuePropName="checked" label="启用">
                <Switch />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="is_default" valuePropName="checked" label="设为默认">
                <Switch />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>
    </div>
  );
};

// Prometheus配置组件 - 极简一行布局
const PrometheusClusterList = ({ darkMode }) => {
  const [cluster, setCluster] = useState(null);
  const [loading, setLoading] = useState(true);
  const [testing, setTesting] = useState(false);
  const [url, setUrl] = useState('');
  const [status, setStatus] = useState(null);

  useEffect(() => {
    fetchCluster();
  }, []);

  const fetchCluster = async () => {
    try {
      setLoading(true);
      const res = await getDefaultCluster();
      const clusterData = res.data?.data;
      setCluster(clusterData);
      if (clusterData) {
        setUrl(clusterData.url);
        setStatus(clusterData.status);
      }
    } catch (error) {
      console.error('获取Prometheus配置失败', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!url.trim()) {
      message.error('请输入 Prometheus 地址');
      return;
    }
    try {
      if (cluster?.id) {
        console.log('Saving cluster:', { id: cluster.id, name: cluster.name, url: url.trim() });
        const response = await updateCluster(cluster.id, {
          name: cluster.name || 'Prometheus',
          url: url.trim(),
          status: 1
        });
        console.log('Save response:', response.data);
        message.success('配置已保存');
        // 延迟一下再获取，确保数据库已更新
        setTimeout(() => {
          fetchCluster();
        }, 500);
      }
    } catch (error) {
      console.error('Save failed:', error);
      message.error('保存失败: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleTest = async () => {
    if (!cluster?.id) {
      message.warning('请先配置 Prometheus 地址');
      return;
    }
    try {
      setTesting(true);
      const res = await testCluster(cluster.id);
      if (res.data?.success) {
        setStatus(1);
        message.success('连接正常');
      } else {
        setStatus(2);
        message.error('连接失败');
      }
    } catch (error) {
      setStatus(2);
      message.error('连接失败');
    } finally {
      setTesting(false);
    }
  };

  if (loading) {
    return <Spin size="small" />;
  }

  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'center',
        gap: '12px',
        padding: '16px 20px',
        background: darkMode ? '#1f1f1f' : '#fff',
        borderRadius: '8px',
        border: `1px solid ${darkMode ? '#303030' : '#e8e8e8'}`,
      }}
    >
      <Input
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        placeholder="http://localhost:9090"
        style={{ width: '320px' }}
        prefix={<GlobalOutlined style={{ color: '#1677ff' }} />}
        onPressEnter={handleSave}
      />
      <Button type="primary" loading={testing} onClick={handleTest}>
        测试连接
      </Button>
      <Button onClick={handleSave}>保存</Button>
      {status === 1 ? (
        <Tag color="success" style={{ margin: 0 }}>连接正常</Tag>
      ) : status === 2 ? (
        <Tag color="error" style={{ margin: 0 }}>连接失败</Tag>
      ) : (
        <Tag style={{ margin: 0 }}>未测试</Tag>
      )}
    </div>
  );
};

export default ConfigCenter;
