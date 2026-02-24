<template>
  <el-container class="app-container">
    <el-header class="app-header">
      <div class="header-content">
        <h1>🔍 自动化侦查平台</h1>
        <el-menu mode="horizontal" :default-active="activeMenu" router>
          <el-menu-item index="/">
            <el-icon><Monitor /></el-icon>
            <span>仪表盘</span>
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
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Monitor, Aim, Connection, Bell, List } from '@element-plus/icons-vue'

const route = useRoute()
const activeMenu = computed(() => route.path)
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background: #f5f7fa;
}

.app-container {
  height: 100vh;
}

.app-header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 0;
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  height: 100%;
}

.header-content h1 {
  font-size: 20px;
  margin-right: 40px;
  color: #409eff;
}

.app-main {
  max-width: 1400px;
  margin: 20px auto;
  padding: 20px;
}

.el-menu--horizontal {
  border-bottom: none;
}
</style>
