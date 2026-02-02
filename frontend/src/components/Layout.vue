<template>
  <div class="layout">
    <!-- 侧边栏 -->
    <el-aside width="200px" class="layout-aside">
      <div class="logo">
        <img src="@/assets/lingproxy-logo.svg" alt="LingProxy Logo" class="logo-img">
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
import { ElMessage, ElMessageBox } from 'element-plus'
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
import { STORAGE_KEYS } from '../utils/constants'

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
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    localStorage.removeItem(STORAGE_KEYS.TOKEN)
    localStorage.removeItem(STORAGE_KEYS.USER_INFO)
    router.push('/login')
    ElMessage.success('退出登录成功')
  } catch {
    // 用户取消操作
  }
}

// 组件挂载时获取用户信息
onMounted(() => {
  // 从localStorage获取用户信息
  const userInfo = localStorage.getItem(STORAGE_KEYS.USER_INFO)
  if (userInfo) {
    try {
      const info = JSON.parse(userInfo)
      username.value = info.username || info.name || '管理员'
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
  background: linear-gradient(180deg, #1e293b 0%, #0f172a 100%);
  color: #fff;
  overflow-y: auto;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  padding: 0 20px;
  height: 64px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.03);
  transition: all 0.3s ease;
}

.logo:hover {
  background: rgba(255, 255, 255, 0.05);
}

.logo-img {
  width: 40px;
  height: 40px;
  margin-right: 12px;
  flex-shrink: 0;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  margin: 0;
  background: linear-gradient(135deg, #2563eb 0%, #7c3aed 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.layout-menu {
  border-right: none;
  background-color: transparent;
}

.layout-menu :deep(.el-menu-item) {
  color: rgba(255, 255, 255, 0.85);
  height: 56px;
  line-height: 56px;
  margin: 4px 12px;
  border-radius: 8px;
  transition: all 0.3s ease;
  font-size: 14px;
}

.layout-menu :deep(.el-menu-item:hover) {
  background: linear-gradient(90deg, rgba(37, 99, 235, 0.2) 0%, rgba(124, 58, 237, 0.2) 100%);
  color: #fff;
  transform: translateX(4px);
}

.layout-menu :deep(.el-menu-item.is-active) {
  background: linear-gradient(90deg, rgba(37, 99, 235, 0.3) 0%, rgba(124, 58, 237, 0.3) 100%);
  color: #fff;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.3);
}

.layout-menu :deep(.el-menu-item .el-icon) {
  font-size: 18px;
  margin-right: 8px;
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