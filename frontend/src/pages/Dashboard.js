import React, { useState, useEffect } from 'react';
import { Card, Row, Col, Statistic, Badge, Spin } from 'antd';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  AreaChart,
  Area,
} from 'recharts';
import {
  AlertOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  DashboardOutlined,
} from '@ant-design/icons';
import { getAlerts, getDashboardMetrics } from '../services/api';

const Dashboard = ({ darkMode }) => {
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    totalAlerts: 0,
    criticalAlerts: 0,
    warningAlerts: 0,
    healthyServices: 0,
  });
  const [cpuData, setCpuData] = useState([]);
  const [memoryData, setMemoryData] = useState([]);

  useEffect(() => {
    fetchDashboardData();
    const interval = setInterval(fetchDashboardData, 30000);
    return () => clearInterval(interval);
  }, []);

  const fetchDashboardData = async () => {
    try {
      setLoading(true);
      const [alertsResponse, metricsResponse] = await Promise.all([
        getAlerts({ status: 'firing' }),
        getDashboardMetrics(),
      ]);

      const alerts = alertsResponse.data?.data || [];
      const metrics = metricsResponse.data?.data || {};

      setStats({
        totalAlerts: alerts.length,
        criticalAlerts: alerts.filter((a) => a.severity === 'critical').length,
        warningAlerts: alerts.filter((a) => a.severity === 'warning').length,
        healthyServices: 12,
      });

      // Generate sample chart data
      const now = new Date();
      const cpuPoints = [];
      const memoryPoints = [];
      for (let i = 0; i < 20; i++) {
        const time = new Date(now - (19 - i) * 60000);
        cpuPoints.push({
          time: time.toLocaleTimeString(),
          value: 40 + Math.random() * 30,
        });
        memoryPoints.push({
          time: time.toLocaleTimeString(),
          value: 50 + Math.random() * 25,
        });
      }
      setCpuData(cpuPoints);
      setMemoryData(memoryPoints);
    } catch (error) {
      console.error('Failed to fetch dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  const cardStyle = {
    borderRadius: '16px',
    boxShadow: darkMode
      ? '0 4px 20px rgba(0, 0, 0, 0.3)'
      : '0 4px 20px rgba(0, 0, 0, 0.08)',
    background: darkMode ? '#1f1f1f' : '#fff',
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px' }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div>
      <h2 style={{ marginBottom: '24px', color: darkMode ? '#fff' : '#1d1d1f' }}>
        <DashboardOutlined style={{ marginRight: '8px' }} />
        首页
      </h2>

      <Row gutter={[24, 24]}>
        <Col xs={24} sm={12} lg={6}>
          <Card style={cardStyle}>
            <Statistic
              title={<span style={{ color: darkMode ? '#aaa' : '#86868b' }}>Total Alerts</span>}
              value={stats.totalAlerts}
              prefix={<AlertOutlined />}
              valueStyle={{ color: '#1677ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card style={cardStyle}>
            <Statistic
              title={<span style={{ color: darkMode ? '#aaa' : '#86868b' }}>Critical</span>}
              value={stats.criticalAlerts}
              prefix={<ExclamationCircleOutlined />}
              valueStyle={{ color: '#ff3b30' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card style={cardStyle}>
            <Statistic
              title={<span style={{ color: darkMode ? '#aaa' : '#86868b' }}>Warnings</span>}
              value={stats.warningAlerts}
              prefix={<Badge status="warning" />}
              valueStyle={{ color: '#ff9500' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card style={cardStyle}>
            <Statistic
              title={<span style={{ color: darkMode ? '#aaa' : '#86868b' }}>Healthy Services</span>}
              value={stats.healthyServices}
              prefix={<CheckCircleOutlined />}
              valueStyle={{ color: '#34c759' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[24, 24]} style={{ marginTop: '24px' }}>
        <Col xs={24} lg={12}>
          <Card
            title={<span style={{ color: darkMode ? '#fff' : '#1d1d1f' }}>CPU Usage</span>}
            style={cardStyle}
          >
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={cpuData}>
                <defs>
                  <linearGradient id="cpuGradient" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#1677ff" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="#1677ff" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke={darkMode ? '#333' : '#eee'} />
                <XAxis dataKey="time" stroke={darkMode ? '#666' : '#999'} />
                <YAxis stroke={darkMode ? '#666' : '#999'} />
                <Tooltip
                  contentStyle={{
                    background: darkMode ? '#1f1f1f' : '#fff',
                    border: 'none',
                    borderRadius: '8px',
                    boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
                  }}
                />
                <Area
                  type="monotone"
                  dataKey="value"
                  stroke="#1677ff"
                  fillOpacity={1}
                  fill="url(#cpuGradient)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card
            title={<span style={{ color: darkMode ? '#fff' : '#1d1d1f' }}>Memory Usage</span>}
            style={cardStyle}
          >
            <ResponsiveContainer width="100%" height={300}>
              <AreaChart data={memoryData}>
                <defs>
                  <linearGradient id="memoryGradient" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#34c759" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="#34c759" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke={darkMode ? '#333' : '#eee'} />
                <XAxis dataKey="time" stroke={darkMode ? '#666' : '#999'} />
                <YAxis stroke={darkMode ? '#666' : '#999'} />
                <Tooltip
                  contentStyle={{
                    background: darkMode ? '#1f1f1f' : '#fff',
                    border: 'none',
                    borderRadius: '8px',
                    boxShadow: '0 4px 12px rgba(0,0,0,0.15)',
                  }}
                />
                <Area
                  type="monotone"
                  dataKey="value"
                  stroke="#34c759"
                  fillOpacity={1}
                  fill="url(#memoryGradient)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
