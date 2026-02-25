<template>
  <div class="targets">
    <div class="page-header">
      <h2>目标管理</h2>
      <el-button type="primary" @click="showAddDialog = true">
        <el-icon><Plus /></el-icon>
        添加目标
      </el-button>
    </div>

    <el-table :data="targets" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="domain" label="域名" min-width="200" />
      <el-table-column prop="description" label="描述" min-width="200" />
      <el-table-column prop="enabled" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.enabled ? 'success' : 'info'">
            {{ row.enabled ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="triggerScan(row)">
            <el-icon><Refresh /></el-icon>
            扫描
          </el-button>
          <el-button size="small" type="primary" @click="editTarget(row)">
            编辑
          </el-button>
          <el-button size="small" type="danger" @click="deleteTarget(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="showAddDialog"
      :title="editingTarget ? '编辑目标' : '添加目标'"
      width="500px"
    >
      <el-form :model="form" label-width="80px">
        <el-form-item label="域名">
          <el-input v-model="form.domain" placeholder="example.com" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="状态" v-if="editingTarget">
          <el-switch v-model="form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveTarget">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import api from '../api'

const targets = ref([])
const loading = ref(false)
const showAddDialog = ref(false)
const editingTarget = ref(null)
const form = ref({
  domain: '',
  description: '',
  enabled: true
})

const loadTargets = async () => {
  loading.value = true
  try {
    const res = await api.getTargets()
    targets.value = res.data
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const editTarget = (target) => {
  editingTarget.value = target
  form.value = {
    domain: target.domain,
    description: target.description,
    enabled: target.enabled
  }
  showAddDialog.value = true
}

const saveTarget = async () => {
  try {
    if (editingTarget.value) {
      await api.updateTarget(editingTarget.value.id, form.value)
      ElMessage.success('更新成功')
    } else {
      await api.createTarget(form.value)
      ElMessage.success('添加成功')
    }
    showAddDialog.value = false
    editingTarget.value = null
    form.value = { domain: '', description: '', enabled: true }
    loadTargets()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const deleteTarget = async (target) => {
  try {
    await ElMessageBox.confirm('确定要删除这个目标吗？', '警告', {
      type: 'warning'
    })
    await api.deleteTarget(target.id)
    ElMessage.success('删除成功')
    loadTargets()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const triggerScan = async (target) => {
  try {
    await api.triggerScan(target.id)
    ElMessage.success('扫描任务已触发')
  } catch (error) {
    ElMessage.error('触发失败')
  }
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  loadTargets()
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
