/**
 * 菜单配置
 * 统一管理侧边栏菜单项
 */
import {
  DataLine,
  Key,
  Cpu,
  Message,
  Operation,
  Setting,
  Document,
  DataAnalysis
} from '@element-plus/icons-vue'

/**
 * 菜单项类型
 * @typedef {Object} MenuItem
 * @property {string} index - 路由路径
 * @property {string} title - 菜单标题
 * @property {string|Component} icon - 图标组件
 * @property {string} [meta.title] - 页面标题（用于面包屑）
 */

/**
 * 菜单配置列表
 */
export const menuItems = [
  {
    index: '/dashboard',
    title: '仪表盘',
    icon: DataLine,
    meta: {
      title: '仪表盘'
    }
  },
  {
    index: '/tokens',
    title: 'Token管理',
    icon: Key,
    meta: {
      title: 'Token管理'
    }
  },
  {
    index: '/llm-resources/list',
    title: 'LLM资源管理',
    icon: Cpu,
    meta: {
      title: 'LLM资源管理'
    }
  },
  {
    index: '/llm-resources/usage',
    title: '用量统计',
    icon: DataAnalysis,
    meta: {
      title: '用量统计'
    }
  },
  {
    index: '/requests',
    title: '请求管理',
    icon: Message,
    meta: {
      title: '请求管理'
    }
  },
  {
    index: '/policies',
    title: '策略管理',
    icon: Operation,
    meta: {
      title: '策略管理'
    }
  },
  {
    index: '/settings',
    title: '系统设置',
    icon: Setting,
    meta: {
      title: '系统设置'
    }
  },
  {
    index: '/logs',
    title: '日志管理',
    icon: Document,
    meta: {
      title: '日志管理'
    }
  }
]

/**
 * 根据路径查找菜单项
 * @param {string} path - 路由路径
 * @returns {MenuItem|null}
 */
export function findMenuItemByPath(path) {
  for (const item of menuItems) {
    if (item.index === path) {
      return item
    }
  }
  return null
}

/**
 * 获取菜单项的页面标题
 * @param {string} path - 路由路径
 * @returns {string}
 */
export function getMenuTitle(path) {
  const item = findMenuItemByPath(path)
  if (item) {
    return item.meta?.title || item.title
  }
  return '首页'
}
