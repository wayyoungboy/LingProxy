import { createRouter, createWebHistory } from 'vue-router'
import { STORAGE_KEYS, ROUTE_PATHS } from '../utils/constants'
import { getMenuTitle } from '../config/menu'
import i18n from '../locales'

const routes = [
  {
    path: ROUTE_PATHS.LOGIN,
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { 
      requiresAuth: false,
      titleKey: 'login.title'
    }
  },
  {
    path: '/',
    name: 'Home',
    redirect: ROUTE_PATHS.DASHBOARD,
    component: () => import('../components/Layout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { titleKey: 'menu.dashboard' }
      },
      {
        path: 'tokens',
        name: 'Tokens',
        component: () => import('../views/Tokens.vue'),
        meta: { titleKey: 'menu.tokens' }
      },
      {
        path: 'llm-resources/list',
        name: 'LLMResources',
        component: () => import('../views/LLMResources.vue'),
        meta: { titleKey: 'menu.llmResources' }
      },
      {
        path: 'llm-resources/usage',
        name: 'LLMResourceUsage',
        component: () => import('../views/LLMResourceUsage.vue'),
        meta: { titleKey: 'menu.llmResourceUsage' }
      },
      {
        path: 'requests',
        name: 'Requests',
        component: () => import('../views/Requests.vue'),
        meta: { titleKey: 'menu.requests' }
      },
      {
        path: 'policies',
        name: 'Policies',
        component: () => import('../views/Policies.vue'),
        meta: { titleKey: 'menu.policies' }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../views/Settings.vue'),
        meta: { titleKey: 'menu.settings' }
      },
      {
        path: 'logs',
        name: 'Logs',
        component: () => import('../views/Logs.vue'),
        meta: { titleKey: 'menu.logs' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: ROUTE_PATHS.DASHBOARD
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const token = localStorage.getItem(STORAGE_KEYS.TOKEN)

  // 设置页面标题（使用国际化）
  if (to.meta.titleKey) {
    const title = i18n.global.t(to.meta.titleKey)
    document.title = `${title} - LingProxy`
  } else if (to.meta.title) {
    // 兼容旧代码
    document.title = `${to.meta.title} - LingProxy`
  }

  if (requiresAuth && !token) {
    // 需要认证但未登录，跳转到登录页
    next({
      path: ROUTE_PATHS.LOGIN,
      query: { redirect: to.fullPath }
    })
  } else if (to.path === ROUTE_PATHS.LOGIN && token) {
    // 已登录访问登录页，跳转到首页
    next(ROUTE_PATHS.DASHBOARD)
  } else {
    next()
  }
})

export default router