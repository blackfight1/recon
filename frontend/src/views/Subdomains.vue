<template>
  <div class="subdomains">
    <div class="page-header">
      <h2>子域名列表</h2>
      <el-select v-model="statusFilter" placeholder="筛选状态" @change="loadSubdomains" clearable>
        <el-option label="全部" value="" />
        <el-option label="新增" value="new" />
        <el-option label="存活" value="alive" />
        <el-option label="失效" value="dead" />
      </el-select>
    </div>

    <el-table :data="subdomains" style="width: 100%" v-loading="loading">
      <el-table-column prop="subdomain" label="子域名" min-width="250" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status_code" label="状态码" width="100" />
      <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
      <el-table-column prop="server" label="服务器" width="150" />
      <el-table-column prop="ip" label="IP" width="150" />
      <el-table-column prop="last_seen" label="最后发现" width="180">
        <template #default="{ row }">
          {{ formatDate(row.last_seen) }}
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[20, 50, 100]"
      layout="total, sizes, prev, pager, next"
      @current-change="loadSubdomains"
      @size-change="loadSubdomains"
      style="margin-top: 20px; justify-content: center;"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../api'

const subdomains = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(50)
const total = ref(0)
const statusFilter = ref('')

const loadSubdomains = async () => {
  loading.value = true
  try {
    const res = await api.getSubdomains({
      page: currentPage.value,
      page_size: pageSize.value,
      status: statusFilter.value
    })
    subdomains.value = res.data
    total.value = res.total
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const getStatusType = (status) => {
  const types = {
    new: 'warning',
    alive: 'success',
    dead: 'info'
  }
  return types[status] || ''
}

const getStatusText = (status) => {
  const texts = {
    new: '新增',
    alive: '存活',
    dead: '失效'
  }
  return texts[status] || status
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  loadSubdomains()
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  color: var(--text-primary);
}
</style>
