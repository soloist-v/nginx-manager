<template>
  <div>
    <v-row class="mb-4">
      <v-col cols="12">
        <v-card>
          <v-card-title class="d-flex justify-space-between align-center">
            <div>
              <v-icon class="mr-2">mdi-text-box-search</v-icon>
              Nginx 日志查看器
            </div>
            <div>
              <v-btn-toggle v-model="logType" mandatory>
                <v-btn value="access">
                  <v-icon class="mr-1">mdi-web</v-icon>
                  访问日志
                </v-btn>
                <v-btn value="error">
                  <v-icon class="mr-1">mdi-alert</v-icon>
                  错误日志
                </v-btn>
              </v-btn-toggle>
            </div>
          </v-card-title>
        </v-card>
      </v-col>
    </v-row>

    <v-row class="mb-4">
      <v-col cols="12" md="8">
        <v-text-field
          v-model="searchQuery"
          label="搜索日志内容"
          prepend-inner-icon="mdi-magnify"
          clearable
          variant="outlined"
          @input="filterLogs"
        ></v-text-field>
      </v-col>
      <v-col cols="12" md="4">
        <v-select
          v-model="logLevel"
          :items="logLevels"
          label="日志级别"
          variant="outlined"
          @update:model-value="filterLogs"
        ></v-select>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title class="d-flex justify-space-between align-center">
            <div>
              {{ logType === 'access' ? '访问日志' : '错误日志' }}
              <v-chip size="small" class="ml-2">{{ filteredLogs.length }} 条</v-chip>
            </div>
            <div>
              <v-btn
                variant="outlined"
                class="mr-2"
                @click="refreshLogs"
                :loading="loading"
              >
                <v-icon class="mr-1">mdi-refresh</v-icon>
                刷新
              </v-btn>
              <v-btn
                variant="outlined"
                @click="clearLogs"
              >
                <v-icon class="mr-1">mdi-delete</v-icon>
                清空显示
              </v-btn>
            </div>
          </v-card-title>
          
          <v-divider></v-divider>
          
          <v-card-text class="pa-0">
            <v-virtual-scroll
              :items="filteredLogs"
              height="500"
              item-height="48"
            >
              <template v-slot:default="{ item }">
                <v-list-item
                  :class="getLogItemClass(item)"
                  class="log-item"
                >
                  <template v-slot:prepend>
                    <v-icon :color="getLogLevelColor(item.level)">
                      {{ getLogLevelIcon(item.level) }}
                    </v-icon>
                  </template>
                  
                  <v-list-item-title class="log-content">
                    <span class="log-time">{{ formatTime(item.time) }}</span>
                    <span class="log-message">{{ item.message }}</span>
                  </v-list-item-title>
                  
                  <template v-slot:append>
                    <v-chip
                      :color="getLogLevelColor(item.level)"
                      size="x-small"
                      variant="outlined"
                    >
                      {{ item.level }}
                    </v-chip>
                  </template>
                </v-list-item>
              </template>
            </v-virtual-scroll>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'

// 状态
const loading = ref(false)
const logType = ref('access')
const searchQuery = ref('')
const logLevel = ref('all')
const logs = ref([])

// 日志级别选项
const logLevels = [
  { title: '全部', value: 'all' },
  { title: '信息', value: 'info' },
  { title: '警告', value: 'warn' },
  { title: '错误', value: 'error' },
  { title: '调试', value: 'debug' }
]

// 模拟日志数据
const mockLogs = {
  access: [
    { time: new Date(), level: 'info', message: '127.0.0.1 - - [25/Dec/2024:14:39:02 +0800] "GET / HTTP/1.1" 200 612 "-" "Mozilla/5.0"' },
    { time: new Date(Date.now() - 60000), level: 'info', message: '127.0.0.1 - - [25/Dec/2024:14:38:02 +0800] "GET /api/nginx/status HTTP/1.1" 200 156 "-" "axios/1.6.2"' },
    { time: new Date(Date.now() - 120000), level: 'info', message: '127.0.0.1 - - [25/Dec/2024:14:37:02 +0800] "POST /api/nginx/start HTTP/1.1" 200 89 "-" "axios/1.6.2"' },
    { time: new Date(Date.now() - 180000), level: 'info', message: '127.0.0.1 - - [25/Dec/2024:14:36:02 +0800] "GET /config HTTP/1.1" 200 2048 "-" "Mozilla/5.0"' },
  ],
  error: [
    { time: new Date(), level: 'error', message: '2024/12/25 14:39:02 [error] 1234#0: *1 connect() failed (111: Connection refused) while connecting to upstream' },
    { time: new Date(Date.now() - 300000), level: 'warn', message: '2024/12/25 14:34:02 [warn] 1234#0: server name "localhost" is not defined for IP address' },
    { time: new Date(Date.now() - 600000), level: 'error', message: '2024/12/25 14:29:02 [error] 1234#0: *2 open() "/var/www/html/favicon.ico" failed (2: No such file or directory)' },
  ]
}

// 计算属性
const filteredLogs = computed(() => {
  let filtered = logs.value

  // 按级别过滤
  if (logLevel.value !== 'all') {
    filtered = filtered.filter(log => log.level === logLevel.value)
  }

  // 按搜索词过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(log => 
      log.message.toLowerCase().includes(query)
    )
  }

  return filtered.sort((a, b) => new Date(b.time) - new Date(a.time))
})

// 监听日志类型变化
watch(logType, () => {
  loadLogs()
})

// 生命周期
onMounted(() => {
  loadLogs()
  // 模拟实时日志更新
  startLogPolling()
})

onUnmounted(() => {
  stopLogPolling()
})

let pollingInterval = null

// 加载日志
const loadLogs = () => {
  loading.value = true
  
  // 模拟API调用
  setTimeout(() => {
    logs.value = [...mockLogs[logType.value]]
    loading.value = false
  }, 500)
}

// 刷新日志
const refreshLogs = () => {
  loadLogs()
}

// 清空显示
const clearLogs = () => {
  logs.value = []
}

// 过滤日志
const filterLogs = () => {
  // 触发计算属性重新计算
}

// 开始日志轮询
const startLogPolling = () => {
  pollingInterval = setInterval(() => {
    // 模拟新日志
    if (Math.random() > 0.7) {
      const newLog = {
        time: new Date(),
        level: Math.random() > 0.8 ? 'error' : 'info',
        message: logType.value === 'access' 
          ? `127.0.0.1 - - [${new Date().toISOString()}] "GET /test HTTP/1.1" 200 123 "-" "TestAgent"`
          : `${new Date().toISOString()} [info] Test log message`
      }
      logs.value.unshift(newLog)
      
      // 限制日志数量
      if (logs.value.length > 100) {
        logs.value = logs.value.slice(0, 100)
      }
    }
  }, 5000)
}

// 停止日志轮询
const stopLogPolling = () => {
  if (pollingInterval) {
    clearInterval(pollingInterval)
    pollingInterval = null
  }
}

// 获取日志项样式类
const getLogItemClass = (item) => {
  const classes = ['border-s-4']
  switch (item.level) {
    case 'error':
      classes.push('border-s-error')
      break
    case 'warn':
      classes.push('border-s-warning')
      break
    case 'info':
      classes.push('border-s-info')
      break
    default:
      classes.push('border-s-grey')
  }
  return classes
}

// 获取日志级别颜色
const getLogLevelColor = (level) => {
  switch (level) {
    case 'error': return 'error'
    case 'warn': return 'warning'
    case 'info': return 'info'
    case 'debug': return 'grey'
    default: return 'grey'
  }
}

// 获取日志级别图标
const getLogLevelIcon = (level) => {
  switch (level) {
    case 'error': return 'mdi-alert-circle'
    case 'warn': return 'mdi-alert'
    case 'info': return 'mdi-information'
    case 'debug': return 'mdi-bug'
    default: return 'mdi-circle'
  }
}

// 格式化时间
const formatTime = (time) => {
  return new Date(time).toLocaleString()
}
</script>

<style scoped>
.log-item {
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.log-content {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.log-time {
  color: #666;
  margin-right: 8px;
}

.log-message {
  word-break: break-all;
}

.v-virtual-scroll {
  background-color: #fafafa;
}
</style> 