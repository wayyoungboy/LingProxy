<template>
  <div class="layout">
    <!-- 侧边栏 -->
    <Sidebar />

    <!-- 主内容区 -->
    <el-container class="layout-container">
      <!-- 顶部导航 -->
      <el-header class="layout-header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">{{ $t('common.home') }}</el-breadcrumb-item>
            <el-breadcrumb-item>{{ currentTitle }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <!-- 语言切换 -->
          <el-dropdown @command="handleLanguageChange" class="language-dropdown">
            <span class="language-selector">
              <el-icon><Setting /></el-icon>
              <span>{{ currentLanguageLabel }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="zh" :class="{ 'is-active': currentLocale === 'zh' }">
                  中文
                </el-dropdown-item>
                <el-dropdown-item command="en" :class="{ 'is-active': currentLocale === 'en' }">
                  English
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>

          <!-- 用户信息 -->
          <el-dropdown>
            <span class="user-info">
              <el-icon class="user-icon"><UserFilled /></el-icon>
              <span>{{ username }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleLogout">
                  <el-icon><Switch /></el-icon>
                  <span>{{ $t('common.logout') }}</span>
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UserFilled, Switch, Setting } from '@element-plus/icons-vue'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import Sidebar from './Sidebar.vue'
import { STORAGE_KEYS } from '../utils/constants'
import { getMenuTitle } from '../config/menu'

const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()
const username = ref('管理员')

// 当前语言
const currentLocale = computed(() => locale.value)

// 当前语言标签
const currentLanguageLabel = computed(() => {
  return locale.value === 'zh' ? '中文' : 'English'
})

// 计算当前页面标题
const currentTitle = computed(() => getMenuTitle(route.path, t))

// 处理语言切换
const handleLanguageChange = lang => {
  locale.value = lang
  localStorage.setItem('lingproxy_locale', lang)

  // 更新 Element Plus 的语言
  ElementPlus.locale(lang === 'en' ? en : zhCn)

  // 更新页面标题
  const titleKey = route.meta.titleKey || 'common.home'
  document.title = `${t(titleKey)} - LingProxy`

  ElMessage.success(t('common.languageSwitched'))
}

// 监听语言变化，更新 Element Plus 的语言
watch(
  () => locale.value,
  newLang => {
    ElementPlus.locale(newLang === 'en' ? en : zhCn)
  }
)

// 处理退出登录
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm(t('common.confirmLogout'), t('common.info'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    localStorage.removeItem(STORAGE_KEYS.TOKEN)
    localStorage.removeItem(STORAGE_KEYS.USER_INFO)
    router.push('/login')
    ElMessage.success(t('common.logoutSuccess'))
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
/* Claude Layout */
.layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
  background-color: var(--claude-parchment);
}

.layout-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}

/* Claude Header - Ivory with warm borders */
.layout-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  height: 72px !important;
  background-color: var(--claude-ivory);
  border-bottom: 1px solid var(--claude-border-warm);
  z-index: 10;
}

.header-left {
  flex: 1;
}

/* Breadcrumb styling */
:deep(.el-breadcrumb__item) {
  font-size: 14px;
  font-family: var(--font-sans);
}

:deep(.el-breadcrumb__inner.is-link:hover) {
  color: var(--claude-terracotta);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* Claude interactive elements - Warm Sand hover */
.language-selector,
.user-info {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  border-radius: var(--radius-generous);
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  color: var(--claude-text-secondary);
  font-weight: 500;
  font-family: var(--font-sans);
  border: 1px solid transparent;
}

.language-selector:hover,
.user-info:hover {
  background-color: var(--claude-warm-sand);
  color: var(--claude-text-primary);
  border-color: var(--claude-border-cream);
  box-shadow: var(--shadow-whisper);
}

.language-selector .el-icon,
.user-icon {
  font-size: 18px;
  margin-right: 8px;
}

/* Claude Main - Parchment background */
.layout-main {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  background-color: var(--claude-parchment);
}

/* Page transitions */
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@media (max-width: 768px) {
  .layout-header {
    padding: 0 16px;
  }

  .language-selector span,
  .user-info span:last-child {
    display: none;
  }
}
</style>
