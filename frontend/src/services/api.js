import axios from 'axios';

// API基础URL - 使用相对路径让nginx代理处理
const API_BASE_URL = '/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000, // 30秒超时
});

// Add request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    // 防止缓存
    config.headers['Cache-Control'] = 'no-cache, no-store, must-revalidate';
    config.headers['Pragma'] = 'no-cache';
    config.headers['Expires'] = '0';
    // Add timestamp to prevent browser caching for GET requests
    if (config.method === 'get') {
      config.params = {
        ...config.params,
        _t: Date.now(),
      };
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // 处理网络错误
    if (!error.response) {
      console.error('Network error:', error);
      return Promise.reject(new Error('网络连接失败，请检查网络'));
    }

    // 处理401错误
    if (error.response?.status === 401) {
      // 排除登录接口和获取系统名称接口
      const url = error.config.url;
      if (!url.includes('/auth/login') && !url.includes('/config/items/app.name')) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  }
);

// Auth APIs
export const getCaptcha = () => api.get('/core/auth/captcha');
export const login = (data) => api.post('/core/auth/login', data);
export const logout = () => api.post('/core/auth/logout');
export const getCurrentUser = () => api.get('/core/auth/me');
export const changePassword = (data) => api.post('/core/auth/change-password', data);
export const refreshToken = () => api.post('/core/auth/refresh');

// User Management APIs
export const getUsers = (params) => api.get('/core/users', { params });
export const getUser = (id) => api.get(`/core/users/${id}`);
export const createUser = (data) => api.post('/core/users', data);
export const updateUser = (id, data) => api.put(`/core/users/${id}`, data);
export const deleteUser = (id) => api.delete(`/core/users/${id}`);
export const resetUserPassword = (id, newPassword) =>
  api.post(`/core/users/${id}/reset-password`, { new_password: newPassword });

// Alert APIs
export const getAlerts = (params) => api.get('/core/alerts', { params });
export const getAlert = (id) => api.get(`/core/alerts/${id}`);
export const createAlert = (data) => api.post('/core/alerts', data);
export const acknowledgeAlert = (id, data) => api.put(`/core/alerts/${id}/ack`, data);

// Rule APIs
export const getRules = () => api.get('/core/rules');
export const getRule = (id) => api.get(`/core/rules/${id}`);
export const createRule = (data) => api.post('/core/rules', data);
export const updateRule = (id, data) => api.put(`/core/rules/${id}`, data);
export const deleteRule = (id) => api.delete(`/core/rules/${id}`);

// Cluster APIs
export const getClusters = () => api.get('/core/clusters');
export const getCluster = (id) => api.get(`/core/clusters/${id}`);
export const createCluster = (data) => api.post('/core/clusters', data);
export const updateCluster = (id, data) => api.put(`/core/clusters/${id}`, data);
export const deleteCluster = (id) => api.delete(`/core/clusters/${id}`);
export const testCluster = (id) => api.post(`/core/clusters/${id}/test`);
export const setDefaultCluster = (id) => api.put(`/core/clusters/${id}/default`);
export const getDefaultCluster = () => api.get('/core/clusters/default');

// Analysis APIs
export const getAnalysis = (params) => api.get('/core/analysis', { params });
export const getAnalysisById = (id) => api.get(`/core/analysis/${id}`);
export const createAnalysis = (data) => api.post('/core/analysis', data);
export const deleteAnalysis = (id) => api.delete(`/core/analysis/${id}`);
export const archiveAnalysis = (id) => api.put(`/core/analysis/${id}/archive`);
export const getAnalysisStats = (params) => api.get('/core/analysis/stats', { params });
export const compareClusters = (data) => api.post('/core/analysis/compare', data);

// Dashboard APIs
export const getDashboard = (params) => api.get('/core/dashboard', { params });
export const getDashboardMetrics = () => api.get('/core/dashboard/metrics');

// Config APIs (Legacy)
export const getConfigs = () => api.get('/core/configs');
export const reloadConfigs = () => api.post('/core/configs/reload');

// New Hierarchical Config APIs
export const getConfigTree = () => api.get('/core/config/tree');
export const getConfigItems = (category, subCategory) =>
  api.get('/core/config/items', { params: { category, sub_category: subCategory } });
export const getConfigItem = (key) => api.get(`/core/config/items/${key}`);
export const getConfigValue = (key) => api.get(`/core/config/items/${key}`);
export const updateConfigValue = (key, value) =>
  api.put(`/core/config/items/${key}`, { value });
export const updateMultipleConfigs = (configs) =>
  api.post('/core/config/batch', { configs });
export const resetConfigToDefault = (key) =>
  api.post(`/core/config/items/${key}/reset`);
export const getSystemSettings = () => api.get('/core/config/settings/system');
export const getAISettings = () => api.get('/core/config/settings/ai');
export const getNotificationSettings = () => api.get('/core/config/settings/notification');
export const exportConfig = () => api.get('/core/config/export');
export const importConfig = (data) => api.post('/core/config/import', data);

// Prometheus APIs
export const queryPrometheus = (params) => api.get('/core/prometheus/query', { params });
export const queryPrometheusRange = (params) =>
  api.get('/core/prometheus/query_range', { params });

// Runbook APIs
export const getRunbooks = (params) => api.get('/runbook/runbooks', { params });
export const getRunbook = (id) => api.get(`/runbook/runbooks/${id}`);
export const createRunbook = (data) => api.post('/runbook/runbooks', data);
export const updateRunbook = (id, data) => api.put(`/runbook/runbooks/${id}`, data);
export const deleteRunbook = (id) => api.delete(`/runbook/runbooks/${id}`);
export const searchRunbooks = (params) => api.get('/runbook/runbooks/search', { params });

// Tenant APIs
export const getTenants = () => api.get('/tenant/tenants');
export const getTenant = (id) => api.get(`/tenant/tenants/${id}`);
export const createTenant = (data) => api.post('/tenant/tenants', data);
export const updateTenant = (id, data) => api.put(`/tenant/tenants/${id}`, data);
export const deleteTenant = (id) => api.delete(`/tenant/tenants/${id}`);

// AI Model Config APIs
export const getAIModelConfigs = () => api.get('/core/ai/configs');
export const getAIModelConfig = (id) => api.get(`/core/ai/configs/${id}`);
export const createAIModelConfig = (data) => api.post('/core/ai/configs', data);
export const updateAIModelConfig = (id, data) => api.put(`/core/ai/configs/${id}`, data);
export const deleteAIModelConfig = (id) => api.delete(`/core/ai/configs/${id}`);
export const testAIModelConfig = (id) => api.post(`/core/ai/configs/${id}/test`);
export const setDefaultAIModelConfig = (id) => api.put(`/core/ai/configs/${id}/default`);
export const getActiveAIModelConfig = () => api.get('/core/ai/configs/active');

// AI Chat APIs
export const chatWithAI = (data) => api.post('/core/ai/chat', data);
export const chatWithAIStream = (data) => api.post('/core/ai/chat/stream', data);
export const getAIHealth = () => api.get('/core/ai/health');
export const getAIModelInfo = () => api.get('/core/ai/model');

export default api;
