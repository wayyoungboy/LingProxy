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
      />
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
      <el-menu-item v-for="item in menuItems" :key="item.index" :index="item.index">
        <el-icon v-if="item.icon">
          <component :is="item.icon" />
        </el-icon>
        <template #title>{{ getMenuTitle(item) }}</template>
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
import { useI18n } from 'vue-i18n'
import { Expand, Fold, Cpu } from '@element-plus/icons-vue'
import { menuItems } from '../config/menu'

const route = useRoute()
const { t } = useI18n()
const isCollapsed = ref(false)

// 计算当前活动菜单
const activeMenu = computed(() => {
  return route.path
})

// 获取菜单标题（国际化）
const getMenuTitle = item => {
  return t(item.titleKey)
}

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
/* Claude Sidebar - Anthropic Near Black */
.sidebar {
  background: var(--claude-near-black);
  color: var(--claude-text-light-on-dark);
  overflow-y: auto;
  overflow-x: hidden;
  box-shadow: none;
  border-right: 1px solid var(--claude-border-dark);
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 100;
}

.sidebar-collapsed {
  width: 64px !important;
}

/* Logo区域 - Warm dark tone */
.sidebar-logo {
  display: flex;
  align-items: center;
  padding: 0 16px;
  height: 72px;
  background: rgba(255, 255, 255, 0.02);
  transition: all 0.3s ease;
  cursor: pointer;
  user-select: none;
  border-bottom: 1px solid var(--claude-border-dark);
}

.sidebar-logo:hover {
  background: rgba(255, 255, 255, 0.05);
}

.logo-img {
  width: 32px;
  height: 32px;
  margin-right: 12px;
  flex-shrink: 0;
}

.logo-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: var(--claude-terracotta);
}

/* Claude Logo text - Warm serif */
.logo-text {
  font-family: var(--font-serif);
  font-size: 20px;
  font-weight: 500;
  margin: 0;
  color: var(--claude-ivory);
  white-space: nowrap;
  letter-spacing: -0.5px;
}

.sidebar-collapsed .logo-text {
  display: none;
}

/* Claude Menu styling */
.sidebar-menu {
  flex: 1;
  border-right: none;
  background-color: transparent;
  padding-top: 12px;
}

.sidebar-menu :deep(.el-menu-item) {
  color: var(--claude-text-light-on-dark);
  height: 50px;
  line-height: 50px;
  margin: 4px 12px;
  border-radius: var(--radius-comfortable);
  transition: all 0.2s ease;
  font-size: 14px;
  font-weight: 500;
  font-family: var(--font-sans);
}

.sidebar-menu :deep(.el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.05) !important;
  color: var(--claude-ivory);
}

/* Active state - Terracotta */
.sidebar-menu :deep(.el-menu-item.is-active) {
  background: var(--claude-terracotta) !important;
  color: var(--claude-ivory) !important;
  box-shadow: 0px 4px 12px rgba(201, 100, 66, 0.3);
}

.sidebar-menu :deep(.el-menu-item .el-icon) {
  font-size: 20px;
  margin-right: 12px;
  color: inherit;
}

.sidebar-collapsed .sidebar-menu :deep(.el-menu-item) {
  padding-left: 0 !important;
  display: flex;
  justify-content: center;
  margin: 4px 8px;
}

.sidebar-collapsed .sidebar-menu :deep(.el-menu-item .el-icon) {
  margin-right: 0;
}

/* 底部折叠按钮 */
.sidebar-footer {
  padding: 16px;
  border-top: 1px solid var(--claude-border-dark);
  display: flex;
  justify-content: center;
}

.collapse-btn {
  background: var(--claude-dark-surface) !important;
  border: none !important;
  color: var(--claude-text-light-on-dark) !important;
  transition: all 0.3s ease;
}

.collapse-btn:hover {
  background: rgba(255, 255, 255, 0.1) !important;
  color: var(--claude-ivory) !important;
  transform: scale(1.05);
}

/* Claude Scrollbar */
.sidebar::-webkit-scrollbar {
  width: 4px;
}

.sidebar::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.15);
  border-radius: 10px;
}
</style>