<template>
  <el-config-provider :locale="zhCn">
    <el-container class="app-container">
      <el-header class="app-header">
        <div class="header-content">
          <div class="logo">
            <el-icon class="logo-icon"><Radar /></el-icon>
            <h1>自动化侦查平台</h1>
          </div>
          <el-menu
            mode="horizontal"
            :default-active="activeMenu"
            router
            class="header-menu"
          >
            <el-menu-item index="/">
              <el-icon><Monitor /></el-icon>
              <span>仪表盘</span>
            </el-menu-item>
            <el-menu-item index="/quick-scan">
              <el-icon><Search /></el-icon>
              <span>快速扫描</span>
            </el-menu-item>
            <el-menu-item index="/targets">
              <el-icon><Aim /></el-icon>
              <span>目标管理</span>
            </el-menu-item>
            <el-menu-item index="/subdomains">
              <el-icon><Connection /></el-icon>
              <span>子域名</span>
            </el-menu-item>
            <el-menu-item index="/changes">
              <el-icon><Bell /></el-icon>
              <span>变更中心</span>
            </el-menu-item>
            <el-menu-item index="/tasks">
              <el-icon><List /></el-icon>
              <span>任务列表</span>
            </el-menu-item>
          </el-menu>
        </div>
      </el-header>
      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>
  </el-config-provider>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Monitor, Search, Aim, Connection, Bell, List, Radar } from '@element-plus/icons-vue'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'

const route = useRoute()
const activeMenu = computed(() => route.path)
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:root {
  --bg-primary: #0a0e27;
  --bg-secondary: #151932;
  --bg-tertiary: #1e2139;
  --border-color: #2a2f4a;
  --text-primary: #e4e7eb;
  --text-secondary: #9ca3af;
  --accent-primary: #3b82f6;
  --accent-hover: #2563eb;
  --success: #10b981;
  --warning: #f59e0b;
  --danger: #ef4444;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background: var(--bg-primary);
  color: var(--text-primary);
}

.app-container {
  min-height: 100vh;
  background: var(--bg-primary);
}

.app-header {
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  padding: 0;
  height: 64px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.header-content {
  max-width: 1600px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  height: 100%;
  padding: 0 24px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-right: 48px;
}

.logo-icon {
  font-size: 28px;
  color: var(--accent-primary);
}

.logo h1 {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  background: linear-gradient(135deg, var(--accent-primary), #8b5cf6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.header-menu {
  flex: 1;
  background: transparent;
  border: none;
}

.header-menu .el-menu-item {
  color: var(--text-secondary);
  border-bottom: 2px solid transparent;
  transition: all 0.3s;
}

.header-menu .el-menu-item:hover {
  color: var(--accent-primary);
  background: rgba(59, 130, 246, 0.1);
}

.header-menu .el-menu-item.is-active {
  color: var(--accent-primary);
  border-bottom-color: var(--accent-primary);
  background: rgba(59, 130, 246, 0.1);
}

.app-main {
  max-width: 1600px;
  margin: 0 auto;
  padding: 24px;
  background: var(--bg-primary);
}

/* Element Plus 暗色主题覆盖 */
.el-card {
  background: var(--bg-secondary);
  border-color: var(--border-color);
  color: var(--text-primary);
}

.el-card__header {
  border-bottom-color: var(--border-color);
  color: var(--text-primary);
}

.el-table {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.el-table th.el-table__cell {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  border-bottom-color: var(--border-color);
}

.el-table td.el-table__cell {
  border-bottom-color: var(--border-color);
}

.el-table tr {
  background: var(--bg-secondary);
}

.el-table tr:hover > td {
  background: var(--bg-tertiary) !important;
}

.el-input__wrapper {
  background: var(--bg-tertiary);
  border-color: var(--border-color);
  box-shadow: none;
}

.el-input__wrapper:hover {
  border-color: var(--accent-primary);
}

.el-input__wrapper.is-focus {
  border-color: var(--accent-primary);
  box-shadow: 0 0 0 1px var(--accent-primary) inset;
}

.el-input__inner {
  color: var(--text-primary);
}

.el-button--primary {
  background: var(--accent-primary);
  border-color: var(--accent-primary);
}

.el-button--primary:hover {
  background: var(--accent-hover);
  border-color: var(--accent-hover);
}

.el-tag {
  border-color: var(--border-color);
}

.el-tag--success {
  background: rgba(16, 185, 129, 0.1);
  border-color: var(--success);
  color: var(--success);
}

.el-tag--warning {
  background: rgba(245, 158, 11, 0.1);
  border-color: var(--warning);
  color: var(--warning);
}

.el-tag--danger {
  background: rgba(239, 68, 68, 0.1);
  border-color: var(--danger);
  color: var(--danger);
}

.el-tag--info {
  background: rgba(59, 130, 246, 0.1);
  border-color: var(--accent-primary);
  color: var(--accent-primary);
}

.el-progress__text {
  color: var(--text-primary);
}

.el-dialog {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
}

.el-dialog__header {
  border-bottom: 1px solid var(--border-color);
}

.el-dialog__title {
  color: var(--text-primary);
}

.el-form-item__label {
  color: var(--text-secondary);
}

.el-textarea__inner {
  background: var(--bg-tertiary);
  border-color: var(--border-color);
  color: var(--text-primary);
}

.el-select__wrapper {
  background: var(--bg-tertiary);
  border-color: var(--border-color);
}

.el-select-dropdown {
  background: var(--bg-secondary);
  border-color: var(--border-color);
}

.el-select-dropdown__item {
  color: var(--text-primary);
}

.el-select-dropdown__item:hover {
  background: var(--bg-tertiary);
}

.el-pagination {
  color: var(--text-primary);
}

.el-pagination button {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.el-pagination button:hover {
  color: var(--accent-primary);
}

.el-pager li {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.el-pager li:hover {
  color: var(--accent-primary);
}

.el-pager li.is-active {
  background: var(--accent-primary);
  color: white;
}
</style>
