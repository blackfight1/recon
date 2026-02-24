<template>
  <div class="dashboard">
    <h2>仪表盘</h2>
    
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" color="#409eff"><Aim /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.total_targets }}</div>
              <div class="stat-label">监控目标</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" color="#67c23a"><Connection /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.alive_subdomains }}</div>
              <div class="stat-label">存活子域名</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" color="#e6a23c"><Star /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.new_subdomains }}</div>
              <div class="stat-label">新增子域名</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" color="#f56c6c"><Bell /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.recent_changes }}</div>
              <div class="stat-label">7天变更</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="info-card">
      <template #header>
        <span>快速开始</span>
      </template>
      <el-steps :active="1" align-center>
        <el-step title="添加目标" description="在目标管理中添加要监控的主域名" />
        <el-step title="等待扫描" description="系统每6小时自动扫描一次" />
        <el-step title="查看变更" description="在变更中心查看资产变化" />
      </el-steps>
      
      <div class="quick-actions">
        <el-button type="primary" @click="$router.push('/targets')">
          <el-icon><Plus /></el-icon>
          添加监控目标
        </el-button>
        <el-button @click="$router.push('/changes')">
          <el-icon><Bell /></el-icon>
          查看变更
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Aim, Connection, Star, Bell, Plus } from '@element-plus/icons-vue'
import api from '../api'

const stats = ref({
  total_targets: 0,
  enabled_targets: 0,
  total_subdomains: 0,
  alive_subdomains: 0,
  new_subdomains: 0,
  recent_changes: 0
})

const loadStats = async () => {
  try {
    const res = await api.getStats()
    stats.value = res.data
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
}

onMounted(() => {
  loadStats()
  // 每30秒刷新一次
  setInterval(loadStats, 30000)
})
</script>

<style scoped>
.dashboard h2 {
  margin-bottom: 20px;
  color: #303133;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-5px);
}

.stat-content {
  display: flex;
  align-items: center;
  padding: 10px 0;
}

.stat-icon {
  font-size: 48px;
  margin-right: 20px;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}

.info-card {
  margin-top: 20px;
}

.quick-actions {
  margin-top: 30px;
  text-align: center;
}

.quick-actions .el-button {
  margin: 0 10px;
}
</style>
