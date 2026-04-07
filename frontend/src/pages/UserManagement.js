import React, { useState, useEffect } from 'react';
import {
  Table,
  Card,
  Button,
  Modal,
  Form,
  Input,
  Select,
  Tag,
  Space,
  message,
  Popconfirm,
  Row,
  Col,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  KeyOutlined,
  UserOutlined,
  MailOutlined,
  PhoneOutlined,
} from '@ant-design/icons';
import {
  getUsers,
  createUser,
  updateUser,
  deleteUser,
  resetUserPassword,
} from '../services/api';

const { Option } = Select;

const UserManagement = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [resetPwdModalVisible, setResetPwdModalVisible] = useState(false);
  const [editingUser, setEditingUser] = useState(null);
  const [resetPwdUser, setResetPwdUser] = useState(null);
  const [form] = Form.useForm();
  const [resetPwdForm] = Form.useForm();
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });

  // 获取用户列表
  const fetchUsers = async (page = 1, pageSize = 10) => {
    setLoading(true);
    try {
      const res = await getUsers({ page, page_size: pageSize });
      if (res.data?.data) {
        setUsers(res.data.data.items || []);
        setPagination({
          current: page,
          pageSize: pageSize,
          total: res.data.data.total || 0,
        });
      }
    } catch (error) {
      message.error('获取用户列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  // 表格列定义
  const columns = [
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      render: (text) => (
        <Space>
          <UserOutlined />
          {text}
        </Space>
      ),
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
      render: (text) => text || '-',
    },
    {
      title: '手机号',
      dataIndex: 'phone',
      key: 'phone',
      render: (text) => text || '-',
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      render: (role) => {
        const roleMap = {
          admin: { color: 'red', text: '管理员' },
          operator: { color: 'blue', text: '运维人员' },
          viewer: { color: 'green', text: '访客' },
        };
        const config = roleMap[role] || { color: 'default', text: role };
        return <Tag color={config.color}>{config.text}</Tag>;
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={status === 1 ? 'success' : 'default'}>
          {status === 1 ? '正常' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '最后登录',
      dataIndex: 'last_login_at',
      key: 'last_login_at',
      render: (text) => text ? new Date(text).toLocaleString() : '从未登录',
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            icon={<KeyOutlined />}
            onClick={() => handleResetPwd(record)}
          >
            重置密码
          </Button>
          <Popconfirm
            title="确认删除"
            description={`确定要删除用户 "${record.username}" 吗？`}
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  // 处理编辑
  const handleEdit = (record) => {
    setEditingUser(record);
    form.setFieldsValue({
      email: record.email,
      phone: record.phone,
      role: record.role,
      status: record.status,
    });
    setModalVisible(true);
  };

  // 处理新增
  const handleAdd = () => {
    setEditingUser(null);
    form.resetFields();
    setModalVisible(true);
  };

  // 处理提交
  const handleSubmit = async (values) => {
    try {
      if (editingUser) {
        // 更新用户
        await updateUser(editingUser.id, {
          email: values.email,
          phone: values.phone,
          role: values.role,
          status: values.status,
        });
        message.success('用户更新成功');
      } else {
        // 创建用户
        await createUser({
          username: values.username,
          password: values.password,
          email: values.email,
          phone: values.phone,
          role: values.role,
        });
        message.success('用户创建成功');
      }
      setModalVisible(false);
      fetchUsers(pagination.current, pagination.pageSize);
    } catch (error) {
      message.error(error.response?.data?.error || '操作失败');
    }
  };

  // 处理删除
  const handleDelete = async (id) => {
    try {
      await deleteUser(id);
      message.success('用户删除成功');
      fetchUsers(pagination.current, pagination.pageSize);
    } catch (error) {
      message.error(error.response?.data?.error || '删除失败');
    }
  };

  // 处理重置密码
  const handleResetPwd = (record) => {
    setResetPwdUser(record);
    resetPwdForm.resetFields();
    setResetPwdModalVisible(true);
  };

  // 提交重置密码
  const handleResetPwdSubmit = async (values) => {
    try {
      await resetUserPassword(resetPwdUser.id, values.new_password);
      message.success('密码重置成功');
      setResetPwdModalVisible(false);
    } catch (error) {
      message.error(error.response?.data?.error || '密码重置失败');
    }
  };

  // 处理分页变化
  const handleTableChange = (newPagination) => {
    fetchUsers(newPagination.current, newPagination.pageSize);
  };

  return (
    <div>
      <Card
        title="用户管理"
        extra={
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={handleAdd}
          >
            新增用户
          </Button>
        }
      >
        <Table
          columns={columns}
          dataSource={users}
          rowKey="id"
          loading={loading}
          pagination={pagination}
          onChange={handleTableChange}
        />
      </Card>

      {/* 新增/编辑用户模态框 */}
      <Modal
        title={editingUser ? '编辑用户' : '新增用户'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={() => form.submit()}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          style={{ marginTop: '20px' }}
        >
          {!editingUser && (
            <>
              <Form.Item
                name="username"
                label="用户名"
                rules={[
                  { required: true, message: '请输入用户名' },
                  { min: 3, message: '用户名至少3个字符' },
                ]}
              >
                <Input
                  prefix={<UserOutlined />}
                  placeholder="请输入用户名"
                />
              </Form.Item>
              <Form.Item
                name="password"
                label="密码"
                rules={[
                  { required: true, message: '请输入密码' },
                  { min: 6, message: '密码至少6个字符' },
                ]}
              >
                <Input.Password placeholder="请输入密码" />
              </Form.Item>
            </>
          )}
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="email"
                label="邮箱"
                rules={[{ type: 'email', message: '请输入有效的邮箱' }]}
              >
                <Input
                  prefix={<MailOutlined />}
                  placeholder="请输入邮箱"
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="phone" label="手机号">
                <Input
                  prefix={<PhoneOutlined />}
                  placeholder="请输入手机号"
                />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="role"
                label="角色"
                rules={[{ required: true, message: '请选择角色' }]}
              >
                <Select placeholder="请选择角色">
                  <Option value="admin">管理员</Option>
                  <Option value="operator">运维人员</Option>
                  <Option value="viewer">访客</Option>
                </Select>
              </Form.Item>
            </Col>
            {editingUser && (
              <Col span={12}>
                <Form.Item
                  name="status"
                  label="状态"
                  rules={[{ required: true, message: '请选择状态' }]}
                >
                  <Select placeholder="请选择状态">
                    <Option value={1}>正常</Option>
                    <Option value={0}>禁用</Option>
                  </Select>
                </Form.Item>
              </Col>
            )}
          </Row>
        </Form>
      </Modal>

      {/* 重置密码模态框 */}
      <Modal
        title="重置密码"
        open={resetPwdModalVisible}
        onCancel={() => setResetPwdModalVisible(false)}
        onOk={() => resetPwdForm.submit()}
      >
        <Form
          form={resetPwdForm}
          layout="vertical"
          onFinish={handleResetPwdSubmit}
          style={{ marginTop: '20px' }}
        >
          <p>
            正在为 <strong>{resetPwdUser?.username}</strong> 重置密码
          </p>
          <Form.Item
            name="new_password"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '密码至少6个字符' },
            ]}
          >
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>
          <Form.Item
            name="confirm_password"
            label="确认密码"
            rules={[
              { required: true, message: '请确认密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('new_password') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'));
                },
              }),
            ]}
          >
            <Input.Password placeholder="请确认密码" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default UserManagement;
