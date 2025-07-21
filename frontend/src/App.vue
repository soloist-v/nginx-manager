<template>
  <v-app>
    <v-navigation-drawer
      v-model="drawer"
      app
      :rail="rail"
      @click="rail = false"
    >
      <v-list-item
        prepend-icon="mdi-server"
        title="Nginx Manager"
        subtitle="配置管理工具"
      ></v-list-item>

      <v-divider></v-divider>

      <v-list density="compact" nav>
        <v-list-item
          v-for="item in menuItems"
          :key="item.name"
          :to="item.path"
          :prepend-icon="item.icon"
          :title="item.title"
          color="primary"
        ></v-list-item>
      </v-list>

      <template v-slot:append>
        <div class="pa-2">
          <v-btn
            icon="mdi-chevron-left"
            variant="text"
            @click.stop="rail = !rail"
          ></v-btn>
        </div>
      </template>
    </v-navigation-drawer>

    <v-app-bar app color="primary" dark>
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      
      <v-toolbar-title>
        {{ currentPageTitle }}
      </v-toolbar-title>

      <v-spacer></v-spacer>

      <!-- 连接状态指示器 -->
      <v-chip
        :color="connected ? 'success' : 'error'"
        size="small"
        class="mr-4"
      >
        <v-icon start :icon="connected ? 'mdi-wifi' : 'mdi-wifi-off'"></v-icon>
        {{ connected ? '已连接' : '连接断开' }}
      </v-chip>

      <!-- nginx状态指示器 -->
      <v-chip
        :color="statusColor"
        size="small"
        class="mr-4"
      >
        <v-icon start :icon="isRunning ? 'mdi-play' : 'mdi-stop'"></v-icon>
        {{ statusText }}
      </v-chip>

      <v-btn icon="mdi-refresh" @click="refreshStatus"></v-btn>
    </v-app-bar>

    <v-main>
      <v-container fluid>
        <router-view />
      </v-container>
    </v-main>

    <!-- 全局通知 -->
    <v-snackbar
      v-model="snackbar.show"
      :color="snackbar.color"
      :timeout="snackbar.timeout"
      location="top right"
    >
      {{ snackbar.text }}
      <template v-slot:actions>
        <v-btn
          color="white"
          variant="text"
          @click="snackbar.show = false"
        >
          关闭
        </v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, provide } from 'vue'
import { useRoute } from 'vue-router'
import { useNginxStore } from '@/stores/nginx'

const route = useRoute()
const nginxStore = useNginxStore()

// 抽屉和导航
const drawer = ref(true)
const rail = ref(false)

// 菜单项
const menuItems = [
  { name: 'Dashboard', path: '/', title: '仪表盘', icon: 'mdi-view-dashboard' },
  { name: 'ConfigEditor', path: '/config', title: '配置编辑', icon: 'mdi-file-document-edit' },
  { name: 'LogViewer', path: '/logs', title: '日志查看', icon: 'mdi-text-box-search' },
  { name: 'BackupManager', path: '/backup', title: '备份管理', icon: 'mdi-backup-restore' }
]

// 全局通知
const snackbar = ref({
  show: false,
  text: '',
  color: 'info',
  timeout: 4000
})

// 计算属性
const currentPageTitle = computed(() => {
  const item = menuItems.find(item => item.path === route.path)
  return item ? item.title : 'Nginx Manager'
})

const isRunning = computed(() => nginxStore.isRunning)
const statusText = computed(() => nginxStore.statusText)
const statusColor = computed(() => nginxStore.statusColor)
const connected = computed(() => nginxStore.connected)

// 方法
const refreshStatus = () => {
  nginxStore.fetchStatus()
}

const showNotification = (text, color = 'info', timeout = 4000) => {
  snackbar.value = {
    show: true,
    text,
    color,
    timeout
  }
}

// 提供全局通知方法
provide('showNotification', showNotification)

// 生命周期
onMounted(() => {
  nginxStore.connectWebSocket()
  nginxStore.fetchStatus()
})

onUnmounted(() => {
  nginxStore.disconnectWebSocket()
})
</script>

<style scoped>
/* 只保留最基本的样式 */
</style> 