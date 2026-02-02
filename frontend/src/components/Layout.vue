<template>
  <div class="layout">
    <!-- 侧边栏 -->
    <el-aside width="200px" class="layout-aside">
      <div class="logo">
        <img src="/src/assets/vue.svg" alt="Logo" class="logo-img">
        <h1 class="logo-text">LingProxy</h1>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="layout-menu"
        router
        :collapse-transition="false"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataLine /></el-icon>
          <template #title>仪表盘</template>
        </el-menu-item>
        <el-menu-item index="/tokens">
          <el-icon><Key /></el-icon>
          <template #title>Token管理</template>
        </el-menu-item>
        <el-menu-item index="/llm-resources">
          <el-icon><Cpu /></el-icon>
          <template #title>LLM资源管理</template>
        </el-menu-item>
        <el-menu-item index="/models">
          <el-icon><Grid /></el-icon>
          <template #title>模型管理</template>
        </el-menu-item>
        <el-menu-item index="/requests">
          <el-icon><Message /></el-icon>
          <template #title>请求管理</template>
        </el-menu-item>
        <el-menu-item index="/policies">
          <el-icon><Operation /></el-icon>
          <template #title>策略管理</template>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <template #title>系统设置</template>
        </el-menu-item>
        <el-menu-item index="/logs">
          <el-icon><Document /></el-icon>
          <template #title>日志管理</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 主内容区 -->
    <el-container class="layout-container">
      <!-- 顶部导航 -->
      <el-header class="layout-header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>{{ currentTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="user-info">
              <el-icon class="user-icon"><UserFilled /></el-icon>
              <span>{{ username }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleLogout">
                  <el-icon><Switch /></el-icon>
                  <span>退出登录</span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 内容区域 -->
      <el-main class="layout-main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  DataLine,
  Key,
  Cpu,
  Grid,
  Message,
  Operation,
  Setting,
  UserFilled,
  Switch,
  Document
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const username = ref('管理员')

// 计算当前活动菜单
const activeMenu = computed(() => {
  return route.path
})

// 计算当前页面标题
const currentTitle = computed(() => {
  const matched = route.matched
  if (matched.length > 0) {
    const lastMatched = matched[matched.length - 1]
    return lastMatched.meta.title || '首页'
  }
  return '首页'
})

// 处理退出登录
const handleLogout = () => {
  localStorage.removeItem('token')
  router.push('/login')
  ElMessage.success('退出登录成功')
}

// 组件挂载时获取用户信息
onMounted(() => {
  // 这里可以从后端获取用户信息
  const userInfo = localStorage.getItem('userInfo')
  if (userInfo) {
    try {
      const info = JSON.parse(userInfo)
      username.value = info.username || '管理员'
    } catch (error) {
      console.error('解析用户信息失败:', error)
    }
  }
})
</script>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.layout-aside {
  background-color: #001529;
  color: #fff;
  overflow-y: auto;
}

.logo {
  display: flex;
  align-items: center;
  padding: 0 20px;
  height: 64px;
  border-bottom: 1px solid #1f2d3d;
}

.logo-img {
  width: 32px;
  height: 32px;
  margin-right: 10px;
}

.logo-text {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
}

.layout-menu {
  border-right: none;
  background-color: transparent;
}

.layout-menu :deep(.el-menu-item) {
  color: rgba(255, 255, 255, 0.85);
  height: 60px;
  line-height: 60px;
  margin: 0 10px;
  border-radius: 6px;
}

.layout-menu :deep(.el-menu-item:hover) {
  background-color: rgba(255, 255, 255, 0.1);
}

.layout-menu :deep(.el-menu-item.is-active) {
  background-color: rgba(255, 255, 255, 0.2);
  color: #fff;
}

.layout-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.layout-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  background-color: #fff;
  border-bottom: 1px solid #e8e8e8;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-left {
  flex: 1;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.user-icon {
  margin-right: 8px;
}

.layout-main {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #f5f5f5;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .layout-aside {
    width: 64px !important;
  }

  .logo-text {
    display: none;
  }

  .layout-menu :deep(.el-menu-item__title) {
    display: none;
  }

  .layout-menu :deep(.el-menu-item) {
    justify-content: center;
  }
}
</style>