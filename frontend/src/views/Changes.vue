<template>
  <div class="changes">
    <h2>变更中心</h2>

    <el-timeline>
      <el-timeline-item
        v-for="change in changes"
        :key="change.id"
        :timestamp="formatDate(change.created_at)"
        placement="top"
        :color="getChangeColor(change.change_type)"
      >
        <el-card>
          <div class="change-item">
            <el-icon :size="24" :color="getChangeColor(change.change_type)">
              <component :is="getChangeIcon(change.change_type)" />
            </el-icon>
            <div class="change-content">
              <div class="change-type">{{ getChangeTypeText(change.change_type) }}</div>
              <div class="change-detail">{{ change.content }}</div>
            </div>
          </div>
        </el-card>
      </el-timeline-item>
    </el-timeline>

    <el-pagination
      v-if="total > pageSize"
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[20, 50, 100]"
      layout="total, sizes, prev, pager, next"
      @current-change="loadChanges"
      @size-change="loadChanges"
      style="margin-top: 20px; justify-content: center;"
    />

    <el-empty v-if="!loading && changes.length === 0" description="暂无变更记录" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Check, Close } from '@element-plus/icons-vue'
import api from '../api'

const changes = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(50)
const total = ref(0)

const loadChanges = async () => {
  loading.value = true
  try {
    const res = await api.getChanges({
      page: currentPage.value,
      page_size: pageSize.value
    })
    changes.value = res.data
    total.value = res.total
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const getChangeTypeText = (type) => {
  const texts = {
    subdomain_new: '新增子域名',
    subdomain_alive: '子域名恢复',
    subdomain_dead: '子域名失效'
  }
  return texts[type] || type
}

const getChangeColor = (type) => {
  const colors = {
    subdomain_new: '#e6a23c',
    subdomain_alive: '#67c23a',
    subdomain_dead: '#909399'
  }
  return colors[type] || '#409eff'
}

const getChangeIcon = (type) => {
  const icons = {
    subdomain_new: Plus,
    subdomain_alive: Check,
    subdomain_dead: Close
  }
  return icons[type] || Plus
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  loadChanges()
})
</script>

<style scoped>
.changes h2 {
  margin-bottom: 30px;
  color: var(--text-primary);
}

.change-item {
  display: flex;
  align-items: center;
  gap: 15px;
}

.change-content {
  flex: 1;
}

.change-type {
  font-weight: bold;
  color: var(--text-primary);
  margin-bottom: 5px;
}

.change-detail {
  color: var(--text-secondary);
  font-size: 14px;
}
</style>
