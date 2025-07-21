import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '@/views/Dashboard.vue'
import ConfigEditor from '@/views/ConfigEditor.vue'
import LogViewer from '@/views/LogViewer.vue'
import BackupManager from '@/views/BackupManager.vue'

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard,
    meta: {
      title: '仪表盘',
      icon: 'mdi-view-dashboard'
    }
  },
  {
    path: '/config',
    name: 'ConfigEditor',
    component: ConfigEditor,
    meta: {
      title: '配置编辑',
      icon: 'mdi-file-document-edit'
    }
  },
  {
    path: '/logs',
    name: 'LogViewer',
    component: LogViewer,
    meta: {
      title: '日志查看',
      icon: 'mdi-text-box-search'
    }
  },
  {
    path: '/backup',
    name: 'BackupManager',
    component: BackupManager,
    meta: {
      title: '备份管理',
      icon: 'mdi-backup-restore'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router 