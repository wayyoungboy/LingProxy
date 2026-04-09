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
  DataAnalysis,
  Collection,
  User,
  Connection
} from '@element-plus/icons-vue'

/**
 * 菜单项类型
 * @typedef {Object} MenuItem
 * @property {string} index - 路由路径
 * @property {string} titleKey - 菜单标题的i18n key
 * @property {string|Component} icon - 图标组件
 * @property {string} [meta.titleKey] - 页面标题的i18n key（用于面包屑）
 */

/**
 * 菜单配置列表
 */
export const menuItems = [
  {
    index: '/dashboard',
    titleKey: 'menu.dashboard',
    icon: DataLine,
    meta: {
      titleKey: 'menu.dashboard'
    }
  },
  {
    index: '/tokens',
    titleKey: 'menu.tokens',
    icon: Key,
    meta: {
      titleKey: 'menu.tokens'
    }
  },
  {
    index: '/llm-resources/list',
    titleKey: 'menu.llmResources',
    icon: Cpu,
    meta: {
      titleKey: 'menu.llmResources'
    }
  },
  {
    index: '/llm-resources/usage',
    titleKey: 'menu.llmResourceUsage',
    icon: DataAnalysis,
    meta: {
      titleKey: 'menu.llmResourceUsage'
    }
  },
  {
    index: '/requests',
    titleKey: 'menu.requests',
    icon: Message,
    meta: {
      titleKey: 'menu.requests'
    }
  },
  {
    index: '/policies',
    titleKey: 'menu.policies',
    icon: Operation,
    meta: {
      titleKey: 'menu.policies'
    }
  },
  {
    index: '/settings',
    titleKey: 'menu.settings',
    icon: Setting,
    meta: {
      titleKey: 'menu.settings'
    }
  },
  {
    index: '/logs',
    titleKey: 'menu.logs',
    icon: Document,
    meta: {
      titleKey: 'menu.logs'
    }
  },
  {
    index: '/models',
    titleKey: 'menu.models',
    icon: Collection,
    meta: {
      titleKey: 'menu.models'
    }
  },
  {
    index: '/users',
    titleKey: 'menu.users',
    icon: User,
    meta: {
      titleKey: 'menu.users'
    }
  },
  {
    index: '/endpoints',
    titleKey: 'menu.endpoints',
    icon: Connection,
    meta: {
      titleKey: 'menu.endpoints'
    }
  },
  {
    index: '/monitoring',
    titleKey: 'menu.monitoring',
    icon: DataLine,
    meta: {
      titleKey: 'menu.monitoring'
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
 * 获取菜单项的页面标题（需要传入i18n实例）
 * @param {string} path - 路由路径
 * @param {Function} t - i18n的t函数
 * @returns {string}
 */
export function getMenuTitle(path, t) {
  const item = findMenuItemByPath(path)
  if (item && t) {
    return t(item.meta?.titleKey || item.titleKey)
  }
  // 兼容旧代码，如果没有传入t函数，返回默认值
  if (item) {
    return item.titleKey || 'common.home'
  }
  return 'common.home'
}
