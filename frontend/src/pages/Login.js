import React, { useState, useEffect } from 'react';
import { Form, Input, Button, Card, message, Row, Col, Alert } from 'antd';
import { UserOutlined, LockOutlined, SafetyOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { login, getCaptcha, getSystemName } from '../services/api';

const Login = ({ onLoginSuccess }) => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [captcha, setCaptcha] = useState(null);
  const [errorMsg, setErrorMsg] = useState('');
  const [systemName, setSystemName] = useState('SRE AI Platform');

  // 组件挂载时从 localStorage 读取缓存（避免闪烁）
  useEffect(() => {
    const cachedName = localStorage.getItem('systemName');
    if (cachedName) {
      setSystemName(cachedName);
    }
  }, []);
  const navigate = useNavigate();

  // 加载系统名称
  useEffect(() => {
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

  // 监听系统名称变化
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

  // 获取验证码
  const fetchCaptcha = async () => {
    try {
      const res = await getCaptcha();
      if (res.data?.data) {
        setCaptcha(res.data.data);
        setErrorMsg('');
      }
    } catch (error) {
      console.error('Failed to get captcha:', error);
      setErrorMsg('获取验证码失败，请刷新页面重试');
    }
  };

  useEffect(() => {
    fetchCaptcha();
    // 清除之前的登录状态
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }, []);

  // 处理登录
  const handleSubmit = async (values) => {
    setLoading(true);
    setErrorMsg('');

    try {
      const res = await login({
        username: values.username,
        password: values.password,
        captcha: values.captcha,
        captcha_id: captcha?.id,
      });

      if (res.data?.data) {
        const { token, user } = res.data.data;
        // 保存token和用户信息
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(user));

        // 调用父组件回调更新登录状态
        if (onLoginSuccess) {
          onLoginSuccess(user);
        }

        // 先跳转再显示成功消息
        navigate('/', { replace: true });
        message.success('登录成功', 1.5);
      }
    } catch (error) {
      const errorText = error.response?.data?.error || '登录失败，请检查用户名和密码';
      setErrorMsg(errorText);
      message.error(errorText);
      // 刷新验证码
      fetchCaptcha();
      form.setFieldsValue({ captcha: '' });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        padding: '20px',
      }}
    >
      <Row justify="center" align="middle" style={{ width: '100%' }}>
        <Col xs={24} sm={20} md={16} lg={12} xl={8}>
          <Card
            style={{
              borderRadius: '16px',
              boxShadow: '0 20px 60px rgba(0,0,0,0.3)',
              border: 'none',
              overflow: 'hidden',
            }}
            bodyStyle={{ padding: '40px' }}
          >
            {/* Logo区域 */}
            <div style={{ textAlign: 'center', marginBottom: '40px' }}>
              <div
                style={{
                  width: '100px',
                  height: '100px',
                  margin: '0 auto 20px',
                  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                  borderRadius: '24px',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  boxShadow: '0 8px 24px rgba(102, 126, 234, 0.4)',
                }}
              >
                <svg
                  width="60"
                  height="60"
                  viewBox="0 0 24 24"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M12 2L2 7L12 12L22 7L12 2Z"
                    stroke="white"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    fill="rgba(255,255,255,0.2)"
                  />
                  <path
                    d="M2 17L12 22L22 17"
                    stroke="white"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M2 12L12 17L22 12"
                    stroke="white"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <circle
                    cx="12"
                    cy="12"
                    r="3"
                    fill="white"
                    fillOpacity="0.9"
                  />
                </svg>
              </div>
              <h1
                style={{
                  fontSize: '28px',
                  fontWeight: 'bold',
                  margin: '0 0 8px',
                  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                  WebkitBackgroundClip: 'text',
                  WebkitTextFillColor: 'transparent',
                }}
              >
                {systemName}
              </h1>
              <p style={{ color: '#666', margin: 0, fontSize: '14px' }}>
                智能运维监控平台
              </p>
            </div>

            {/* 错误提示 */}
            {errorMsg && (
              <Alert
                message={errorMsg}
                type="error"
                showIcon
                style={{ marginBottom: '20px' }}
                closable
                onClose={() => setErrorMsg('')}
              />
            )}

            {/* 登录表单 */}
            <Form
              form={form}
              name="login"
              onFinish={handleSubmit}
              autoComplete="off"
              size="large"
            >
              <Form.Item
                name="username"
                rules={[{ required: true, message: '请输入用户名' }]}
              >
                <Input
                  prefix={<UserOutlined style={{ color: '#bfbfbf' }} />}
                  placeholder="用户名"
                  style={{ borderRadius: '8px', height: '44px' }}
                />
              </Form.Item>

              <Form.Item
                name="password"
                rules={[{ required: true, message: '请输入密码' }]}
              >
                <Input.Password
                  prefix={<LockOutlined style={{ color: '#bfbfbf' }} />}
                  placeholder="密码"
                  style={{ borderRadius: '8px', height: '44px' }}
                />
              </Form.Item>

              <Form.Item
                name="captcha"
                rules={[{ required: true, message: '请输入验证码' }]}
              >
                <Row gutter={8}>
                  <Col flex="auto">
                    <Input
                      prefix={<SafetyOutlined style={{ color: '#bfbfbf' }} />}
                      placeholder="验证码"
                      style={{ borderRadius: '8px', height: '44px' }}
                    />
                  </Col>
                  <Col>
                    <div
                      onClick={!loading ? fetchCaptcha : undefined}
                      style={{
                        cursor: loading ? 'not-allowed' : 'pointer',
                        height: '44px',
                        borderRadius: '8px',
                        overflow: 'hidden',
                        border: '1px solid #d9d9d9',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        minWidth: '120px',
                        background: '#f5f5f5',
                        opacity: loading ? 0.6 : 1,
                      }}
                    >
                      {captcha?.image ? (
                        <img
                          src={captcha.image}
                          alt="captcha"
                          style={{ height: '42px', width: '120px', objectFit: 'contain' }}
                        />
                      ) : (
                        <span style={{ color: '#999', padding: '0 10px', fontSize: '12px' }}>
                          点击刷新
                        </span>
                      )}
                    </div>
                  </Col>
                </Row>
              </Form.Item>

              <Form.Item style={{ marginTop: '24px', marginBottom: '0' }}>
                <Button
                  type="primary"
                  htmlType="submit"
                  loading={loading}
                  block
                  style={{
                    height: '48px',
                    borderRadius: '8px',
                    fontSize: '16px',
                    fontWeight: '500',
                    background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                    border: 'none',
                    boxShadow: '0 4px 12px rgba(102, 126, 234, 0.4)',
                  }}
                >
                  登录
                </Button>
              </Form.Item>
            </Form>

            {/* 底部提示 */}
            <div
              style={{
                marginTop: '24px',
                textAlign: 'center',
                fontSize: '12px',
                color: '#999',
              }}
            >
              请使用您的账号密码登录系统
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Login;
