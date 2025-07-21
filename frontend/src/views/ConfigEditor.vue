<template>
  <div>
    <v-row class="mb-4">
      <v-col cols="12">
        <v-card>
          <v-card-title class="d-flex justify-space-between align-center">
            <div>
              <v-icon class="mr-2">mdi-file-document-edit</v-icon>
              Nginx 配置编辑器
              <v-chip size="x-small" variant="outlined" class="ml-2">
                <v-icon size="x-small" class="mr-1">mdi-keyboard</v-icon>
                F1 查看快捷键
              </v-chip>
            </div>
            <div>
              <v-btn
                color="primary"
                variant="outlined"
                class="mr-2"
                @click="loadTemplate"
                :disabled="loading"
              >
                <v-icon class="mr-1">mdi-file-code</v-icon>
                加载模板
              </v-btn>
              <v-btn
                color="info"
                variant="outlined"
                class="mr-2"
                @click="validateConfig"
                :disabled="loading"
              >
                <v-icon class="mr-1">mdi-check-circle</v-icon>
                验证配置
              </v-btn>
                              <v-btn
                  color="success"
                  @click="saveConfig"
                  :disabled="loading"
                  :loading="saving"
                >
                  <v-icon class="mr-1">mdi-content-save</v-icon>
                  保存配置
                </v-btn>
                <v-btn
                  icon="mdi-keyboard"
                  variant="outlined"
                  @click="showShortcutsDialog = true"
                  class="ml-2"
                >
                </v-btn>
            </div>
          </v-card-title>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-text class="pa-0">
            <!-- Monaco编辑器容器 -->
            <div
              ref="editorContainer"
              style="height: 600px; width: 100%;"
            ></div>
          </v-card-text>
          
          <!-- 编辑器状态栏 -->
          <v-card-actions class="pa-2 bg-grey-lighten-5">
            <v-chip size="x-small" variant="text">
              <v-icon size="x-small" class="mr-1">mdi-cursor-text</v-icon>
              行 {{ cursorPosition.line }}，列 {{ cursorPosition.column }}
            </v-chip>
            
            <v-chip size="x-small" variant="text" class="ml-2">
              <v-icon size="x-small" class="mr-1">mdi-text</v-icon>
              {{ documentStats.lines }} 行，{{ documentStats.characters }} 字符
            </v-chip>
            
            <v-chip 
              size="x-small" 
              variant="text" 
              class="ml-2"
              :color="documentStats.modified ? 'warning' : 'success'"
            >
              <v-icon size="x-small" class="mr-1">
                {{ documentStats.modified ? 'mdi-circle' : 'mdi-check-circle' }}
              </v-icon>
              {{ documentStats.modified ? '已修改' : '已保存' }}
            </v-chip>
            
            <v-spacer></v-spacer>
            
            <v-chip size="x-small" variant="text">
              <v-icon size="x-small" class="mr-1">mdi-code-tags</v-icon>
              nginx
            </v-chip>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!-- 验证结果对话框 -->
    <v-dialog v-model="validationDialog" max-width="600">
      <v-card>
        <v-card-title class="d-flex align-center">
          <v-icon :color="validationResult.valid ? 'success' : 'error'" class="mr-2">
            {{ validationResult.valid ? 'mdi-check-circle' : 'mdi-alert-circle' }}
          </v-icon>
          配置验证结果
        </v-card-title>
        <v-card-text>
          <v-alert
            :type="validationResult.valid ? 'success' : 'error'"
            :text="validationResult.message"
          ></v-alert>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="validationDialog = false">关闭</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- 快捷键帮助对话框 -->
    <v-dialog v-model="showShortcutsDialog" max-width="600">
      <v-card>
        <v-card-title class="d-flex align-center">
          <v-icon color="info" class="mr-2">mdi-keyboard</v-icon>
          快捷键帮助
        </v-card-title>
        <v-card-text>
          <v-list density="compact">
            <v-list-subheader>编辑器快捷键</v-list-subheader>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+S</v-chip>
              </template>
              <v-list-item-title>保存配置</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+T</v-chip>
              </template>
              <v-list-item-title>验证配置</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+Shift+T</v-chip>
              </template>
              <v-list-item-title>加载模板</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">F5</v-chip>
              </template>
              <v-list-item-title>重新加载配置</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+Shift+S</v-chip>
              </template>
              <v-list-item-title>保存并验证</v-list-item-title>
            </v-list-item>

            <v-divider class="my-2"></v-divider>
            <v-list-subheader>编辑功能</v-list-subheader>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+Z</v-chip>
              </template>
              <v-list-item-title>撤销</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+Y</v-chip>
              </template>
              <v-list-item-title>重做</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+F</v-chip>
              </template>
              <v-list-item-title>查找</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+H</v-chip>
              </template>
              <v-list-item-title>查找替换</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+Shift+I</v-chip>
              </template>
              <v-list-item-title>格式化文档</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+/</v-chip>
              </template>
              <v-list-item-title>切换行注释</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Alt+Shift+F</v-chip>
              </template>
              <v-list-item-title>格式化选中内容</v-list-item-title>
            </v-list-item>

            <v-divider class="my-2"></v-divider>
            <v-list-subheader>高级编辑</v-list-subheader>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+G</v-chip>
              </template>
              <v-list-item-title>跳转到指定行</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+Shift+K</v-chip>
              </template>
              <v-list-item-title>删除当前行</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+D</v-chip>
              </template>
              <v-list-item-title>复制当前行</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Alt+↑/↓</v-chip>
              </template>
              <v-list-item-title>移动行</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">F3</v-chip>
              </template>
              <v-list-item-title>查找下一个</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Shift+F3</v-chip>
              </template>
              <v-list-item-title>查找上一个</v-list-item-title>
            </v-list-item>

            <v-divider class="my-2"></v-divider>
            <v-list-subheader>帮助与导航</v-list-subheader>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">F1</v-chip>
              </template>
              <v-list-item-title>显示快捷键帮助</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Ctrl+?</v-chip>
              </template>
              <v-list-item-title>显示快捷键帮助</v-list-item-title>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-chip size="small" variant="outlined">Esc</v-chip>
              </template>
              <v-list-item-title>关闭对话框</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="showShortcutsDialog = false">关闭</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, inject, nextTick } from 'vue'
import { configAPI } from '@/api/nginx'
import * as monaco from 'monaco-editor'

const showNotification = inject('showNotification')

// 状态
const loading = ref(false)
const saving = ref(false)
const editorContainer = ref(null)
const validationDialog = ref(false)
const validationResult = ref({ valid: false, message: '' })
const showShortcutsDialog = ref(false)

// 编辑器状态
const cursorPosition = ref({ line: 1, column: 1 })
const documentStats = ref({ 
  lines: 0, 
  characters: 0, 
  modified: false 
})

// Monaco编辑器实例
let editor = null

// 生命周期
onMounted(async () => {
  await nextTick()
  initEditor()
  await loadConfig()
  
  // 添加全局键盘事件监听器
  window.addEventListener('keydown', handleGlobalKeyDown)
})

onUnmounted(() => {
  if (editor) {
    editor.dispose()
  }
  
  // 移除全局键盘事件监听器
  window.removeEventListener('keydown', handleGlobalKeyDown)
})

// 处理全局键盘事件
const handleGlobalKeyDown = (event) => {
  // 检查是否在编辑器区域内
  const isInEditor = document.activeElement && 
    document.activeElement.closest('.monaco-editor')
  
  if (!isInEditor) return
  
  // 防止F5刷新页面
  if (event.key === 'F5') {
    event.preventDefault()
    loadConfig()
    showNotification('配置已重新加载', 'info', 2000)
    return
  }
  
  // F1 显示快捷键帮助
  if (event.key === 'F1') {
    event.preventDefault()
    showShortcutsDialog.value = true
    return
  }
  
  // Ctrl+? 也可以显示快捷键帮助
  if (event.ctrlKey && event.key === '?') {
    event.preventDefault()
    showShortcutsDialog.value = true
    return
  }
  
  // Escape 关闭对话框
  if (event.key === 'Escape') {
    if (showShortcutsDialog.value) {
      showShortcutsDialog.value = false
      event.preventDefault()
    }
    if (validationDialog.value) {
      validationDialog.value = false
      event.preventDefault()
    }
  }
}

// 初始化Monaco编辑器
const initEditor = () => {
  if (!editorContainer.value) return

  // 配置nginx语法高亮和语言特性
  monaco.languages.register({ id: 'nginx' })
  
  // 设置语言配置（支持注释等功能）
  monaco.languages.setLanguageConfiguration('nginx', {
    comments: {
      lineComment: '#'
    },
    brackets: [
      ['{', '}'],
      ['[', ']'],
      ['(', ')']
    ],
    autoClosingPairs: [
      { open: '{', close: '}' },
      { open: '[', close: ']' },
      { open: '(', close: ')' },
      { open: '"', close: '"' },
      { open: "'", close: "'" }
    ],
    surroundingPairs: [
      { open: '{', close: '}' },
      { open: '[', close: ']' },
      { open: '(', close: ')' },
      { open: '"', close: '"' },
      { open: "'", close: "'" }
    ],
    folding: {
      markers: {
        start: new RegExp('^\\s*#\\s*region\\b'),
        end: new RegExp('^\\s*#\\s*endregion\\b')
      }
    }
  })
  
  // 设置语法高亮
  monaco.languages.setMonarchTokensProvider('nginx', {
    tokenizer: {
      root: [
        [/#.*$/, 'comment'],
        [/\b(server|location|upstream|proxy_pass|listen|server_name|root|index|error_page|access_log|error_log|worker_processes|worker_connections|sendfile|keepalive_timeout|gzip|ssl_certificate|ssl_certificate_key)\b/, 'keyword'],
        [/\b(http|events|stream|mail)\b/, 'type'],
        [/".*?"/, 'string'],
        [/'.*?'/, 'string'],
        [/\d+[kmg]?/, 'number'],
        [/[{}();]/, 'delimiter'],
        [/[a-zA-Z_]\w*/, 'identifier']
      ]
    }
  })

  // 创建编辑器
  editor = monaco.editor.create(editorContainer.value, {
    value: '',
    language: 'nginx',
    theme: 'vs-dark',
    automaticLayout: true,
    minimap: { enabled: true },
    scrollBeyondLastLine: false,
    fontSize: 14,
    wordWrap: 'on',
    lineNumbers: 'on',
    folding: true,
    selectOnLineNumbers: true,
    roundedSelection: false,
    readOnly: false,
    cursorStyle: 'line',
    automaticLayout: true,
  })

  // 添加自定义快捷键
  setupKeyboardShortcuts()
  
  // 设置编辑器事件监听器
  setupEditorListeners()
}

// 设置编辑器事件监听器
const setupEditorListeners = () => {
  if (!editor) return

  // 监听光标位置变化
  editor.onDidChangeCursorPosition((e) => {
    cursorPosition.value = {
      line: e.position.lineNumber,
      column: e.position.column
    }
  })

  // 监听文档内容变化
  editor.onDidChangeModelContent(() => {
    updateDocumentStats()
    documentStats.value.modified = true
  })

  // 初始化文档统计
  updateDocumentStats()
}

// 更新文档统计信息
const updateDocumentStats = () => {
  if (!editor) return
  
  const model = editor.getModel()
  if (model) {
    const content = model.getValue()
    documentStats.value.lines = model.getLineCount()
    documentStats.value.characters = content.length
  }
}

// 设置键盘快捷键
const setupKeyboardShortcuts = () => {
  if (!editor) return

  // Ctrl+S: 保存配置
  editor.addAction({
    id: 'save-config',
    label: 'Save Config',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 1.5,
    run: function(ed) {
      showNotification('正在保存配置...', 'info', 1500)
      saveConfig()
      return null
    }
  })

  // Ctrl+T: 验证配置
  editor.addAction({
    id: 'validate-config',
    label: 'Validate Config',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyT],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 1.6,
    run: function(ed) {
      showNotification('正在验证配置...', 'info', 1500)
      validateConfig()
      return null
    }
  })

  // Ctrl+Shift+T: 加载模板
  editor.addAction({
    id: 'load-template',
    label: 'Load Template',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyMod.Shift | monaco.KeyCode.KeyT],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 1.7,
    run: function(ed) {
      if (confirm('加载模板将覆盖当前内容，确定继续吗？')) {
        showNotification('正在加载模板...', 'info', 1500)
        loadTemplate()
      }
      return null
    }
  })

  // F5: 重新加载配置
  editor.addAction({
    id: 'reload-config',
    label: 'Reload Config',
    keybindings: [monaco.KeyCode.F5],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 1.8,
    run: function(ed) {
      showNotification('正在重新加载配置...', 'info', 1500)
      loadConfig()
      return null
    }
  })

  // Ctrl+Shift+S: 保存并验证
  editor.addAction({
    id: 'save-and-validate',
    label: 'Save and Validate',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyMod.Shift | monaco.KeyCode.KeyS],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 1.9,
    run: async function(ed) {
      showNotification('正在保存并验证配置...', 'info', 2000)
      await saveConfig()
      setTimeout(() => validateConfig(), 500)
      return null
    }
  })

  // Ctrl+Shift+I: 格式化文档
  editor.addAction({
    id: 'format-document',
    label: 'Format Document',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyMod.Shift | monaco.KeyCode.KeyI],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 2.0,
    run: function(ed) {
      ed.trigger('keyboard', 'editor.action.formatDocument', {})
      return null
    }
  })

  // Alt+Shift+F: 格式化选中内容
  editor.addAction({
    id: 'format-selection',
    label: 'Format Selection',
    keybindings: [monaco.KeyMod.Alt | monaco.KeyMod.Shift | monaco.KeyCode.KeyF],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 2.1,
    run: function(ed) {
      ed.trigger('keyboard', 'editor.action.formatSelection', {})
      return null
    }
  })

  // Ctrl+/: 切换行注释
  editor.addAction({
    id: 'toggle-line-comment',
    label: 'Toggle Line Comment',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.Slash],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 2.2,
    run: function(ed) {
      // 使用编辑器内置的注释切换功能
      ed.trigger('keyboard', 'editor.action.commentLine', {})
      return null
    }
  })

  // Ctrl+G: 跳转到指定行
  editor.addAction({
    id: 'go-to-line',
    label: 'Go to Line',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyG],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 2.3,
    run: function(ed) {
      ed.trigger('keyboard', 'editor.action.gotoLine', {})
      return null
    }
  })

  // Ctrl+Shift+K: 删除当前行
  editor.addAction({
    id: 'delete-line',
    label: 'Delete Line',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyMod.Shift | monaco.KeyCode.KeyK],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 2.4,
    run: function(ed) {
      ed.trigger('keyboard', 'editor.action.deleteLines', {})
      return null
    }
  })

  // Alt+Up/Down: 移动行
  editor.addAction({
    id: 'move-line-up',
    label: 'Move Line Up',
    keybindings: [monaco.KeyMod.Alt | monaco.KeyCode.UpArrow],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 2.5,
    run: function(ed) {
      ed.trigger('keyboard', 'editor.action.moveLinesUpAction', {})
      return null
    }
  })

  editor.addAction({
    id: 'move-line-down',
    label: 'Move Line Down',
    keybindings: [monaco.KeyMod.Alt | monaco.KeyCode.DownArrow],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 2.6,
    run: function(ed) {
      ed.trigger('keyboard', 'editor.action.moveLinesDownAction', {})
      return null
    }
  })

  // Ctrl+D: 复制当前行（修改为标准快捷键）
  editor.addAction({
    id: 'copy-line-down',
    label: 'Copy Line Down',
    keybindings: [monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyD],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 2.7,
    run: function(ed) {
      ed.trigger('keyboard', 'editor.action.copyLinesDownAction', {})
      return null
    }
  })

  // F3/Shift+F3: 查找下一个/上一个
  editor.addAction({
    id: 'find-next',
    label: 'Find Next',
    keybindings: [monaco.KeyCode.F3],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 2.8,
    run: function(ed) {
      ed.trigger('keyboard', 'editor.action.nextMatchFindAction', {})
      return null
    }
  })

  editor.addAction({
    id: 'find-previous',
    label: 'Find Previous',
    keybindings: [monaco.KeyMod.Shift | monaco.KeyCode.F3],
    precondition: null,
    keybindingContext: null,
    contextMenuGroupId: 'navigation',
    contextMenuOrder: 2.9,
    run: function(ed) {
      ed.trigger('keyboard', 'editor.action.previousMatchFindAction', {})
      return null
    }
  })
}

// 加载配置文件
const loadConfig = async () => {
  try {
    loading.value = true
    const response = await configAPI.getConfig()
    if (response.success && editor) {
      editor.setValue(response.data)
      // 重置修改状态
      documentStats.value.modified = false
      updateDocumentStats()
    }
  } catch (error) {
    showNotification('加载配置失败: ' + error.message, 'error')
  } finally {
    loading.value = false
  }
}

// 保存配置文件
const saveConfig = async () => {
  if (!editor) return

  try {
    saving.value = true
    const content = editor.getValue()
    
    // 检查内容是否为空
    if (!content.trim()) {
      showNotification('配置内容不能为空', 'warning')
      return
    }
    
    const response = await configAPI.saveConfig(content)
    
    if (response.success) {
      showNotification('配置保存成功', 'success', 3000)
      // 标记文档为已保存状态
      markDocumentAsSaved()
    } else {
      showNotification(response.message, 'error')
    }
  } catch (error) {
    showNotification('保存配置失败: ' + error.message, 'error')
  } finally {
    saving.value = false
  }
}

// 标记文档为已保存状态
const markDocumentAsSaved = () => {
  documentStats.value.modified = false
  if (editor) {
    // 清除修改标记
    const model = editor.getModel()
    if (model) {
      model.pushStackElement()
    }
  }
}

// 验证配置
const validateConfig = async () => {
  if (!editor) return

  try {
    loading.value = true
    const content = editor.getValue()
    const response = await configAPI.validateConfig(content)
    
    validationResult.value = {
      valid: response.valid,
      message: response.message
    }
    validationDialog.value = true
  } catch (error) {
    validationResult.value = {
      valid: false,
      message: '验证失败: ' + error.message
    }
    validationDialog.value = true
  } finally {
    loading.value = false
  }
}

// 加载配置模板
const loadTemplate = async () => {
  try {
    loading.value = true
    const response = await configAPI.getTemplate()
    
    if (response.success && editor) {
      editor.setValue(response.data)
      // 标记为已修改（因为加载了新内容）
      documentStats.value.modified = true
      updateDocumentStats()
      showNotification('模板加载成功', 'success')
    }
  } catch (error) {
    showNotification('加载模板失败: ' + error.message, 'error')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.monaco-editor {
  border: 1px solid #ddd;
}
</style> 