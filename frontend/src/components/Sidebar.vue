<template>
  <el-aside 
    :width="isCollapsed ? '64px' : '200px'" 
    class="sidebar"
    :class="{ 'sidebar-collapsed': isCollapsed }"
  >
    <!-- Logo区域 -->
    <div class="sidebar-logo" @click="toggleCollapse">
      <img 
        v-if="!isCollapsed" 
        src="@/assets/lingproxy-logo.svg" 
        alt="LingProxy Logo" 
        class="logo-img"
      >
      <div v-else class="logo-icon">
        <el-icon><Cpu /></el-icon>
      </div>
      <h1 v-if="!isCollapsed" class="logo-text">LingProxy</h1>
    </div>

    <!-- 菜单区域 -->
    <el-menu
      :default-active="activeMenu"
      class="sidebar-menu"
      router
      :collapse="isCollapsed"
      :collapse-transition="false"
    >
      <el-menu-item
        v-for="item in menuItems"
        :key="item.index"
        :index="item.index"
      >
        <el-icon v-if="item.icon">
          <component :is="item.icon" />
        </el-icon>
        <template #title>{{ item.title }}</template>
      </el-menu-item>
    </el-menu>

    <!-- 折叠按钮 -->
    <div class="sidebar-footer">
      <el-button
        :icon="isCollapsed ? Expand : Fold"
        circle
        size="small"
        @click="toggleCollapse"
        class="collapse-btn"
      />
    </div>
  </el-aside>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Expand, Fold, Cpu } from '@element-plus/icons-vue'
import { menuItems } from '../config/menu'

const route = useRoute()
const isCollapsed = ref(false)

// 计算当前活动菜单
const activeMenu = computed(() => {
  return route.path
})

// 切换折叠状态
const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
  // 保存折叠状态到localStorage
  localStorage.setItem('sidebar-collapsed', isCollapsed.value.toString())
}

// 初始化时从localStorage读取折叠状态
onMounted(() => {
  const savedCollapsed = localStorage.getItem('sidebar-collapsed')
  if (savedCollapsed !== null) {
    isCollapsed.value = savedCollapsed === 'true'
  }
})
</script>

<style scoped>
.sidebar {
  background: linear-gradient(180deg, #1e293b 0%, #0f172a 100%);
  color: #fff;
  overflow-y: auto;
  overflow-x: hidden;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
  transition: width 0.3s ease;
  display: flex;
  flex-direction: column;
  position: relative;
}

.sidebar-collapsed {
  width: 64px !important;
}

/* Logo区域 */
.sidebar-logo {
  display: flex;
  align-items: center;
  padding: 0 20px;
  height: 64px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.03);
  transition: all 0.3s ease;
  cursor: pointer;
  user-select: none;
}

.sidebar-logo:hover {
  background: rgba(255, 255, 255, 0.05);
}

.logo-img {
  width: 40px;
  height: 40px;
  margin-right: 12px;
  flex-shrink: 0;
}

.logo-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: #2563eb;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  margin: 0;
  background: linear-gradient(135deg, #2563eb 0%, #7c3aed 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  white-space: nowrap;
  overflow: hidden;
}

.sidebar-collapsed .logo-text {
  display: none;
}

/* 菜单样式 */
.sidebar-menu {
  flex: 1;
  border-right: none;
  background-color: transparent;
  overflow-y: auto;
  overflow-x: hidden;
}

.sidebar-menu :deep(.el-menu-item) {
  color: rgba(255, 255, 255, 0.85);
  height: 56px;
  line-height: 56px;
  margin: 4px 12px;
  border-radius: 8px;
  transition: all 0.3s ease;
  font-size: 14px;
}

.sidebar-menu :deep(.el-menu-item:hover) {
  background: linear-gradient(90deg, rgba(37, 99, 235, 0.2) 0%, rgba(124, 58, 237, 0.2) 100%);
  color: #fff;
  transform: translateX(4px);
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background: linear-gradient(90deg, rgba(37, 99, 235, 0.3) 0%, rgba(124, 58, 237, 0.3) 100%);
  color: #fff;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.3);
}

.sidebar-menu :deep(.el-menu-item .el-icon) {
  font-size: 18px;
  margin-right: 8px;
  color: inherit;
}

/* 折叠状态下的样式调整 */
.sidebar-collapsed .sidebar-menu :deep(.el-menu-item) {
  padding-left: 20px !important;
  justify-content: center;
}

.sidebar-collapsed .sidebar-menu :deep(.el-menu-item__title) {
  display: none;
}

/* 底部折叠按钮 */
.sidebar-footer {
  padding: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  justify-content: center;
  align-items: center;
}

.collapse-btn {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: rgba(255, 255, 255, 0.85);
  transition: all 0.3s ease;
}

.collapse-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  transform: scale(1.1);
}

/* 滚动条样式 */
.sidebar-menu::-webkit-scrollbar {
  width: 6px;
}

.sidebar-menu::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
}

.sidebar-menu::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
}

.sidebar-menu::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}

.sidebar::-webkit-scrollbar {
  width: 6px;
}

.sidebar::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
}

.sidebar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
}

.sidebar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0;
    top: 0;
    height: 100vh;
    z-index: 1000;
    transform: translateX(0);
  }

  .sidebar-collapsed {
    transform: translateX(-100%);
  }
}
</style>
