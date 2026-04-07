<template>
  <div class="login-container">
    <div class="login-form-wrapper">
      <!-- Claude style logo -->
      <div class="login-logo">
        <div class="logo-icon-wrapper">
          <img src="@/assets/lingproxy-logo.svg" alt="LingProxy Logo" class="logo-img" />
        </div>
        <h1 class="logo-text">LingProxy</h1>
        <p class="logo-desc">{{ $t('login.subtitle') }}</p>
      </div>

      <!-- Claude style login card -->
      <el-card class="login-form-card">
        <template #header>
          <h2 class="login-title">{{ $t('login.title') }}</h2>
        </template>

        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          label-position="top"
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

      <!-- Footer -->
      <div class="login-footer">
        <p>&copy; 2026 LingProxy. 保留所有权利。</p>
      </div>
    </div>

    <!-- Claude style decorative illustration -->
    <div class="login-decoration">
      <div class="decoration-circle decoration-1"></div>
      <div class="decoration-circle decoration-2"></div>
      <div class="decoration-circle decoration-3"></div>
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
    await loginFormRef.value.validate()
    loading.value = true

    const response = await api.login({
      username: loginForm.username.trim(),
      password: loginForm.password
    })

    const token = response?.token || response?.data?.token
    const userInfo = response?.user || response?.data?.user

    if (token) {
      localStorage.setItem(STORAGE_KEYS.TOKEN, token)
      if (userInfo) {
        localStorage.setItem(STORAGE_KEYS.USER_INFO, JSON.stringify(userInfo))
      }
      ElMessage.success(t('login.loginSuccess'))
      router.push('/')
    } else {
      ElMessage.error(t('login.loginFailed') + ': ' + t('login.noToken'))
    }
  } catch (error) {
    console.error('登录失败:', error)
    if (!error.response) {
      ElMessage.error(t('api.networkError'))
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* Claude Style Login Container */
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: var(--claude-parchment);
  padding: 40px 20px;
  position: relative;
  overflow: hidden;
}

/* Claude style decorative circles */
.login-decoration {
  position: absolute;
  inset: 0;
  overflow: hidden;
  z-index: 0;
}

.decoration-circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.3;
}

.decoration-1 {
  width: 300px;
  height: 300px;
  background: var(--claude-terracotta);
  top: -150px;
  right: -100px;
}

.decoration-2 {
  width: 200px;
  height: 200px;
  background: var(--claude-coral);
  bottom: -100px;
  left: -50px;
}

.decoration-3 {
  width: 150px;
  height: 150px;
  background: #7a9a6d;
  top: 50%;
  left: 10%;
}

.login-form-wrapper {
  width: 100%;
  max-width: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
  z-index: 1;
}

/* Claude Style Logo */
.login-logo {
  text-align: center;
  margin-bottom: 32px;
}

.logo-icon-wrapper {
  width: 72px;
  height: 72px;
  margin: 0 auto 16px;
  background: var(--claude-ivory);
  border-radius: var(--radius-very);
  padding: 16px;
  box-shadow: var(--shadow-whisper);
  border: 1px solid var(--claude-border-cream);
  transition: all 0.3s ease;
}

.logo-icon-wrapper:hover {
  transform: scale(1.05);
  box-shadow: rgba(0, 0, 0, 0.08) 0px 8px 32px;
}

.logo-img {
  width: 40px;
  height: 40px;
}

.logo-text {
  font-family: var(--font-serif);
  font-size: 36px;
  font-weight: 500;
  line-height: 1.1;
  color: var(--claude-text-primary);
  margin: 0 0 8px 0;
}

.logo-desc {
  font-family: var(--font-sans);
  font-size: 16px;
  color: var(--claude-text-secondary);
  margin: 0;
  line-height: 1.6;
}

/* Claude Style Login Card */
.login-form-card {
  width: 100%;
  border-radius: var(--radius-very) !important;
  background: var(--claude-ivory) !important;
  border: 1px solid var(--claude-border-cream) !important;
  box-shadow: var(--shadow-whisper) !important;
}

.login-title {
  font-family: var(--font-serif);
  font-size: 24px;
  font-weight: 500;
  line-height: 1.2;
  color: var(--claude-text-primary);
  margin: 0;
  text-align: center;
}

.login-form {
  padding: 0;
}

/* Claude Style Login Button */
.login-button {
  width: 100%;
  height: 48px !important;
  font-size: 16px !important;
  font-weight: 500 !important;
  font-family: var(--font-sans) !important;
  border-radius: var(--radius-generous) !important;
  margin-top: 16px;
}

/* Claude Style Form Labels */
:deep(.el-form-item__label) {
  font-family: var(--font-sans);
  font-weight: 500;
  color: var(--claude-text-secondary);
  padding-bottom: 8px;
}

/* Footer */
.login-footer {
  margin-top: 24px;
  text-align: center;
}

.login-footer p {
  font-family: var(--font-sans);
  color: var(--claude-text-tertiary);
  font-size: 14px;
  margin: 0;
}

/* Responsive */
@media (max-width: 480px) {
  .login-form-wrapper {
    max-width: 100%;
  }

  .logo-text {
    font-size: 28px;
  }
}
</style>