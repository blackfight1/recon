<template>
  <div class="quick-scan">
    <h2>快速扫描</h2>
    <p class="description">输入域名立即扫描，无需添加到监控列表</p>

    <el-card class="scan-card">
      <el-form :model="form" @submit.prevent="startScan">
        <el-form-item label="目标域名">
          <el-input
            v-model="form.domain"
            placeholder="例如: example.com"
            :disabled="scanning"
          >
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
    <el-card v-if="taskId" class="progress-card">
      <template #header>
        <div class="card-header">
          <span>扫描进度</span>
          <el-tag :type="getStatusType(task.status)">
            {{ getStatusText(task.status) }}
          </el-tag>
        </div>
      </template>

      <el-progress
        :percentage="task.progress"
        :status="task.status === 'completed' ? 'success' : task.status === 'failed' ? 'exception' : ''"
      />
      <div class="current-step">
        当前步骤: {{ task.current_step || '准备中...' }}
      </div>

      <!-- 实时日志 -->
      <div class="logs-container">
        <div class="logs-header">
          <span>扫描日志</span>
          <el-button size="small" @click="loadLogs">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
        <div class="logs-content">
          <div
            v-for="log in logs"
            :key="log.id"
            :class="['log-item', `log-${log.level}`]"
          >
            <span class="log-time">{{ formatTime(log.created_at) }}</span>
            <span class="log-step">[{{ log.step }}]</span>
            <span class="log-message">{{ log.message }}</span>
          </div>
          <div v-if="logs.length === 0" class="no-logs">
            暂无日志...
          </div>
        </div>
      </div>

      <!-- 扫描结果 -->
      <div v-if="task.status === 'completed'" class="result-section">
        <el-divider />
        <h3>扫描结果</h3>
        <div class="result-content">
          {{ task.result }}
        </div>
        <el-button type="primary" @click="viewSubdomains">
          查看子域名列表
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
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
  result: ''
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
  // 跳转到子域名列表
  window.location.href = '/subdomains'
}

// 组件卸载时停止轮询
onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.quick-scan h2 {
  margin-bottom: 10px;
  color: #303133;
}

.description {
  color: #909399;
  margin-bottom: 20px;
}

.scan-card {
  margin-bottom: 20px;
}

.progress-card {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.current-step {
  margin-top: 10px;
  color: #606266;
  font-size: 14px;
}

.logs-container {
  margin-top: 20px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
}

.logs-content {
  max-height: 400px;
  overflow-y: auto;
  padding: 10px;
  background: #fff;
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.log-item {
  padding: 5px 0;
  border-bottom: 1px solid #f0f0f0;
}

.log-item:last-child {
  border-bottom: none;
}

.log-time {
  color: #909399;
  margin-right: 10px;
}

.log-step {
  color: #409eff;
  margin-right: 10px;
  font-weight: bold;
}

.log-message {
  color: #303133;
}

.log-info .log-message {
  color: #606266;
}

.log-success .log-message {
  color: #67c23a;
}

.log-warning .log-message {
  color: #e6a23c;
}

.log-error .log-message {
  color: #f56c6c;
}

.no-logs {
  text-align: center;
  color: #909399;
  padding: 20px;
}

.result-section {
  margin-top: 20px;
}

.result-section h3 {
  margin-bottom: 10px;
}

.result-content {
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;
  margin-bottom: 15px;
}
</style>
