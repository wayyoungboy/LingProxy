import { createRouter, createWebHistory } from 'vue-router'
import { STORAGE_KEYS, ROUTE_PATHS } from '../utils/constants'

const routes = [
  {
    path: ROUTE_PATHS.LOGIN,
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { 
      requiresAuth: false,
      title: '登录'
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
        meta: { title: '仪表盘' }
      },
      {
        path: 'tokens',
        name: 'Tokens',
        component: () => import('../views/Tokens.vue'),
        meta: { title: 'Token管理' }
      },
      {
        path: 'llm-resources',
        name: 'LLMResources',
        component: () => import('../views/LLMResources.vue'),
        meta: { title: 'LLM资源管理' }
      },
      {
        path: 'models',
        name: 'Models',
        component: () => import('../views/Models.vue'),
        meta: { title: '模型管理' }
      },
      {
        path: 'requests',
        name: 'Requests',
        component: () => import('../views/Requests.vue'),
        meta: { title: '请求管理' }
      },
      {
        path: 'policies',
        name: 'Policies',
        component: () => import('../views/Policies.vue'),
        meta: { title: '策略管理' }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../views/Settings.vue'),
        meta: { title: '系统设置' }
      },
      {
        path: 'logs',
        name: 'Logs',
        component: () => import('../views/Logs.vue'),
        meta: { title: '日志管理' }
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

  // 设置页面标题
  if (to.meta.title) {
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