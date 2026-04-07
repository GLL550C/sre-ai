import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  Button,
  Space,
  Tag,
  message,
  Modal,
  Form,
  Input,
  InputNumber,
  Select,
  Switch,
  Tooltip,
  Popconfirm,
  Badge,
  Descriptions,
  Divider,
  Typography,
  Alert,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  ThunderboltOutlined,
  RobotOutlined,
  GlobalOutlined,
  KeyOutlined,
  SettingOutlined,
  InfoCircleOutlined,
} from '@ant-design/icons';
import {
  getAIModelConfigs,
  createAIModelConfig,
  updateAIModelConfig,
  deleteAIModelConfig,
  testAIModelConfig,
  setDefaultAIModelConfig,
  getActiveAIModelConfig,
} from '../services/api';

const { Title, Text, Paragraph } = Typography;
const { Option } = Select;
const { TextArea } = Input;

const AIModelConfig = ({ darkMode }) => {
  const [configs, setConfigs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalVisible, setModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [selectedConfig, setSelectedConfig] = useState(null);
  const [form] = Form.useForm();
  const [isEditing, setIsEditing] = useState(false);
  const [activeConfig, setActiveConfig] = useState(null);
  const [testingId, setTestingId] = useState(null);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [configsRes, activeRes] = await Promise.all([
        getAIModelConfigs(),
        getActiveAIModelConfig().catch(() => null),
      ]);
      setConfigs(configsRes.data?.data || []);
      if (activeRes?.data?.data) {
        setActiveConfig(activeRes.data.data);
      }
    } catch (error) {
      console.error('Failed to fetch AI configs:', error);
      message.error('Failed to fetch AI configs');
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setIsEditing(false);
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
    setIsEditing(true);
    setSelectedConfig(record);
    form.setFieldsValue({
      ...record,
    });
    setModalVisible(true);
  };

  const handleDelete = async (id) => {
    try {
      await deleteAIModelConfig(id);
      message.success('Config deleted successfully');
      fetchData();
    } catch (error) {
      message.error('Failed to delete config');
    }
  };

  const handleSubmit = async (values) => {
    try {
      if (isEditing && selectedConfig) {
        await updateAIModelConfig(selectedConfig.id, values);
        message.success('Config updated successfully');
      } else {
        await createAIModelConfig(values);
        message.success('Config created successfully');
      }
      setModalVisible(false);
      fetchData();
    } catch (error) {
      message.error(isEditing ? 'Failed to update config' : 'Failed to create config');
    }
  };

  const handleTest = async (id) => {
    try {
      setTestingId(id);
      const res = await testAIModelConfig(id);
      if (res.data?.success) {
        message.success(res.data?.message || 'Connection successful');
      } else {
        message.error(res.data?.message || 'Connection failed');
      }
    } catch (error) {
      message.error('Test failed: ' + (error.response?.data?.error || error.message));
    } finally {
      setTestingId(null);
    }
  };

  const handleSetDefault = async (id) => {
    try {
      await setDefaultAIModelConfig(id);
      message.success('Default config set successfully');
      fetchData();
    } catch (error) {
      message.error('Failed to set default config');
    }
  };

  const handleViewDetail = (record) => {
    setSelectedConfig(record);
    setDetailModalVisible(true);
  };

  const getProviderColor = (provider) => {
    const colors = {
      openai: 'green',
      claude: 'purple',
      azure: 'blue',
      custom: 'orange',
    };
    return colors[provider] || 'default';
  };

  const getProviderLabel = (provider) => {
    const labels = {
      openai: 'OpenAI',
      claude: 'Claude',
      azure: 'Azure OpenAI',
      custom: 'Custom',
    };
    return labels[provider] || provider;
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
      render: (text, record) => (
        <Space>
          <span>{text}</span>
          {record.is_default && (
            <Tag color="gold" icon={<CheckCircleOutlined />}>Default</Tag>
          )}
        </Space>
      ),
    },
    {
      title: 'Provider',
      dataIndex: 'provider',
      key: 'provider',
      render: (provider) => (
        <Tag color={getProviderColor(provider)}>
          {getProviderLabel(provider)}
        </Tag>
      ),
    },
    {
      title: 'Model',
      dataIndex: 'model',
      key: 'model',
    },
    {
      title: 'Status',
      key: 'status',
      render: (_, record) => (
        <Space>
          {record.is_enabled ? (
            <Badge status="success" text="Enabled" />
          ) : (
            <Badge status="error" text="Disabled" />
          )}
        </Space>
      ),
    },
    {
      title: 'Actions',
      key: 'actions',
      width: 280,
      render: (_, record) => (
        <Space size="small">
          <Button
            icon={<ThunderboltOutlined />}
            size="small"
            loading={testingId === record.id}
            onClick={() => handleTest(record.id)}
          >
            Test
          </Button>
          {!record.is_default && (
            <Button
              icon={<CheckCircleOutlined />}
              size="small"
              onClick={() => handleSetDefault(record.id)}
            >
              Set Default
            </Button>
          )}
          <Button
            icon={<EditOutlined />}
            size="small"
            onClick={() => handleEdit(record)}
          >
            Edit
          </Button>
          <Popconfirm
            title="Delete this config?"
            onConfirm={() => handleDelete(record.id)}
          >
            <Button icon={<DeleteOutlined />} size="small" danger>
              Delete
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const cardStyle = {
    borderRadius: '16px',
    boxShadow: darkMode
      ? '0 4px 20px rgba(0, 0, 0, 0.3)'
      : '0 4px 20px rgba(0, 0, 0, 0.08)',
    background: darkMode ? '#1f1f1f' : '#fff',
    marginBottom: '24px',
  };

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
        <h2 style={{ margin: 0, color: darkMode ? '#fff' : '#1d1d1f' }}>
          <RobotOutlined style={{ marginRight: '8px' }} />
          AI Model Configuration
        </h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          Add Config
        </Button>
      </div>

      {/* Active Config Alert */}
      {activeConfig && (
        <Alert
          message={
            <Space>
              <CheckCircleOutlined style={{ color: '#52c41a' }} />
              <span>
                Active AI Model: <strong>{activeConfig.name}</strong> ({getProviderLabel(activeConfig.provider)} - {activeConfig.model})
              </span>
            </Space>
          }
          type="success"
          showIcon={false}
          style={{ marginBottom: '24px', borderRadius: '8px' }}
        />
      )}

      {/* Configs Table */}
      <Card style={cardStyle}>
        <Table
          columns={columns}
          dataSource={configs}
          rowKey="id"
          loading={loading}
          pagination={{ pageSize: 10 }}
        />
      </Card>

      {/* Create/Edit Modal */}
      <Modal
        title={isEditing ? 'Edit AI Model Config' : 'Add AI Model Config'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={() => form.submit()}
        width={700}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="name"
            label="Config Name"
            rules={[{ required: true, message: 'Please enter config name' }]}
          >
            <Input placeholder="e.g., OpenAI GPT-4" />
          </Form.Item>

          <Form.Item
            name="provider"
            label="Provider"
            rules={[{ required: true, message: 'Please select provider' }]}
          >
            <Select placeholder="Select provider">
              <Option value="openai">OpenAI</Option>
              <Option value="claude">Claude (Anthropic)</Option>
              <Option value="azure">Azure OpenAI</Option>
              <Option value="custom">Custom</Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="model"
            label="Model"
            rules={[{ required: true, message: 'Please enter model name' }]}
          >
            <Input placeholder="e.g., gpt-4, claude-3-opus-20240229" />
          </Form.Item>

          <Form.Item
            name="api_key"
            label="API Key"
            rules={[{ required: true, message: 'Please enter API key' }]}
          >
            <Input.Password placeholder="Enter API key" />
          </Form.Item>

          <Form.Item
            name="base_url"
            label="Base URL (Optional)"
          >
            <Input placeholder="e.g., https://api.openai.com/v1" />
          </Form.Item>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '16px' }}>
            <Form.Item
              name="max_tokens"
              label="Max Tokens"
            >
              <InputNumber min={1} max={32000} style={{ width: '100%' }} />
            </Form.Item>

            <Form.Item
              name="temperature"
              label="Temperature"
            >
              <InputNumber min={0} max={2} step={0.1} style={{ width: '100%' }} />
            </Form.Item>

            <Form.Item
              name="timeout"
              label="Timeout (s)"
            >
              <InputNumber min={1} max={300} style={{ width: '100%' }} />
            </Form.Item>
          </div>

          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
            <Form.Item
              name="is_enabled"
              valuePropName="checked"
              label="Enabled"
            >
              <Switch />
            </Form.Item>

            <Form.Item
              name="is_default"
              valuePropName="checked"
              label="Set as Default"
            >
              <Switch />
            </Form.Item>
          </div>

          <Form.Item
            name="description"
            label="Description"
          >
            <TextArea rows={2} placeholder="Optional description" />
          </Form.Item>
        </Form>
      </Modal>

      {/* Detail Modal */}
      <Modal
        title="AI Model Config Details"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailModalVisible(false)}>
            Close
          </Button>,
        ]}
        width={600}
      >
        {selectedConfig && (
          <Descriptions bordered column={1}>
            <Descriptions.Item label="ID">{selectedConfig.id}</Descriptions.Item>
            <Descriptions.Item label="Name">{selectedConfig.name}</Descriptions.Item>
            <Descriptions.Item label="Provider">
              <Tag color={getProviderColor(selectedConfig.provider)}>
                {getProviderLabel(selectedConfig.provider)}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="Model">{selectedConfig.model}</Descriptions.Item>
            <Descriptions.Item label="API Key">
              <Text type="secondary">{selectedConfig.api_key?.substring(0, 10)}****</Text>
            </Descriptions.Item>
            <Descriptions.Item label="Base URL">{selectedConfig.base_url || 'Default'}</Descriptions.Item>
            <Descriptions.Item label="Max Tokens">{selectedConfig.max_tokens}</Descriptions.Item>
            <Descriptions.Item label="Temperature">{selectedConfig.temperature}</Descriptions.Item>
            <Descriptions.Item label="Timeout">{selectedConfig.timeout}s</Descriptions.Item>
            <Descriptions.Item label="Enabled">
              {selectedConfig.is_enabled ? <Badge status="success" text="Yes" /> : <Badge status="error" text="No" />}
            </Descriptions.Item>
            <Descriptions.Item label="Default">
              {selectedConfig.is_default ? <Badge status="success" text="Yes" /> : <Badge status="default" text="No" />}
            </Descriptions.Item>
            <Descriptions.Item label="Description">{selectedConfig.description || '-'}</Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </div>
  );
};

export default AIModelConfig;
