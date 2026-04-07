import React, { useState, useEffect, useRef } from 'react';
import {
  Card,
  Table,
  Button,
  Space,
  Tag,
  message,
  Modal,
  Form,
  Select,
  Input,
  Row,
  Col,
  Statistic,
  Tooltip,
  Popconfirm,
  Badge,
  Descriptions,
  Divider,
  List,
  Typography,
  Empty,
  Alert,
  Avatar,
  Spin,
  Tabs,
  Drawer,
} from 'antd';
import {
  LineChartOutlined,
  BulbOutlined,
  RobotOutlined,
  DeleteOutlined,
  PlusOutlined,
  EyeOutlined,
  BarChartOutlined,
  ClusterOutlined,
  ReloadOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  SendOutlined,
  MessageOutlined,
  SettingOutlined,
  CloseOutlined,
} from '@ant-design/icons';
import {
  getAnalysis,
  createAnalysis,
  deleteAnalysis,
  archiveAnalysis,
  getAnalysisStats,
  compareClusters,
  getClusters,
  chatWithAI,
  getAIHealth,
  getAIModelInfo,
  getConfigValue,
} from '../services/api';

const { Title, Text, Paragraph } = Typography;
const { Option } = Select;
const { TextArea } = Input;
const { TabPane } = Tabs;

const Analysis = ({ darkMode }) => {
  const [analyses, setAnalyses] = useState([]);
  const [clusters, setClusters] = useState([]);
  const [stats, setStats] = useState({});
  const [loading, setLoading] = useState(true);
  const [modalVisible, setModalVisible] = useState(false);
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [compareModalVisible, setCompareModalVisible] = useState(false);
  const [chatDrawerVisible, setChatDrawerVisible] = useState(false);
  const [selectedAnalysis, setSelectedAnalysis] = useState(null);
  const [compareResult, setCompareResult] = useState(null);
  const [form] = Form.useForm();
  const [compareForm] = Form.useForm();
  const [filters, setFilters] = useState({
    cluster_id: undefined,
    type: undefined,
    status: 1,
  });

  // Chat state
  const [chatMessages, setChatMessages] = useState([]);
  const [chatInput, setChatInput] = useState('');
  const [chatLoading, setChatLoading] = useState(false);
  const [aiEnabled, setAiEnabled] = useState(false);
  const [aiModelInfo, setAiModelInfo] = useState(null);
  const [systemName, setSystemName] = useState('SRE AI');
  const chatEndRef = useRef(null);

  useEffect(() => {
    fetchData();
    checkAIStatus();
    fetchSystemName();
  }, [filters]);

  useEffect(() => {
    scrollToBottom();
  }, [chatMessages]);

  const scrollToBottom = () => {
    chatEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const fetchData = async () => {
    try {
      setLoading(true);
      const [analysisRes, clustersRes, statsRes] = await Promise.all([
        getAnalysis(filters),
        getClusters(),
        getAnalysisStats(filters.cluster_id ? { cluster_id: filters.cluster_id } : {}),
      ]);
      setAnalyses(analysisRes.data?.data || []);
      setClusters(clustersRes.data?.data || []);
      setStats(statsRes.data?.data || {});
    } catch (error) {
      console.error('Failed to fetch analysis data:', error);
      message.error('Failed to fetch analysis data');
    } finally {
      setLoading(false);
    }
  };

  const checkAIStatus = async () => {
    try {
      const [healthRes, modelRes] = await Promise.all([
        getAIHealth().catch(() => ({ data: { status: 'unhealthy' } })),
        getAIModelInfo().catch(() => ({ data: { data: null } })),
      ]);
      setAiEnabled(healthRes.data?.status === 'healthy');
      setAiModelInfo(modelRes.data?.data);
    } catch (error) {
      setAiEnabled(false);
    }
  };

  // 从后端获取系统名称
  const fetchSystemName = async () => {
    try {
      const res = await getConfigValue('app.name');
      if (res.data?.data?.value) {
        setSystemName(res.data.data.value);
      }
    } catch (error) {
      console.error('Failed to fetch system name:', error);
    }
  };

  const handleCreateAnalysis = async (values) => {
    try {
      await createAnalysis({
        ...values,
        alert_fingerprint: `manual-${Date.now()}`,
      });
      message.success('Analysis created successfully');
      setModalVisible(false);
      form.resetFields();
      fetchData();
    } catch (error) {
      message.error('Failed to create analysis');
    }
  };

  const handleDelete = async (id) => {
    try {
      await deleteAnalysis(id);
      message.success('Analysis deleted');
      fetchData();
    } catch (error) {
      message.error('Failed to delete analysis');
    }
  };

  const handleArchive = async (id) => {
    try {
      await archiveAnalysis(id);
      message.success('Analysis archived');
      fetchData();
    } catch (error) {
      message.error('Failed to archive analysis');
    }
  };

  const handleViewDetail = (record) => {
    setSelectedAnalysis(record);
    setDetailModalVisible(true);
  };

  const handleCompareClusters = async (values) => {
    try {
      const res = await compareClusters({ cluster_ids: values.cluster_ids });
      setCompareResult(res.data?.data);
      message.success('Clusters compared successfully');
    } catch (error) {
      message.error('Failed to compare clusters');
    }
  };

  // Chat functions
  const handleOpenChat = () => {
    if (!aiEnabled) {
      message.warning('AI service is not available. Please configure AI model first.');
      return;
    }
    setChatDrawerVisible(true);
    if (chatMessages.length === 0) {
      setChatMessages([
        {
          role: 'assistant',
          content: `Hello! I am your ${systemName} assistant. I can help you analyze system issues, troubleshoot problems, and provide recommendations. How can I help you today?`,
          timestamp: new Date().toISOString(),
        },
      ]);
    }
  };

  const handleSendMessage = async () => {
    if (!chatInput.trim()) return;

    const userMessage = {
      role: 'user',
      content: chatInput,
      timestamp: new Date().toISOString(),
    };

    setChatMessages((prev) => [...prev, userMessage]);
    setChatInput('');
    setChatLoading(true);

    try {
      const messages = chatMessages.map((m) => ({ role: m.role, content: m.content }));
      messages.push({ role: 'user', content: chatInput });

      const res = await chatWithAI({ messages });
      const assistantMessage = {
        role: 'assistant',
        content: res.data?.data?.choices?.[0]?.message?.content || 'Sorry, I could not process your request.',
        timestamp: new Date().toISOString(),
      };

      setChatMessages((prev) => [...prev, assistantMessage]);
    } catch (error) {
      message.error('Failed to get AI response');
    } finally {
      setChatLoading(false);
    }
  };

  const handleQuickAnalyze = async (type, clusterId) => {
    if (!clusterId) {
      message.warning('Please select a cluster first');
      return;
    }

    try {
      const res = await createAnalysis({
        cluster_id: clusterId,
        analysis_type: type,
        analysis_mode: 'realtime',
        alert_fingerprint: `chat-${Date.now()}`,
      });

      const analysis = res.data?.data;
      const assistantMessage = {
        role: 'assistant',
        content: `I've completed a ${getTypeLabel(type)} for the selected cluster.\n\n**Result:** ${analysis?.result || 'Analysis completed'}\n\n${analysis?.root_cause ? `**Root Cause:** ${analysis.root_cause}\n\n` : ''}${analysis?.suggestions?.length ? `**Suggestions:**\n${analysis.suggestions.map((s) => `- ${s}`).join('\n')}` : ''}`,
        timestamp: new Date().toISOString(),
      };

      setChatMessages((prev) => [...prev, assistantMessage]);
      fetchData();
    } catch (error) {
      message.error('Failed to perform analysis');
    }
  };

  const getTypeColor = (type) => {
    const colors = {
      root_cause: 'blue',
      trend: 'green',
      anomaly: 'orange',
      capacity: 'purple',
      correlation: 'cyan',
      recommendation: 'magenta',
    };
    return colors[type] || 'default';
  };

  const getTypeIcon = (type) => {
    switch (type) {
      case 'root_cause':
        return <BulbOutlined />;
      case 'trend':
        return <LineChartOutlined />;
      case 'anomaly':
        return <ExclamationCircleOutlined />;
      case 'capacity':
        return <BarChartOutlined />;
      case 'correlation':
        return <ClusterOutlined />;
      default:
        return <RobotOutlined />;
    }
  };

  const getTypeLabel = (type) => {
    const labels = {
      root_cause: 'Root Cause Analysis',
      trend: 'Trend Analysis',
      anomaly: 'Anomaly Detection',
      capacity: 'Capacity Planning',
      correlation: 'Alert Correlation',
      recommendation: 'Recommendation',
    };
    return labels[type] || type;
  };

  const getConfidenceBadge = (confidence) => {
    if (!confidence) return <Badge status="default" text="N/A" />;
    if (confidence >= 90) return <Badge status="success" text={`${confidence}%`} />;
    if (confidence >= 70) return <Badge status="warning" text={`${confidence}%`} />;
    return <Badge status="error" text={`${confidence}%`} />;
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: 'Type',
      dataIndex: 'analysis_type',
      key: 'analysis_type',
      render: (type) => (
        <Tag icon={getTypeIcon(type)} color={getTypeColor(type)}>
          {getTypeLabel(type)}
        </Tag>
      ),
    },
    {
      title: 'Cluster',
      dataIndex: 'cluster_id',
      key: 'cluster_id',
      render: (id) => {
        const cluster = clusters.find((c) => c.id === id);
        return cluster ? cluster.name : id ? `Cluster ${id}` : 'N/A';
      },
    },
    {
      title: 'Mode',
      dataIndex: 'analysis_mode',
      key: 'analysis_mode',
      render: (mode) => (
        <Tag color={mode === 'realtime' ? 'green' : mode === 'scheduled' ? 'blue' : 'default'}>
          {mode || 'manual'}
        </Tag>
      ),
    },
    {
      title: 'Result Summary',
      dataIndex: 'result',
      key: 'result',
      ellipsis: true,
      render: (text) => (
        <Tooltip title={text}>
          <span>{text?.substring(0, 50)}...</span>
        </Tooltip>
      ),
    },
    {
      title: 'Confidence',
      dataIndex: 'confidence',
      key: 'confidence',
      render: (confidence) => getConfidenceBadge(confidence),
    },
    {
      title: 'Created',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date) => new Date(date).toLocaleString(),
    },
    {
      title: 'Actions',
      key: 'actions',
      width: 200,
      render: (_, record) => (
        <Space>
          <Button icon={<EyeOutlined />} size="small" onClick={() => handleViewDetail(record)}>
            View
          </Button>
          <Popconfirm title="Archive this analysis?" onConfirm={() => handleArchive(record.id)}>
            <Button icon={<CheckCircleOutlined />} size="small">
              Archive
            </Button>
          </Popconfirm>
          <Popconfirm title="Delete this analysis?" onConfirm={() => handleDelete(record.id)}>
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
    boxShadow: darkMode ? '0 4px 20px rgba(0, 0, 0, 0.3)' : '0 4px 20px rgba(0, 0, 0, 0.08)',
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
          AI分析
        </h2>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={fetchData}>
            Refresh
          </Button>
          <Button onClick={() => setCompareModalVisible(true)}>Compare Clusters</Button>
          <Button type="primary" icon={<MessageOutlined />} onClick={handleOpenChat}>
            AI Chat
          </Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
            New Analysis
          </Button>
        </Space>
      </div>

      {/* AI Status Alert */}
      {aiEnabled ? (
        <Alert
          message={
            <Space>
              <CheckCircleOutlined style={{ color: '#52c41a' }} />
              <span>
                AI Service is active - Model: <strong>{aiModelInfo?.model || 'Unknown'}</strong> ({aiModelInfo?.provider || 'Unknown'})
              </span>
            </Space>
          }
          type="success"
          showIcon={false}
          style={{ marginBottom: '24px', borderRadius: '8px' }}
        />
      ) : (
        <Alert
          message={
            <Space>
              <ExclamationCircleOutlined style={{ color: '#ff4d4f' }} />
              <span>AI Service is not available. Please configure AI model in AI Config page.</span>
            </Space>
          }
          type="error"
          showIcon={false}
          style={{ marginBottom: '24px', borderRadius: '8px' }}
        />
      )}

      {/* Statistics Cards */}
      <Row gutter={[16, 16]} style={{ marginBottom: '24px' }}>
        <Col xs={24} sm={12} md={6}>
          <Card style={cardStyle}>
            <Statistic title="Total Analyses" value={stats.total || 0} prefix={<BarChartOutlined />} />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card style={cardStyle}>
            <Statistic title="Avg Confidence" value={stats.avg_confidence || '0%'} prefix={<CheckCircleOutlined />} />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card style={cardStyle}>
            <Statistic title="Root Cause" value={stats.by_type?.root_cause || 0} prefix={<BulbOutlined />} />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card style={cardStyle}>
            <Statistic title="Anomalies" value={stats.by_type?.anomaly || 0} prefix={<ExclamationCircleOutlined />} />
          </Card>
        </Col>
      </Row>

      {/* Filters */}
      <Card style={cardStyle}>
        <Space wrap>
          <Select
            placeholder="Filter by Cluster"
            value={filters.cluster_id}
            onChange={(value) => setFilters({ ...filters, cluster_id: value })}
            style={{ width: 200 }}
            allowClear
          >
            {clusters.map((cluster) => (
              <Option key={cluster.id} value={cluster.id}>
                {cluster.name}
              </Option>
            ))}
          </Select>
          <Select
            placeholder="Filter by Type"
            value={filters.type}
            onChange={(value) => setFilters({ ...filters, type: value })}
            style={{ width: 200 }}
            allowClear
          >
            <Option value="root_cause">Root Cause</Option>
            <Option value="trend">Trend Analysis</Option>
            <Option value="anomaly">Anomaly Detection</Option>
            <Option value="capacity">Capacity Planning</Option>
            <Option value="correlation">Alert Correlation</Option>
          </Select>
          <Select
            placeholder="Filter by Status"
            value={filters.status}
            onChange={(value) => setFilters({ ...filters, status: value })}
            style={{ width: 150 }}
          >
            <Option value={1}>Active</Option>
            <Option value={2}>Archived</Option>
          </Select>
        </Space>
      </Card>

      {/* Analysis Table */}
      <Card style={cardStyle}>
        <Table columns={columns} dataSource={analyses} rowKey="id" loading={loading} pagination={{ pageSize: 10 }} />
      </Card>

      {/* Create Analysis Modal */}
      <Modal
        title="Create New Analysis"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={() => form.submit()}
        width={600}
      >
        <Form form={form} layout="vertical" onFinish={handleCreateAnalysis}>
          <Form.Item name="cluster_id" label="Prometheus Cluster" rules={[{ required: true, message: 'Please select a cluster' }]}>
            <Select placeholder="Select cluster to analyze">
              {clusters.map((cluster) => (
                <Option key={cluster.id} value={cluster.id}>
                  {cluster.name} ({cluster.url})
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item
            name="analysis_type"
            label="Analysis Type"
            rules={[{ required: true, message: 'Please select analysis type' }]}
          >
            <Select placeholder="Select analysis type">
              <Option value="root_cause">Root Cause Analysis</Option>
              <Option value="trend">Trend Analysis</Option>
              <Option value="anomaly">Anomaly Detection</Option>
              <Option value="capacity">Capacity Planning</Option>
              <Option value="correlation">Alert Correlation</Option>
            </Select>
          </Form.Item>
          <Form.Item name="analysis_mode" label="Analysis Mode" initialValue="realtime">
            <Select>
              <Option value="realtime">Real-time</Option>
              <Option value="historical">Historical</Option>
            </Select>
          </Form.Item>
          <Form.Item name="input_data" label="Additional Context (Optional)">
            <TextArea rows={3} placeholder="Enter any additional context for the analysis" />
          </Form.Item>
        </Form>
      </Modal>

      {/* Detail Modal */}
      <Modal
        title="Analysis Details"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailModalVisible(false)}>
            Close
          </Button>,
        ]}
        width={700}
      >
        {selectedAnalysis && (
          <div>
            <Descriptions bordered column={2}>
              <Descriptions.Item label="ID">{selectedAnalysis.id}</Descriptions.Item>
              <Descriptions.Item label="Type">
                <Tag color={getTypeColor(selectedAnalysis.analysis_type)}>{getTypeLabel(selectedAnalysis.analysis_type)}</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="Cluster">
                {clusters.find((c) => c.id === selectedAnalysis.cluster_id)?.name || `Cluster ${selectedAnalysis.cluster_id}`}
              </Descriptions.Item>
              <Descriptions.Item label="Mode">{selectedAnalysis.analysis_mode}</Descriptions.Item>
              <Descriptions.Item label="Confidence">{getConfidenceBadge(selectedAnalysis.confidence)}</Descriptions.Item>
              <Descriptions.Item label="Model">{selectedAnalysis.model_version}</Descriptions.Item>
              <Descriptions.Item label="Created">{new Date(selectedAnalysis.created_at).toLocaleString()}</Descriptions.Item>
              <Descriptions.Item label="Created By">{selectedAnalysis.created_by}</Descriptions.Item>
            </Descriptions>

            <Divider />

            <Title level={5}>Analysis Result</Title>
            <Paragraph style={{ background: darkMode ? '#141414' : '#f5f5f5', padding: '16px', borderRadius: '8px' }}>
              {selectedAnalysis.result}
            </Paragraph>

            {selectedAnalysis.root_cause && (
              <>
                <Title level={5}>Root Cause</Title>
                <Alert message={selectedAnalysis.root_cause} type="warning" showIcon style={{ marginBottom: '16px' }} />
              </>
            )}

            {selectedAnalysis.suggestions && selectedAnalysis.suggestions.length > 0 && (
              <>
                <Title level={5}>Suggestions</Title>
                <List
                  bordered
                  dataSource={selectedAnalysis.suggestions}
                  renderItem={(item) => (
                    <List.Item>
                      <CheckCircleOutlined style={{ marginRight: 8, color: '#52c41a' }} />
                      {item}
                    </List.Item>
                  )}
                />
              </>
            )}
          </div>
        )}
      </Modal>

      {/* Compare Modal */}
      <Modal
        title="Compare Prometheus Clusters"
        open={compareModalVisible}
        onCancel={() => {
          setCompareModalVisible(false);
          setCompareResult(null);
        }}
        footer={null}
        width={800}
      >
        <Form form={compareForm} onFinish={handleCompareClusters}>
          <Form.Item
            name="cluster_ids"
            label="Select Clusters to Compare (min 2)"
            rules={[{ required: true, message: 'Please select at least 2 clusters' }]}
          >
            <Select mode="multiple" placeholder="Select clusters">
              {clusters.map((cluster) => (
                <Option key={cluster.id} value={cluster.id}>
                  {cluster.name}
                </Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit">
              Compare
            </Button>
          </Form.Item>
        </Form>

        {compareResult && (
          <div style={{ marginTop: '24px' }}>
            <Divider />
            <Title level={5}>Comparison Results</Title>
            {compareResult.findings?.length > 0 ? (
              <List
                bordered
                dataSource={compareResult.findings}
                renderItem={(finding) => (
                  <List.Item>
                    <List.Item.Meta title={finding.title} description={finding.description} />
                  </List.Item>
                )}
              />
            ) : (
              <Empty description="No significant differences found" />
            )}
          </div>
        )}
      </Modal>

      {/* AI Chat Drawer */}
      <Drawer
        title={
          <Space>
            <RobotOutlined />
            <span>AI Assistant</span>
            {aiModelInfo && (
              <Tag size="small" color="blue">
                {aiModelInfo.model}
              </Tag>
            )}
          </Space>
        }
        placement="right"
        width={500}
        onClose={() => setChatDrawerVisible(false)}
        open={chatDrawerVisible}
        bodyStyle={{ padding: 0, display: 'flex', flexDirection: 'column', height: '100%' }}
      >
        {/* Quick Actions */}
        <div style={{ padding: '16px', borderBottom: `1px solid ${darkMode ? '#303030' : '#f0f0f0'}` }}>
          <Text type="secondary" style={{ display: 'block', marginBottom: '8px' }}>
            Quick Analysis:
          </Text>
          <Space wrap>
            <Select
              placeholder="Select cluster"
              style={{ width: 150 }}
              onChange={(value) => (window.selectedChatCluster = value)}
              allowClear
            >
              {clusters.map((cluster) => (
                <Option key={cluster.id} value={cluster.id}>
                  {cluster.name}
                </Option>
              ))}
            </Select>
            <Button size="small" icon={<BulbOutlined />} onClick={() => handleQuickAnalyze('root_cause', window.selectedChatCluster)}>
              Root Cause
            </Button>
            <Button size="small" icon={<LineChartOutlined />} onClick={() => handleQuickAnalyze('trend', window.selectedChatCluster)}>
              Trend
            </Button>
            <Button size="small" icon={<ExclamationCircleOutlined />} onClick={() => handleQuickAnalyze('anomaly', window.selectedChatCluster)}>
              Anomaly
            </Button>
          </Space>
        </div>

        {/* Chat Messages */}
        <div
          style={{
            flex: 1,
            overflow: 'auto',
            padding: '16px',
            background: darkMode ? '#141414' : '#f5f5f5',
          }}
        >
          {chatMessages.map((msg, index) => (
            <div
              key={index}
              style={{
                display: 'flex',
                justifyContent: msg.role === 'user' ? 'flex-end' : 'flex-start',
                marginBottom: '16px',
              }}
            >
              <div
                style={{
                  maxWidth: '80%',
                  padding: '12px 16px',
                  borderRadius: '12px',
                  background: msg.role === 'user' ? '#1677ff' : darkMode ? '#1f1f1f' : '#fff',
                  color: msg.role === 'user' ? '#fff' : darkMode ? '#fff' : '#1d1d1f',
                  boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', marginBottom: '4px' }}>
                  <Avatar
                    size="small"
                    icon={msg.role === 'user' ? <SettingOutlined /> : <RobotOutlined />}
                    style={{
                      marginRight: '8px',
                      background: msg.role === 'user' ? '#52c41a' : '#722ed1',
                    }}
                  />
                  <Text style={{ color: msg.role === 'user' ? '#fff' : 'inherit', fontSize: '12px', fontWeight: 'bold' }}>
                    {msg.role === 'user' ? 'You' : 'AI Assistant'}
                  </Text>
                </div>
                <div style={{ whiteSpace: 'pre-wrap' }}>{msg.content}</div>
              </div>
            </div>
          ))}
          {chatLoading && (
            <div style={{ display: 'flex', justifyContent: 'flex-start', marginBottom: '16px' }}>
              <div
                style={{
                  padding: '12px 16px',
                  borderRadius: '12px',
                  background: darkMode ? '#1f1f1f' : '#fff',
                  boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
                }}
              >
                <Spin size="small" />
                <Text style={{ marginLeft: '8px' }}>Thinking...</Text>
              </div>
            </div>
          )}
          <div ref={chatEndRef} />
        </div>

        {/* Chat Input */}
        <div
          style={{
            padding: '16px',
            borderTop: `1px solid ${darkMode ? '#303030' : '#f0f0f0'}`,
            background: darkMode ? '#1f1f1f' : '#fff',
          }}
        >
          <Space.Compact style={{ width: '100%' }}>
            <TextArea
              value={chatInput}
              onChange={(e) => setChatInput(e.target.value)}
              placeholder="Ask me anything about your system..."
              autoSize={{ minRows: 1, maxRows: 4 }}
              onPressEnter={(e) => {
                if (!e.shiftKey) {
                  e.preventDefault();
                  handleSendMessage();
                }
              }}
            />
            <Button type="primary" icon={<SendOutlined />} onClick={handleSendMessage} loading={chatLoading}>
              Send
            </Button>
          </Space.Compact>
          <Text type="secondary" style={{ fontSize: '12px', marginTop: '8px', display: 'block' }}>
            Press Enter to send, Shift+Enter for new line
          </Text>
        </div>
      </Drawer>
    </div>
  );
};

export default Analysis;
