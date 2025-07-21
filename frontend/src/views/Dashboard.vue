<template>
  <div>
    <v-row>
      <!-- 状态卡片 -->
      <v-col cols="12" md="6" lg="3">
        <v-card>
          <v-card-item>
            <div class="d-flex justify-space-between align-center">
              <div>
                <div class="text-overline mb-1">服务状态</div>
                <div class="text-h6">{{ statusText }}</div>
              </div>
              <v-icon :color="statusColor" size="48">
                {{ isRunning ? 'mdi-play-circle' : 'mdi-stop-circle' }}
              </v-icon>
            </div>
          </v-card-item>
        </v-card>
      </v-col>

      <v-col cols="12" md="6" lg="3">
        <v-card>
          <v-card-item>
            <div class="d-flex justify-space-between align-center">
              <div>
                <div class="text-overline mb-1">进程ID</div>
                <div class="text-h6">{{ status.pid || 'N/A' }}</div>
              </div>
              <v-icon color="info" size="48">mdi-identifier</v-icon>
            </div>
          </v-card-item>
        </v-card>
      </v-col>

      <v-col cols="12" md="6" lg="3">
        <v-card>
          <v-card-item>
            <div class="d-flex justify-space-between align-center">
              <div>
                <div class="text-overline mb-1">版本</div>
                <div class="text-h6">{{ status.version || 'Unknown' }}</div>
              </div>
              <v-icon color="warning" size="48">mdi-information</v-icon>
            </div>
          </v-card-item>
        </v-card>
      </v-col>

      <v-col cols="12" md="6" lg="3">
        <v-card>
          <v-card-item>
            <div class="d-flex justify-space-between align-center">
              <div>
                <div class="text-overline mb-1">配置状态</div>
                <div class="text-h6">{{ status.config_valid ? '有效' : '无效' }}</div>
              </div>
              <v-icon :color="status.config_valid ? 'success' : 'error'" size="48">
                {{ status.config_valid ? 'mdi-check-circle' : 'mdi-alert-circle' }}
              </v-icon>
            </div>
          </v-card-item>
        </v-card>
      </v-col>
    </v-row>

    <!-- 操作按钮 -->
    <v-row class="mt-4">
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-cog</v-icon>
            服务控制
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" sm="6" md="3">
                <v-btn
                  :disabled="isRunning || loading"
                  :loading="loading"
                  color="success"
                  block
                  @click="handleStart"
                >
                  <v-icon class="mr-1">mdi-play</v-icon>
                  启动
                </v-btn>
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <v-btn
                  :disabled="!isRunning || loading"
                  :loading="loading"
                  color="error"
                  block
                  @click="handleStop"
                >
                  <v-icon class="mr-1">mdi-stop</v-icon>
                  停止
                </v-btn>
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <v-btn
                  :disabled="!isRunning || loading"
                  :loading="loading"
                  color="warning"
                  block
                  @click="handleRestart"
                >
                  <v-icon class="mr-1">mdi-restart</v-icon>
                  重启
                </v-btn>
              </v-col>
              <v-col cols="12" sm="6" md="3">
                <v-btn
                  :disabled="!isRunning || loading"
                  :loading="loading"
                  color="info"
                  block
                  @click="handleReload"
                >
                  <v-icon class="mr-1">mdi-reload</v-icon>
                  重载配置
                </v-btn>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 快速链接 -->
    <v-row class="mt-4">
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-link</v-icon>
            快速链接
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" sm="6" md="4">
                <v-btn
                  to="/config"
                  color="primary"
                  variant="outlined"
                  block
                >
                  <v-icon class="mr-1">mdi-file-document-edit</v-icon>
                  编辑配置
                </v-btn>
              </v-col>
              <v-col cols="12" sm="6" md="4">
                <v-btn
                  to="/logs"
                  color="primary"
                  variant="outlined"
                  block
                >
                  <v-icon class="mr-1">mdi-text-box-search</v-icon>
                  查看日志
                </v-btn>
              </v-col>
              <v-col cols="12" sm="6" md="4">
                <v-btn
                  to="/backup"
                  color="primary"
                  variant="outlined"
                  block
                >
                  <v-icon class="mr-1">mdi-backup-restore</v-icon>
                  备份管理
                </v-btn>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 系统信息 -->
    <v-row class="mt-4">
      <v-col cols="12">
        <v-card>
          <v-card-title>
            <v-icon class="mr-2">mdi-information</v-icon>
            系统信息
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <div class="text-subtitle-2 mb-1">最后更新</div>
                <div>{{ formatTime(status.updated_at) }}</div>
              </v-col>
              <v-col cols="12" md="6">
                <div class="text-subtitle-2 mb-1">连接状态</div>
                <v-chip :color="connected ? 'success' : 'error'" size="small">
                  {{ connected ? '已连接' : '连接断开' }}
                </v-chip>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useNginxStore } from '@/stores/nginx'

const nginxStore = useNginxStore()
const showNotification = inject('showNotification')

// 计算属性
const status = computed(() => nginxStore.status)
const loading = computed(() => nginxStore.loading)
const isRunning = computed(() => nginxStore.isRunning)
const statusText = computed(() => nginxStore.statusText)
const statusColor = computed(() => nginxStore.statusColor)
const connected = computed(() => nginxStore.connected)

// 方法
const handleStart = async () => {
  const result = await nginxStore.startNginx()
  showNotification(result.message, result.success ? 'success' : 'error')
}

const handleStop = async () => {
  const result = await nginxStore.stopNginx()
  showNotification(result.message, result.success ? 'success' : 'error')
}

const handleRestart = async () => {
  const result = await nginxStore.restartNginx()
  showNotification(result.message, result.success ? 'success' : 'error')
}

const handleReload = async () => {
  const result = await nginxStore.reloadNginx()
  showNotification(result.message, result.success ? 'success' : 'error')
}

const formatTime = (time) => {
  if (!time) return 'N/A'
  return new Date(time).toLocaleString()
}
</script> 