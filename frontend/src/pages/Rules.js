import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  message,
  Popconfirm,
} from 'antd';
import { FileTextOutlined, PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { getRules, createRule, updateRule, deleteRule } from '../services/api';

const { Option } = Select;
const { TextArea } = Input;

const Rules = ({ darkMode }) => {
  const [rules, setRules] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingRule, setEditingRule] = useState(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchRules();
  }, []);

  const fetchRules = async () => {
    try {
      setLoading(true);
      const response = await getRules();
      setRules(response.data?.data || []);
    } catch (error) {
      console.error('Failed to fetch rules:', error);
      message.error('Failed to fetch rules');
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingRule(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (rule) => {
    setEditingRule(rule);
    form.setFieldsValue({
      name: rule.name,
      description: rule.description,
      expr: rule.expr,
      duration: rule.duration,
      severity: rule.severity,
    });
    setModalVisible(true);
  };

  const handleDelete = async (id) => {
    try {
      await deleteRule(id);
      message.success('Rule deleted');
      fetchRules();
    } catch (error) {
      message.error('Failed to delete rule');
    }
  };

  const handleSubmit = async (values) => {
    try {
      if (editingRule) {
        await updateRule(editingRule.id, values);
        message.success('Rule updated');
      } else {
        await createRule(values);
        message.success('Rule created');
      }
      setModalVisible(false);
      fetchRules();
    } catch (error) {
      message.error('Failed to save rule');
    }
  };

  const getSeverityColor = (severity) => {
    switch (severity) {
      case 'critical':
        return 'error';
      case 'warning':
        return 'warning';
      case 'info':
        return 'default';
      default:
        return 'default';
    }
  };

  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Description',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: 'Severity',
      dataIndex: 'severity',
      key: 'severity',
      render: (severity) => (
        <span style={{ color: getSeverityColor(severity) }}>
          {severity.toUpperCase()}
        </span>
      ),
    },
    {
      title: 'Duration',
      dataIndex: 'duration',
      key: 'duration',
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button icon={<EditOutlined />} size="small" onClick={() => handleEdit(record)}>
            Edit
          </Button>
          <Popconfirm
            title="Delete rule?"
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
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <h2 style={{ margin: 0, color: darkMode ? '#fff' : '#1d1d1f' }}>
          <FileTextOutlined style={{ marginRight: '8px' }} />
          Alert Rules
        </h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
          Create Rule
        </Button>
      </div>

      <Card style={cardStyle}>
        <Table
          columns={columns}
          dataSource={rules}
          rowKey="id"
          loading={loading}
          pagination={{ pageSize: 10 }}
        />
      </Card>

      <Modal
        title={editingRule ? 'Edit Rule' : 'Create Rule'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={() => form.submit()}
        width={700}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            name="name"
            label="Name"
            rules={[{ required: true, message: 'Please enter rule name' }]}
          >
            <Input placeholder="e.g., HighCPUUsage" />
          </Form.Item>
          <Form.Item
            name="description"
            label="Description"
            rules={[{ required: true, message: 'Please enter description' }]}
          >
            <TextArea rows={2} placeholder="Rule description" />
          </Form.Item>
          <Form.Item
            name="expr"
            label="Expression (PromQL)"
            rules={[{ required: true, message: 'Please enter expression' }]}
          >
            <TextArea
              rows={3}
              placeholder="e.g., 100 - (avg by (instance) (irate(node_cpu_seconds_total{mode='idle'}[5m])) * 100) > 80"
            />
          </Form.Item>
          <Form.Item
            name="duration"
            label="Duration"
            rules={[{ required: true, message: 'Please enter duration' }]}
          >
            <Input placeholder="e.g., 5m" />
          </Form.Item>
          <Form.Item
            name="severity"
            label="Severity"
            rules={[{ required: true, message: 'Please select severity' }]}
          >
            <Select placeholder="Select severity">
              <Option value="critical">Critical</Option>
              <Option value="warning">Warning</Option>
              <Option value="info">Info</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Rules;