import { createI18n } from 'vue-i18n'
import zh from './zh'
import en from './en'

const messages = {
  zh,
  en
}

// 从 localStorage 获取保存的语言设置，默认为中文
const savedLocale = localStorage.getItem('lingproxy_locale') || 'zh'

const i18n = createI18n({
  legacy: false,
  locale: savedLocale,
  fallbackLocale: 'zh',
  messages
})

export default i18n
