import React, { useState, useEffect } from 'react';
import { Layout, Menu, Switch, ConfigProvider, theme, Avatar, Dropdown, message, Spin } from 'antd';
import {
  DashboardOutlined,
  AlertOutlined,
  SettingOutlined,
  LineChartOutlined,
  BulbOutlined,
  UserOutlined,
  LogoutOutlined,
  LockOutlined,
  TeamOutlined,
} from '@ant-design/icons';
import { Routes, Route, Link, useLocation, useNavigate, Navigate } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import Alerts from './pages/Alerts';
import Analysis from './pages/Analysis';
import ConfigCenter from './pages/ConfigCenter';
import Login from './pages/Login';
import UserManagement from './pages/UserManagement';
import { getSystemName } from './services/api';

const { Header, Sider, Content } = Layout;

function App() {
  const [darkMode, setDarkMode] = useState(false);
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [systemName, setSystemName] = useState('SRE AI Platform');

  // 组件挂载时从 localStorage 读取缓存（避免闪烁）
  useEffect(() => {
    const cachedName = localStorage.getItem('systemName');
    if (cachedName) {
      setSystemName(cachedName);
    }
  }, []);
  const location = useLocation();
  const navigate = useNavigate();

  // 检查登录状态和加载系统名称
  useEffect(() => {
    const checkAuth = () => {
      const storedUser = localStorage.getItem('user');
      const token = localStorage.getItem('token');
      if (storedUser && token) {
        try {
          setUser(JSON.parse(storedUser));
        } catch (e) {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
      }
      setIsLoading(false);
    };
    checkAuth();

    // 加载系统名称
    fetchSystemName();
  }, []);

  // 从后端获取系统名称
  const fetchSystemName = async () => {
    try {
      const res = await getSystemName();
      // 后端返回的数据结构是 { data: "SRE AI Platform" }
      const configValue = res.data?.data;
      if (configValue) {
        setSystemName(configValue);
        localStorage.setItem('systemName', configValue);
      }
    } catch (error) {
      console.error('Failed to fetch system name:', error);
    }
  };

  // 监听系统名称变化事件
  useEffect(() => {
    const handleSystemNameChange = (event) => {
      if (event.detail?.name) {
        setSystemName(event.detail.name);
      }
    };

    window.addEventListener('systemNameChanged', handleSystemNameChange);
    return () => {
      window.removeEventListener('systemNameChanged', handleSystemNameChange);
    };
  }, []);

  // 处理登出
  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setUser(null);
    message.success('已退出登录');
    navigate('/login', { replace: true });
  };

  // 用户菜单
  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人信息',
      onClick: () => {
        message.info('个人信息功能开发中');
      },
    },
    {
      key: 'password',
      icon: <LockOutlined />,
      label: '修改密码',
      onClick: () => {
        message.info('修改密码功能开发中');
      },
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  const menuItems = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: <Link to="/">首页</Link>,
    },
    {
      key: '/alerts',
      icon: <AlertOutlined />,
      label: <Link to="/alerts">告警通知</Link>,
    },
    {
      key: '/analysis',
      icon: <LineChartOutlined />,
      label: <Link to="/analysis">AI分析</Link>,
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
          label: <Link to="/config/ai">AI模型配置</Link>,
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
    ...(user?.role === 'admin'
      ? [
          {
            key: '/users',
            icon: <TeamOutlined />,
            label: <Link to="/users">用户管理</Link>,
          },
        ]
      : []),
  ];

  const currentTheme = {
    algorithm: darkMode ? theme.darkAlgorithm : theme.defaultAlgorithm,
    token: {
      colorPrimary: '#1677ff',
      borderRadius: 8,
    },
  };

  // 加载中
  if (isLoading) {
    return (
      <ConfigProvider theme={currentTheme}>
        <div style={{ height: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <Spin size="large" tip="加载中..." />
        </div>
      </ConfigProvider>
    );
  }

  // 未登录 - 只能访问登录页
  if (!user) {
    return (
      <ConfigProvider theme={currentTheme}>
        <Routes>
          <Route path="/login" element={<Login onLoginSuccess={(userData) => setUser(userData)} />} />
          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </ConfigProvider>
    );
  }

  // 已登录 - 访问登录页时重定向到首页
  if (location.pathname === '/login') {
    return <Navigate to="/" replace />;
  }

  return (
    <ConfigProvider theme={currentTheme}>
      <Layout style={{ minHeight: '100vh' }}>
        <Sider
          theme={darkMode ? 'dark' : 'light'}
          style={{
            boxShadow: '2px 0 8px rgba(0,0,0,0.1)',
          }}
        >
          <Menu
            theme={darkMode ? 'dark' : 'light'}
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            style={{ border: 'none', marginTop: '16px' }}
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
            <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
              <svg
                width="32"
                height="32"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M12 2L2 7L12 12L22 7L12 2Z"
                  stroke={darkMode ? '#fff' : '#1677ff'}
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  fill={darkMode ? 'rgba(255,255,255,0.1)' : 'rgba(22,119,255,0.1)'}
                />
                <path
                  d="M2 17L12 22L22 17"
                  stroke={darkMode ? '#fff' : '#1677ff'}
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
                <path
                  d="M2 12L12 17L22 12"
                  stroke={darkMode ? '#fff' : '#1677ff'}
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
                <circle
                  cx="12"
                  cy="12"
                  r="3"
                  fill={darkMode ? '#fff' : '#1677ff'}
                  fillOpacity="0.9"
                />
              </svg>
              <span
                style={{
                  fontSize: '18px',
                  fontWeight: 600,
                  color: darkMode ? '#fff' : '#1d1d1f',
                }}
              >
                {systemName}
              </span>
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <BulbOutlined style={{ color: darkMode ? '#fff' : '#1d1d1f' }} />
                <Switch
                  checked={darkMode}
                  onChange={setDarkMode}
                  checkedChildren="Dark"
                  unCheckedChildren="Light"
                />
              </div>
              <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
                  <Avatar icon={<UserOutlined />} />
                  <span style={{ color: darkMode ? '#fff' : '#1d1d1f' }}>
                    {user?.username}
                  </span>
                </div>
              </Dropdown>
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
              <Route path="/users" element={<UserManagement />} />
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </Content>
        </Layout>
      </Layout>
    </ConfigProvider>
  );
}

export default App;
