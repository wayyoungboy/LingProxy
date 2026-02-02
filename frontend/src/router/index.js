import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    name: 'Home',
    redirect: '/dashboard',
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
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/dashboard'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const hasToken = localStorage.getItem('token')

  if (requiresAuth && !hasToken) {
    next('/login')
  } else if (to.path === '/login' && hasToken) {
    next('/')
  } else {
    next()
  }
})

export default router