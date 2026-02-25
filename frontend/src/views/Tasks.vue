<template>
  <div class="tasks">
    <h2>任务列表</h2>

    <el-table :data="tasks" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="target_id" label="目标ID" width="100" />
      <el-table-column prop="task_type" label="任务类型" width="120">
        <template #default="{ row }">
          <el-tag>{{ getTaskTypeText(row.task_type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="result" label="结果" min-width="200" show-overflow-tooltip />
      <el-table-column prop="error_msg" label="错误信息" min-width="200" show-overflow-tooltip />
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="completed_at" label="完成时间" width="180">
        <template #default="{ row }">
          {{ row.completed_at ? formatDate(row.completed_at) : '-' }}
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50]"
      layout="total, sizes, prev, pager, next"
      @current-change="loadTasks"
      @size-change="loadTasks"
      style="margin-top: 20px; justify-content: center;"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../api'

const tasks = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const loadTasks = async () => {
  loading.value = true
  try {
    const res = await api.getTasks({
      page: currentPage.value,
      page_size: pageSize.value
    })
    tasks.value = res.data
    total.value = res.total
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const getTaskTypeText = (type) => {
  const texts = {
    subdomain: '子域名扫描',
    port: '端口扫描',
    screenshot: '截图',
    vuln: '漏洞扫描'
  }
  return texts[type] || type
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

const formatDate = (date) => {
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  loadTasks()
  // 每10秒刷新一次
  setInterval(loadTasks, 10000)
})
</script>

<style scoped>
.tasks h2 {
  margin-bottom: 20px;
  color: var(--text-primary);
}
</style>
