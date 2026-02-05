<template>
  <div class="login-container">
    <div class="login-form-wrapper">
      <div class="login-logo">
        <img src="@/assets/lingproxy-logo.svg" alt="LingProxy Logo" class="logo-img">
        <h1 class="logo-text">LingProxy</h1>
        <p class="logo-desc">{{ $t('login.subtitle') }}</p>
      </div>
      
      <el-card class="login-form-card">
        <template #header>
          <div class="login-form-header">
            <h2>{{ $t('login.title') }}</h2>
          </div>
        </template>
        
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          label-width="80px"
          class="login-form"
        >
          <el-form-item :label="$t('login.username')" prop="username">
            <el-input
              v-model="loginForm.username"
              :placeholder="$t('login.usernamePlaceholder')"
              :prefix-icon="User"
              size="large"
              @keyup.enter="handleLogin"
              clearable
            ></el-input>
          </el-form-item>
          
          <el-form-item :label="$t('login.password')" prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              :placeholder="$t('login.passwordPlaceholder')"
              :prefix-icon="Lock"
              show-password
              size="large"
              @keyup.enter="handleLogin"
              clearable
            ></el-input>
          </el-form-item>
          
          <el-form-item>
            <el-button
              type="primary"
              class="login-button"
              :loading="loading"
              @click="handleLogin"
              size="large"
            >
              {{ loading ? $t('login.loggingIn') : $t('login.loginButton') }}
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
      
      <div class="login-footer">
        <p>&copy; 2026 LingProxy. 保留所有权利。</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import api from '../api'
import { STORAGE_KEYS } from '../utils/constants'

const router = useRouter()
const { t } = useI18n()
const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = computed(() => ({
  username: [
    { required: true, message: t('login.usernameRequired'), trigger: 'blur' },
    { min: 2, message: t('login.usernameMinLength'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('login.passwordRequired'), trigger: 'blur' },
    { min: 6, message: t('login.passwordMinLength'), trigger: 'blur' }
  ]
}))

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    // 验证表单
    await loginFormRef.value.validate()
    
    loading.value = true
    
    // 调用登录API
    const response = await api.login({
      username: loginForm.username.trim(),
      password: loginForm.password
    })
    
    // 处理响应数据（根据实际API响应格式调整）
    const token = response?.token || response?.data?.token
    const userInfo = response?.user || response?.data?.user
    
    if (token) {
      // 存储token和用户信息
      localStorage.setItem(STORAGE_KEYS.TOKEN, token)
      if (userInfo) {
        localStorage.setItem(STORAGE_KEYS.USER_INFO, JSON.stringify(userInfo))
      }
      
      ElMessage.success(t('login.loginSuccess'))
      // 跳转到首页
      router.push('/')
    } else {
      ElMessage.error(t('login.loginFailed') + ': ' + t('login.noToken'))
    }
  } catch (error) {
    console.error('登录失败:', error)
    // 错误信息已在API拦截器中处理，这里不需要再次显示
    if (!error.response) {
      ElMessage.error(t('api.networkError'))
    }
  } finally {
    loading.value = false
  }
}

// 支持回车键登录
const handleKeyPress = (event) => {
  if (event.key === 'Enter') {
    handleLogin()
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-form-wrapper {
  width: 100%;
  max-width: 450px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.login-logo {
  text-align: center;
  margin-bottom: 30px;
}

.logo-img {
  width: 80px;
  height: 80px;
  margin-bottom: 20px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease;
}

.logo-img:hover {
  transform: scale(1.05);
}

.logo-text {
  font-size: 32px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px 0;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  letter-spacing: -0.5px;
}

.logo-desc {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
}

.login-form-card {
  width: 100%;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}

.login-form-header {
  text-align: center;
}

.login-form-header h2 {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  color: #333;
}

.login-form {
  padding: 0 30px 30px;
}

.login-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 500;
  margin-top: 10px;
}

.login-footer {
  margin-top: 20px;
  text-align: center;
}

.login-footer p {
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-form-wrapper {
    max-width: 100%;
  }
  
  .login-form {
    padding: 0 20px 20px;
  }
  
  .logo-text {
    font-size: 24px;
  }
  
  .logo-desc {
    font-size: 14px;
  }
}
</style>