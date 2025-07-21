<template>
  <div>
    <v-row class="mb-4">
      <v-col cols="12">
        <v-card>
          <v-card-title class="d-flex justify-space-between align-center">
            <div>
              <v-icon class="mr-2">mdi-backup-restore</v-icon>
              配置文件备份管理
            </div>
            <v-btn
              color="primary"
              @click="refreshBackups"
              :loading="loading"
            >
              <v-icon class="mr-1">mdi-refresh</v-icon>
              刷新列表
            </v-btn>
          </v-card-title>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title>
            备份列表
            <v-chip size="small" class="ml-2">{{ backups.length }} 个备份</v-chip>
          </v-card-title>

          <v-data-table
            :headers="headers"
            :items="backups"
            :loading="loading"
            class="elevation-1"
            item-key="id"
          >
            <template v-slot:item.created_at="{ item }">
              {{ formatTime(item.created_at) }}
            </template>

            <template v-slot:item.size="{ item }">
              {{ formatFileSize(item.size) }}
            </template>

            <template v-slot:item.actions="{ item }">
              <v-btn
                icon="mdi-restore"
                size="small"
                color="success"
                @click="confirmRestore(item)"
                :disabled="loading"
              >
              </v-btn>
              <v-btn
                icon="mdi-download"
                size="small"
                color="info"
                @click="downloadBackup(item)"
                class="ml-1"
                :disabled="loading"
              >
              </v-btn>
              <v-btn
                icon="mdi-delete"
                size="small"
                color="error"
                @click="confirmDelete(item)"
                class="ml-1"
                :disabled="loading"
              >
              </v-btn>
            </template>

            <template v-slot:no-data>
              <div class="text-center py-4">
                <v-icon size="48" color="grey">mdi-backup-restore</v-icon>
                <div class="text-h6 mt-2">暂无备份文件</div>
                <div class="text-body-2 text-grey">配置文件在保存时会自动创建备份</div>
              </div>
            </template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>

    <!-- 恢复确认对话框 -->
    <v-dialog v-model="restoreDialog" max-width="500">
      <v-card>
        <v-card-title class="d-flex align-center">
          <v-icon color="warning" class="mr-2">mdi-alert</v-icon>
          确认恢复备份
        </v-card-title>
        <v-card-text>
          <v-alert type="warning" class="mb-4">
            恢复备份将覆盖当前的nginx配置文件，此操作不可撤销！
          </v-alert>
          <p><strong>备份文件：</strong>{{ selectedBackup?.filename }}</p>
          <p><strong>创建时间：</strong>{{ formatTime(selectedBackup?.created_at) }}</p>
          <p><strong>文件大小：</strong>{{ formatFileSize(selectedBackup?.size) }}</p>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="restoreDialog = false">取消</v-btn>
          <v-btn
            color="warning"
            @click="restoreBackup"
            :loading="restoring"
          >
            确认恢复
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 删除确认对话框 -->
    <v-dialog v-model="deleteDialog" max-width="500">
      <v-card>
        <v-card-title class="d-flex align-center">
          <v-icon color="error" class="mr-2">mdi-delete</v-icon>
          确认删除备份
        </v-card-title>
        <v-card-text>
          <p>确定要删除以下备份文件吗？</p>
          <p><strong>文件名：</strong>{{ selectedBackup?.filename }}</p>
          <p><strong>创建时间：</strong>{{ formatTime(selectedBackup?.created_at) }}</p>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="deleteDialog = false">取消</v-btn>
          <v-btn
            color="error"
            @click="deleteBackup"
            :loading="deleting"
          >
            确认删除
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, inject } from 'vue'
import { backupAPI } from '@/api/nginx'

const showNotification = inject('showNotification')

// 状态
const loading = ref(false)
const restoring = ref(false)
const deleting = ref(false)
const backups = ref([])
const restoreDialog = ref(false)
const deleteDialog = ref(false)
const selectedBackup = ref(null)

// 表格头部
const headers = [
  { title: '文件名', key: 'filename', sortable: true },
  { title: '创建时间', key: 'created_at', sortable: true },
  { title: '文件大小', key: 'size', sortable: true },
  { title: '操作', key: 'actions', sortable: false, width: 150 }
]

// 生命周期
onMounted(() => {
  loadBackups()
})

// 加载备份列表
const loadBackups = async () => {
  try {
    loading.value = true
    const response = await backupAPI.getBackups()
    if (response.success) {
      backups.value = response.data
    }
  } catch (error) {
    showNotification('加载备份列表失败: ' + error.message, 'error')
  } finally {
    loading.value = false
  }
}

// 刷新备份列表
const refreshBackups = () => {
  loadBackups()
}

// 确认恢复
const confirmRestore = (backup) => {
  selectedBackup.value = backup
  restoreDialog.value = true
}

// 恢复备份
const restoreBackup = async () => {
  if (!selectedBackup.value) return

  try {
    restoring.value = true
    const response = await backupAPI.restoreBackup(selectedBackup.value.id)
    
    if (response.success) {
      showNotification(response.message, 'success')
      restoreDialog.value = false
      selectedBackup.value = null
    } else {
      showNotification(response.message, 'error')
    }
  } catch (error) {
    showNotification('恢复备份失败: ' + error.message, 'error')
  } finally {
    restoring.value = false
  }
}

// 确认删除
const confirmDelete = (backup) => {
  selectedBackup.value = backup
  deleteDialog.value = true
}

// 删除备份
const deleteBackup = async () => {
  if (!selectedBackup.value) return

  try {
    deleting.value = true
    const response = await backupAPI.deleteBackup(selectedBackup.value.id)
    
    if (response.success) {
      showNotification(response.message, 'success')
      deleteDialog.value = false
      selectedBackup.value = null
      // 从列表中移除
      const index = backups.value.findIndex(b => b.id === selectedBackup.value?.id)
      if (index > -1) {
        backups.value.splice(index, 1)
      }
    } else {
      showNotification(response.message, 'error')
    }
  } catch (error) {
    showNotification('删除备份失败: ' + error.message, 'error')
  } finally {
    deleting.value = false
  }
}

// 下载备份
const downloadBackup = (backup) => {
  // 模拟下载功能
  const url = `/api/backup/download/${backup.id}`
  const link = document.createElement('a')
  link.href = url
  link.download = backup.filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  
  showNotification('备份下载已开始', 'info')
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return 'N/A'
  return new Date(time).toLocaleString()
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes) return '0 B'
  
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  
  return Math.round((bytes / Math.pow(1024, i)) * 100) / 100 + ' ' + sizes[i]
}
</script>

<style scoped>
.v-data-table {
  border-radius: 8px;
}
</style> 