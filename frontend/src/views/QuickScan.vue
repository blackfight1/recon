<template>
  <div class="quick-scan">
    <div class="page-header">
      <h2><el-icon><Search /></el-icon> 快速扫描</h2>
      <p class="description">输入域名立即扫描，实时查看进度和日志</p>
    </div>

    <el-card class="scan-card" shadow="hover">
      <el-form :model="form" @submit.prevent="startScan">
        <el-form-item>
          <el-input
            v-model="form.domain"
            placeholder="输入域名，例如: example.com"
            :disabled="scanning"
            size="large"
            clearable
          >
            <template #prefix>
              <el-icon><Link /></el-icon>
            </template>
            <template #append>
              <el-button
                type="primary"
                @click="startScan"
                :loading="scanning"
                :disabled="!form.domain"
              >
                <el-icon v-if="!scanning"><Search /></el-icon>
                {{ scanning ? '扫描中...' : '开始扫描' }}
              </el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 扫描进度 -->
    <el-card v-if="taskId" class="progress-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span><el-icon><Loading /></el-icon> 扫描进度</span>
          <el-tag :type="getStatusType(task.status)" size="large">
            {{ getStatusText(task.status) }}
          </el-tag>
        </div>
      </template>

      <div class="progress-section">
        <el-progress
          :percentage="task.progress"
          :status="task.status === 'completed' ? 'success' : task.status === 'failed' ? 'exception' : ''"
          :stroke-width="12"
          striped
          striped-flow
        />
        <div class="current-step">
          <el-icon><Operation /></el-icon>
          当前步骤: <span class="step-text">{{ task.current_step || '准备中...' }}</span>
        </div>
      </div>

      <!-- 实时日志 -->
      <div class="logs-container">
        <div class="logs-header">
          <span><el-icon><Document /></el-icon> 扫描日志</span>
          <el-button size="small" @click="loadLogs" :icon="Refresh">
            刷新
          </el-button>
        </div>
        <div class="logs-content">
          <div
            v-for="log in logs"
            :key="log.id"
            :class="['log-item', `log-${log.level}`]"
          >
            <span class="log-icon">
              <el-icon v-if="log.level === 'success'"><SuccessFilled /></el-icon>
              <el-icon v-else-if="log.level === 'error'"><CircleCloseFilled /></el-icon>
              <el-icon v-else-if="log.level === 'warning'"><WarningFilled /></el-icon>
              <el-icon v-else><InfoFilled /></el-icon>
            </span>
            <span class="log-time">{{ formatTime(log.created_at) }}</span>
            <span class="log-step">[{{ log.step }}]</span>
            <span class="log-message">{{ log.message }}</span>
          </div>
          <div v-if="logs.length === 0" class="no-logs">
            <el-icon><Loading /></el-icon>
            <span>等待日志...</span>
          </div>
        </div>
      </div>

      <!-- 扫描结果 -->
      <div v-if="task.status === 'completed'" class="result-section">
        <el-divider />
        <div class="result-header">
          <el-icon><CircleCheckFilled /></el-icon>
          <h3>扫描完成</h3>
        </div>
        <div class="result-content">
          <el-icon><DataAnalysis /></el-icon>
          {{ task.result }}
        </div>
        <el-button type="primary" @click="viewSubdomains" :icon="View">
          查看子域名列表
        </el-button>
      </div>

      <!-- 扫描失败 -->
      <div v-if="task.status === 'failed'" class="error-section">
        <el-divider />
        <div class="error-header">
          <el-icon><CircleCloseFilled /></el-icon>
          <h3>扫描失败</h3>
        </div>
        <div class="error-content">
          {{ task.error_msg || '未知错误' }}
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Search,
  Link,
  Loading,
  Operation,
  Document,
  Refresh,
  SuccessFilled,
  CircleCloseFilled,
  WarningFilled,
  InfoFilled,
  CircleCheckFilled,
  DataAnalysis,
  View
} from '@element-plus/icons-vue'
import api from '../api'

const form = ref({
  domain: ''
})

const scanning = ref(false)
const taskId = ref(null)
const task = ref({
  progress: 0,
  status: 'pending',
  current_step: '',
  result: '',
  error_msg: ''
})
const logs = ref([])
let progressInterval = null
let logsInterval = null

const startScan = async () => {
  if (!form.value.domain) {
    ElMessage.warning('请输入域名')
    return
  }

  scanning.value = true
  try {
    const res = await api.quickScan(form.value.domain)
    taskId.value = res.task_id
    ElMessage.success('扫描已开始')
    
    // 开始轮询进度和日志
    startPolling()
  } catch (error) {
    ElMessage.error('启动扫描失败')
    scanning.value = false
  }
}

const startPolling = () => {
  // 每2秒更新进度
  progressInterval = setInterval(loadProgress, 2000)
  // 每3秒更新日志
  logsInterval = setInterval(loadLogs, 3000)
  
  // 立即加载一次
  loadProgress()
  loadLogs()
}

const stopPolling = () => {
  if (progressInterval) {
    clearInterval(progressInterval)
    progressInterval = null
  }
  if (logsInterval) {
    clearInterval(logsInterval)
    logsInterval = null
  }
}

const loadProgress = async () => {
  if (!taskId.value) return
  
  try {
    const res = await api.getTaskProgress(taskId.value)
    task.value = res.data
    
    // 如果任务完成或失败，停止轮询
    if (task.value.status === 'completed' || task.value.status === 'failed') {
      stopPolling()
      scanning.value = false
      
      // 最后加载一次日志
      loadLogs()
    }
  } catch (error) {
    console.error('Failed to load progress:', error)
  }
}

const loadLogs = async () => {
  if (!taskId.value) return
  
  try {
    const res = await api.getTaskLogs(taskId.value)
    logs.value = res.data
  } catch (error) {
    console.error('Failed to load logs:', error)
  }
}

const getStatusType = (status) => {
  const types = {
    pending: 'info',
    running: 'warning',
    completed: 'success',
    failed: 'danger'
  }
  return types[status] || ''
}

const getStatusText = (status) => {
  const texts = {
    pending: '等待中',
    running: '运行中',
    completed: '已完成',
    failed: '失败'
  }
  return texts[status] || status
}

const formatTime = (date) => {
  return new Date(date).toLocaleTimeString('zh-CN')
}

const viewSubdomains = () => {
  window.location.href = '/subdomains'
}

// 组件卸载时停止轮询
onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.quick-scan {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h2 {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 28px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.description {
  color: var(--text-secondary);
  font-size: 14px;
}

.scan-card {
  margin-bottom: 24px;
}

.progress-card {
  margin-top: 24px;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
}

.card-header span {
  display: flex;
  align-items: center;
  gap: 8px;
}

.progress-section {
  margin-bottom: 24px;
}

.current-step {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 16px;
  padding: 12px;
  background: var(--bg-tertiary);
  border-radius: 8px;
  color: var(--text-secondary);
  font-size: 14px;
}

.step-text {
  color: var(--accent-primary);
  font-weight: 500;
}

.logs-container {
  margin-top: 24px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: var(--bg-tertiary);
  border-bottom: 1px solid var(--border-color);
  font-weight: 500;
}

.logs-header span {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logs-content {
  max-height: 400px;
  overflow-y: auto;
  padding: 16px;
  background: var(--bg-primary);
  font-family: 'Courier New', 'Consolas', monospace;
  font-size: 13px;
}

.logs-content::-webkit-scrollbar {
  width: 8px;
}

.logs-content::-webkit-scrollbar-track {
  background: var(--bg-tertiary);
}

.logs-content::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 4px;
}

.logs-content::-webkit-scrollbar-thumb:hover {
  background: var(--accent-primary);
}

.log-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid var(--border-color);
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.log-item:last-child {
  border-bottom: none;
}

.log-icon {
  font-size: 16px;
}

.log-time {
  color: var(--text-secondary);
  font-size: 12px;
  min-width: 80px;
}

.log-step {
  color: var(--accent-primary);
  font-weight: 600;
  min-width: 120px;
}

.log-message {
  color: var(--text-primary);
  flex: 1;
}

.log-info .log-icon {
  color: var(--accent-primary);
}

.log-success .log-icon {
  color: var(--success);
}

.log-warning .log-icon {
  color: var(--warning);
}

.log-error .log-icon {
  color: var(--danger);
}

.no-logs {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px;
  color: var(--text-secondary);
}

.no-logs .el-icon {
  font-size: 32px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.result-section,
.error-section {
  margin-top: 24px;
}

.result-header,
.error-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  font-size: 18px;
  font-weight: 600;
}

.result-header {
  color: var(--success);
}

.error-header {
  color: var(--danger);
}

.result-content,
.error-content {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: var(--bg-tertiary);
  border-radius: 8px;
  margin-bottom: 16px;
  font-size: 14px;
}

.result-content {
  color: var(--success);
  border-left: 4px solid var(--success);
}

.error-content {
  color: var(--danger);
  border-left: 4px solid var(--danger);
}
</style>
