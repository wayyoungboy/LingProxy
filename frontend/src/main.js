import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import './style.css'

const app = createApp(App)

// 配置Element Plus
app.use(ElementPlus, {
  locale: zhCn,
  size: 'default'
})

// 使用路由
app.use(router)

// 全局错误处理
app.config.errorHandler = (err, instance, info) => {
  console.error('Global error:', err, info)
}

app.mount('#app')