import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  Tag,
  Button,
  Space,
  Modal,
  message,
  Select,
  Badge,
  Spin,
} from 'antd';
import {
  AlertOutlined,
  CheckCircleOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import { getAlerts, acknowledgeAlert } from '../services/api';

const { Option } = Select;

const Alerts = ({ darkMode }) => {
  const [alerts, setAlerts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [statusFilter, setStatusFilter] = useState('firing');
  const [severityFilter, setSeverityFilter] = useState('');
  const [selectedAlert, setSelectedAlert] = useState(null);
  const [detailVisible, setDetailVisible] = useState(false);

  useEffect(() => {
    fetchAlerts();
    const interval = setInterval(fetchAlerts, 30000);
    return () => clearInterval(interval);
  }, [statusFilter, severityFilter]);

  const fetchAlerts = async () => {
    try {
      setLoading(true);
      const params = {};
      if (statusFilter) params.status = statusFilter;
      if (severityFilter) params.severity = severityFilter;

      const response = await getAlerts(params);
      setAlerts(response.data?.data || []);
    } catch (error) {
      console.error('Failed to fetch alerts:', error);
      message.error('Failed to fetch alerts');
    } finally {
      setLoading(false);
    }
  };

  const handleAcknowledge = async (alert) => {
    try {
      await acknowledgeAlert(alert.id, { user: 'admin' });
      message.success('Alert acknowledged');
      fetchAlerts();
    } catch (error) {
      message.error('Failed to acknowledge alert');
    }
  };

  const showDetail = (alert) => {
    setSelectedAlert(alert);
    setDetailVisible(true);
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

  const getStatusBadge = (status) => {
    switch (status) {
      case 'firing':
        return <Badge status="error" text="Firing" />;
      case 'resolved':
        return <Badge status="success" text="Resolved" />;
      case 'acknowledged':
        return <Badge status="warning" text="Acknowledged" />;
      default:
        return <Badge status="default" text={status} />;
    }
  };

  const columns = [
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status) => getStatusBadge(status),
    },
    {
      title: 'Severity',
      dataIndex: 'severity',
      key: 'severity',
      render: (severity) => (
        <Tag color={getSeverityColor(severity)}>{severity.toUpperCase()}</Tag>
      ),
    },
    {
      title: 'Summary',
      dataIndex: 'summary',
      key: 'summary',
      ellipsis: true,
    },
    {
      title: 'Description',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: 'Started',
      dataIndex: 'starts_at',
      key: 'starts_at',
      render: (date) => new Date(date).toLocaleString(),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space>
          <Button
            icon={<EyeOutlined />}
            size="small"
            onClick={() => showDetail(record)}
          >
            View
          </Button>
          {record.status === 'firing' && (
            <Button
              icon={<CheckCircleOutlined />}
              size="small"
              type="primary"
              onClick={() => handleAcknowledge(record)}
            >
              Ack
            </Button>
          )}
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
      <h2 style={{ marginBottom: '24px', color: darkMode ? '#fff' : '#1d1d1f' }}>
        <AlertOutlined style={{ marginRight: '8px' }} />
        告警通知
      </h2>

      <Card
        style={cardStyle}
        extra={
          <Space>
            <Select
              placeholder="Status"
              value={statusFilter}
              onChange={setStatusFilter}
              style={{ width: 120 }}
              allowClear
            >
              <Option value="firing">Firing</Option>
              <Option value="resolved">Resolved</Option>
              <Option value="acknowledged">Acknowledged</Option>
            </Select>
            <Select
              placeholder="Severity"
              value={severityFilter}
              onChange={setSeverityFilter}
              style={{ width: 120 }}
              allowClear
            >
              <Option value="critical">Critical</Option>
              <Option value="warning">Warning</Option>
              <Option value="info">Info</Option>
            </Select>
          </Space>
        }
      >
        <Table
          columns={columns}
          dataSource={alerts}
          rowKey="id"
          loading={loading}
          pagination={{ pageSize: 10 }}
          scroll={{ x: true }}
        />
      </Card>

      <Modal
        title="Alert Details"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={null}
        width={700}
      >
        {selectedAlert && (
          <div>
            <p>
              <strong>ID:</strong> {selectedAlert.id}
            </p>
            <p>
              <strong>Status:</strong> {getStatusBadge(selectedAlert.status)}
            </p>
            <p>
              <strong>Severity:</strong>{' '}
              <Tag color={getSeverityColor(selectedAlert.severity)}>
                {selectedAlert.severity.toUpperCase()}
              </Tag>
            </p>
            <p>
              <strong>Summary:</strong> {selectedAlert.summary}
            </p>
            <p>
              <strong>Description:</strong> {selectedAlert.description}
            </p>
            <p>
              <strong>Started:</strong>{' '}
              {new Date(selectedAlert.starts_at).toLocaleString()}
            </p>
            <p>
              <strong>Fingerprint:</strong>{' '}
              <code>{selectedAlert.fingerprint}</code>
            </p>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default Alerts;
