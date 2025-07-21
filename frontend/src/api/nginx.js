import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

// 响应拦截器
api.interceptors.response.use(
  response => response.data,
  error => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export const nginxAPI = {
  // 获取nginx状态
  getStatus() {
    return api.get('/nginx/status')
  },

  // 启动nginx
  start() {
    return api.post('/nginx/start')
  },

  // 停止nginx
  stop() {
    return api.post('/nginx/stop')
  },

  // 重启nginx
  restart() {
    return api.post('/nginx/restart')
  },

  // 重新加载配置
  reload() {
    return api.post('/nginx/reload')
  }
}

export const configAPI = {
  // 获取配置文件内容
  getConfig() {
    return api.get('/config')
  },

  // 保存配置文件
  saveConfig(content) {
    return api.put('/config', { content })
  },

  // 验证配置文件
  validateConfig(content) {
    return api.post('/config/validate', { content })
  },

  // 获取配置模板
  getTemplate() {
    return api.get('/config/template')
  }
}

export const backupAPI = {
  // 获取备份列表
  getBackups() {
    return api.get('/backup')
  },

  // 恢复备份
  restoreBackup(backupId) {
    return api.post(`/backup/restore/${backupId}`)
  },

  // 删除备份
  deleteBackup(backupId) {
    return api.delete(`/backup/${backupId}`)
  }
}
