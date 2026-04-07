import React, { useState, useEffect } from 'react';
import { Card, Table, Button, Space, Tag, message, Modal, Form, Input } from 'antd';
import { SettingOutlined, ReloadOutlined, PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { getConfigs, reloadConfigs, getClusters, createCluster, updateCluster, deleteCluster } from '../services/api';

const Config = ({ darkMode }) => {
  const [configs, setConfigs] = useState({});
  const [clusters, setClusters] = useState([]);
  const [loading, setLoading] = useState(true);
  const [clusterModalVisible, setClusterModalVisible] = useState(false);
  const [editingCluster, setEditingCluster] = useState(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [configsRes, clustersRes] = await Promise.all([getConfigs(), getClusters()]);
      setConfigs(configsRes.data?.data || {});
      setClusters(clustersRes.data?.data || []);
    } catch (error) {
      console.error('Failed to fetch config:', error);
      message.error('Failed to fetch configuration');
    } finally {
      setLoading(false);
    }
  };

  const handleReload = async () => {
    try {
      await reloadConfigs();
      message.success('Configuration reloaded');
      fetchData();
    } catch (error) {
      message.error('Failed to reload configuration');
    }
  };

  const handleCreateCluster = () => {
    setEditingCluster(null);
    form.resetFields();
    setClusterModalVisible(true);
  };

  const handleEditCluster = (cluster) => {
    setEditingCluster(cluster);
    form.setFieldsValue({
      name: cluster.name,
      url: cluster.url,
      status: cluster.status,
    });
    setClusterModalVisible(true);
  };

  const handleDeleteCluster = async (id) => {
    try {
      await deleteCluster(id);
      message.success('Cluster deleted');
      fetchData();
    } catch (error) {
      message.error('Failed to delete cluster');
    }
  };

  const handleClusterSubmit = async (values) => {
    try {
      if (editingCluster) {
        await updateCluster(editingCluster.id, values);
        message.success('Cluster updated');
      } else {
        await createCluster(values);
        message.success('Cluster created');
      }
      setClusterModalVisible(false);
      fetchData();
    } catch (error) {
      message.error('Failed to save cluster');
    }
  };

  const configColumns = [
    { title: 'Key', dataIndex: 'key', key: 'key' },
    { title: 'Value', dataIndex: 'value', key: 'value' },
  ];

  const configData = Object.entries(configs).map(([key, value]) => ({
    key,
    value,
  }));

  const clusterColumns = [
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'URL', dataIndex: 'url', key: 'url' },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={status === 1 ? 'success' : 'error'}>
          {status === 1 ? 'Active' : 'Inactive'}
        </Tag>
      ),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button icon={<EditOutlined />} size="small" onClick={() => handleEditCluster(record)}>
            Edit
          </Button>
          <Button icon={<DeleteOutlined />} size="small" danger onClick={() => handleDeleteCluster(record.id)}>
            Delete
          </Button>
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
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <h2 style={{ margin: 0, color: darkMode ? '#fff' : '#1d1d1f' }}>
          <SettingOutlined style={{ marginRight: '8px' }} />
          Configuration
        </h2>
        <Button icon={<ReloadOutlined />} onClick={handleReload}>
          Reload Config
        </Button>
      </div>

      <Card
        title={<span style={{ color: darkMode ? '#fff' : '#1d1d1f' }}>Platform Settings</span>}
        style={cardStyle}
      >
        <Table
          columns={configColumns}
          dataSource={configData}
          rowKey="key"
          loading={loading}
          pagination={false}
          size="small"
        />
      </Card>

      <Card
        title={
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span style={{ color: darkMode ? '#fff' : '#1d1d1f' }}>Prometheus Clusters</span>
            <Button type="primary" icon={<PlusOutlined />} size="small" onClick={handleCreateCluster}>
              Add Cluster
            </Button>
          </div>
        }
        style={cardStyle}
      >
        <Table
          columns={clusterColumns}
          dataSource={clusters}
          rowKey="id"
          loading={loading}
          pagination={{ pageSize: 5 }}
        />
      </Card>

      <Modal
        title={editingCluster ? 'Edit Cluster' : 'Add Cluster'}
        open={clusterModalVisible}
        onCancel={() => setClusterModalVisible(false)}
        onOk={() => form.submit()}
      >
        <Form form={form} layout="vertical" onFinish={handleClusterSubmit}>
          <Form.Item
            name="name"
            label="Name"
            rules={[{ required: true, message: 'Please enter cluster name' }]}
          >
            <Input placeholder="e.g., Production" />
          </Form.Item>
          <Form.Item
            name="url"
            label="URL"
            rules={[{ required: true, message: 'Please enter cluster URL' }]}
          >
            <Input placeholder="e.g., http://prometheus:9090" />
          </Form.Item>
          <Form.Item name="status" hidden>
            <Input type="hidden" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Config;