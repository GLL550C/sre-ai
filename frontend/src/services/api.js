import axios from 'axios';

const API_BASE_URL = '';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor to disable caching
api.interceptors.request.use((config) => {
  config.headers['Cache-Control'] = 'no-cache, no-store, must-revalidate';
  config.headers['Pragma'] = 'no-cache';
  config.headers['Expires'] = '0';
  // Add timestamp to prevent browser caching
  if (config.method === 'get') {
    config.params = {
      ...config.params,
      _t: Date.now(),
    };
  }
  return config;
});

// Alert APIs
export const getAlerts = (params) => api.get('/api/core/alerts', { params });
export const getAlert = (id) => api.get(`/api/core/alerts/${id}`);
export const createAlert = (data) => api.post('/api/core/alerts', data);
export const acknowledgeAlert = (id, data) => api.put(`/api/core/alerts/${id}/ack`, data);

// Rule APIs
export const getRules = () => api.get('/api/core/rules');
export const getRule = (id) => api.get(`/api/core/rules/${id}`);
export const createRule = (data) => api.post('/api/core/rules', data);
export const updateRule = (id, data) => api.put(`/api/core/rules/${id}`, data);
export const deleteRule = (id) => api.delete(`/api/core/rules/${id}`);

// Cluster APIs
export const getClusters = () => api.get('/api/core/clusters');
export const getCluster = (id) => api.get(`/api/core/clusters/${id}`);
export const createCluster = (data) => api.post('/api/core/clusters', data);
export const updateCluster = (id, data) => api.put(`/api/core/clusters/${id}`, data);
export const deleteCluster = (id) => api.delete(`/api/core/clusters/${id}`);
export const testCluster = (id) => api.post(`/api/core/clusters/${id}/test`);
export const setDefaultCluster = (id) => api.put(`/api/core/clusters/${id}/default`);
export const getDefaultCluster = () => api.get('/api/core/clusters/default');

// Analysis APIs
export const getAnalysis = (params) => api.get('/api/core/analysis', { params });
export const getAnalysisById = (id) => api.get(`/api/core/analysis/${id}`);
export const createAnalysis = (data) => api.post('/api/core/analysis', data);
export const deleteAnalysis = (id) => api.delete(`/api/core/analysis/${id}`);
export const archiveAnalysis = (id) => api.put(`/api/core/analysis/${id}/archive`);
export const getAnalysisStats = (params) => api.get('/api/core/analysis/stats', { params });
export const compareClusters = (data) => api.post('/api/core/analysis/compare', data);

// Dashboard APIs
export const getDashboard = (params) => api.get('/api/core/dashboard', { params });
export const getDashboardMetrics = () => api.get('/api/core/dashboard/metrics');

// Config APIs (Legacy)
export const getConfigs = () => api.get('/api/core/configs');
export const reloadConfigs = () => api.post('/api/core/configs/reload');

// New Hierarchical Config APIs
export const getConfigTree = () => api.get('/api/core/config/tree');
export const getConfigItems = (category, subCategory) =>
  api.get('/api/core/config/items', { params: { category, sub_category: subCategory } });
export const getConfigItem = (key) => api.get(`/api/core/config/items/${key}`);
export const updateConfigValue = (key, value) =>
  api.put(`/api/core/config/items/${key}`, { value });
export const updateMultipleConfigs = (configs) =>
  api.post('/api/core/config/batch', { configs });
export const resetConfigToDefault = (key) =>
  api.post(`/api/core/config/items/${key}/reset`);
export const getSystemSettings = () => api.get('/api/core/config/settings/system');
export const getAISettings = () => api.get('/api/core/config/settings/ai');
export const getNotificationSettings = () => api.get('/api/core/config/settings/notification');
export const exportConfig = () => api.get('/api/core/config/export');
export const importConfig = (data) => api.post('/api/core/config/import', data);

// Prometheus APIs
export const queryPrometheus = (params) => api.get('/api/core/prometheus/query', { params });
export const queryPrometheusRange = (params) =>
  api.get('/api/core/prometheus/query_range', { params });

// Runbook APIs
export const getRunbooks = (params) => api.get('/api/runbook/runbooks', { params });
export const getRunbook = (id) => api.get(`/api/runbook/runbooks/${id}`);
export const createRunbook = (data) => api.post('/api/runbook/runbooks', data);
export const updateRunbook = (id, data) => api.put(`/api/runbook/runbooks/${id}`, data);
export const deleteRunbook = (id) => api.delete(`/api/runbook/runbooks/${id}`);
export const searchRunbooks = (params) => api.get('/api/runbook/runbooks/search', { params });

// Tenant APIs
export const getTenants = () => api.get('/api/tenant/tenants');
export const getTenant = (id) => api.get(`/api/tenant/tenants/${id}`);
export const createTenant = (data) => api.post('/api/tenant/tenants', data);
export const updateTenant = (id, data) => api.put(`/api/tenant/tenants/${id}`, data);
export const deleteTenant = (id) => api.delete(`/api/tenant/tenants/${id}`);

// AI Model Config APIs
export const getAIModelConfigs = () => api.get('/api/core/ai/configs');
export const getAIModelConfig = (id) => api.get(`/api/core/ai/configs/${id}`);
export const createAIModelConfig = (data) => api.post('/api/core/ai/configs', data);
export const updateAIModelConfig = (id, data) => api.put(`/api/core/ai/configs/${id}`, data);
export const deleteAIModelConfig = (id) => api.delete(`/api/core/ai/configs/${id}`);
export const testAIModelConfig = (id) => api.post(`/api/core/ai/configs/${id}/test`);
export const setDefaultAIModelConfig = (id) => api.put(`/api/core/ai/configs/${id}/default`);
export const getActiveAIModelConfig = () => api.get('/api/core/ai/configs/active');

// AI Chat APIs
export const chatWithAI = (data) => api.post('/api/core/ai/chat', data);
export const chatWithAIStream = (data) => api.post('/api/core/ai/chat/stream', data);
export const getAIHealth = () => api.get('/api/core/ai/health');
export const getAIModelInfo = () => api.get('/api/core/ai/model');

export default api;