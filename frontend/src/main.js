import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import i18n from './locales'
import './style.css'

const app = createApp(App)

// 配置Element Plus，根据当前语言设置locale
const savedLocale = localStorage.getItem('lingproxy_locale') || 'zh'
app.use(ElementPlus, {
  locale: savedLocale === 'en' ? en : zhCn,
  size: 'default'
})

// 使用i18n
app.use(i18n)

// 使用路由
app.use(router)

// 全局错误处理
app.config.errorHandler = (err, instance, info) => {
  console.error('Global error:', err, info)
}

app.mount('#app')
