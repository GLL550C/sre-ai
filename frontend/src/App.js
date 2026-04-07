import React, { useState } from 'react';
import { Layout, Menu, Switch, ConfigProvider, theme } from 'antd';
import {
  DashboardOutlined,
  AlertOutlined,
  SettingOutlined,
  LineChartOutlined,
  BulbOutlined,
} from '@ant-design/icons';
import { Routes, Route, Link, useLocation } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import Alerts from './pages/Alerts';
import Analysis from './pages/Analysis';
import Rules from './pages/Rules';
import ConfigCenter from './pages/ConfigCenter';

const { Header, Sider, Content } = Layout;

function App() {
  const [darkMode, setDarkMode] = useState(false);
  const location = useLocation();

  const menuItems = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: <Link to="/">Dashboard</Link>,
    },
    {
      key: '/alerts',
      icon: <AlertOutlined />,
      label: <Link to="/alerts">Alerts</Link>,
    },
    {
      key: '/analysis',
      icon: <LineChartOutlined />,
      label: <Link to="/analysis">Analysis</Link>,
    },
    {
      key: 'config',
      icon: <SettingOutlined />,
      label: '配置中心',
      children: [
        {
          key: '/config/platform',
          label: <Link to="/config/platform">平台设置</Link>,
        },
        {
          key: '/config/ai',
          label: <Link to="/config/ai">AI智能</Link>,
        },
        {
          key: '/config/monitoring',
          label: <Link to="/config/monitoring">监控告警</Link>,
        },
        {
          key: '/config/integration',
          label: <Link to="/config/integration">集成</Link>,
        },
      ],
    },
  ];

  const currentTheme = {
    algorithm: darkMode ? theme.darkAlgorithm : theme.defaultAlgorithm,
    token: {
      colorPrimary: '#1677ff',
      borderRadius: 8,
    },
  };

  return (
    <ConfigProvider theme={currentTheme}>
      <Layout style={{ minHeight: '100vh' }}>
        <Sider
          theme={darkMode ? 'dark' : 'light'}
          style={{
            boxShadow: '2px 0 8px rgba(0,0,0,0.1)',
          }}
        >
          <div style={{ padding: '20px', textAlign: 'center' }}>
            <h2 style={{ margin: 0, color: darkMode ? '#fff' : '#1d1d1f' }}>
              SRE AI
            </h2>
          </div>
          <Menu
            theme={darkMode ? 'dark' : 'light'}
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            style={{ border: 'none' }}
          />
        </Sider>
        <Layout>
          <Header
            style={{
              background: darkMode ? '#141414' : '#fff',
              padding: '0 24px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
            }}
          >
            <h3 style={{ margin: 0, color: darkMode ? '#fff' : '#1d1d1f' }}>
              Intelligent Monitoring Platform
            </h3>
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <BulbOutlined style={{ color: darkMode ? '#fff' : '#1d1d1f' }} />
              <Switch
                checked={darkMode}
                onChange={setDarkMode}
                checkedChildren="Dark"
                unCheckedChildren="Light"
              />
            </div>
          </Header>
          <Content
            style={{
              margin: '24px',
              padding: '24px',
              background: darkMode ? '#141414' : '#f5f5f7',
              borderRadius: '16px',
              minHeight: 'calc(100vh - 112px)',
            }}
          >
            <Routes>
              <Route path="/" element={<Dashboard darkMode={darkMode} />} />
              <Route path="/alerts" element={<Alerts darkMode={darkMode} />} />
              <Route path="/analysis" element={<Analysis darkMode={darkMode} />} />
              <Route path="/config" element={<ConfigCenter darkMode={darkMode} />} />
              <Route path="/config/:category" element={<ConfigCenter darkMode={darkMode} />} />
            </Routes>
          </Content>
        </Layout>
      </Layout>
    </ConfigProvider>
  );
}

export default App;
