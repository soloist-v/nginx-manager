import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { nginxAPI } from '@/api/nginx'

export const useNginxStore = defineStore('nginx', () => {
  // 状态
  const status = ref({
    is_running: false,
    pid: 0,
    uptime: '',
    version: '',
    config_valid: false,
    updated_at: null
  })

  const loading = ref(false)
  const error = ref(null)
  const websocket = ref(null)
  const connected = ref(false)

  // 计算属性
  const isRunning = computed(() => status.value.is_running)
  const statusText = computed(() => status.value.is_running ? '运行中' : '已停止')
  const statusColor = computed(() => status.value.is_running ? 'success' : 'error')

  // WebSocket连接
  const connectWebSocket = () => {
    if (websocket.value) {
      websocket.value.close()
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws/status`
    
    websocket.value = new WebSocket(wsUrl)
    
    websocket.value.onopen = () => {
      connected.value = true
      console.log('WebSocket connected')
    }
    
    websocket.value.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data)
        if (message.type === 'status') {
          status.value = message.data
        }
      } catch (err) {
        console.error('Failed to parse WebSocket message:', err)
      }
    }
    
    websocket.value.onclose = () => {
      connected.value = false
      console.log('WebSocket disconnected')
      // 重连
      setTimeout(() => {
        if (!connected.value) {
          connectWebSocket()
        }
      }, 5000)
    }
    
    websocket.value.onerror = (err) => {
      console.error('WebSocket error:', err)
    }
  }

  // 断开WebSocket
  const disconnectWebSocket = () => {
    if (websocket.value) {
      websocket.value.close()
      websocket.value = null
    }
    connected.value = false
  }

  // 获取状态
  const fetchStatus = async () => {
    try {
      loading.value = true
      error.value = null
      const response = await nginxAPI.getStatus()
      if (response.success) {
        status.value = response.data
      }
    } catch (err) {
      error.value = err.message || '获取状态失败'
    } finally {
      loading.value = false
    }
  }

  // 启动nginx
  const startNginx = async () => {
    try {
      loading.value = true
      error.value = null
      const response = await nginxAPI.start()
      if (response.success) {
        await fetchStatus()
        return { success: true, message: response.message }
      } else {
        throw new Error(response.message)
      }
    } catch (err) {
      error.value = err.message || '启动失败'
      return { success: false, message: error.value }
    } finally {
      loading.value = false
    }
  }

  // 停止nginx
  const stopNginx = async () => {
    try {
      loading.value = true
      error.value = null
      const response = await nginxAPI.stop()
      if (response.success) {
        await fetchStatus()
        return { success: true, message: response.message }
      } else {
        throw new Error(response.message)
      }
    } catch (err) {
      error.value = err.message || '停止失败'
      return { success: false, message: error.value }
    } finally {
      loading.value = false
    }
  }

  // 重启nginx
  const restartNginx = async () => {
    try {
      loading.value = true
      error.value = null
      const response = await nginxAPI.restart()
      if (response.success) {
        await fetchStatus()
        return { success: true, message: response.message }
      } else {
        throw new Error(response.message)
      }
    } catch (err) {
      error.value = err.message || '重启失败'
      return { success: false, message: error.value }
    } finally {
      loading.value = false
    }
  }

  // 重新加载配置
  const reloadNginx = async () => {
    try {
      loading.value = true
      error.value = null
      const response = await nginxAPI.reload()
      if (response.success) {
        await fetchStatus()
        return { success: true, message: response.message }
      } else {
        throw new Error(response.message)
      }
    } catch (err) {
      error.value = err.message || '重新加载失败'
      return { success: false, message: error.value }
    } finally {
      loading.value = false
    }
  }

  return {
    // 状态
    status,
    loading,
    error,
    connected,
    // 计算属性
    isRunning,
    statusText,
    statusColor,
    // 方法
    connectWebSocket,
    disconnectWebSocket,
    fetchStatus,
    startNginx,
    stopNginx,
    restartNginx,
    reloadNginx
  }
}) 