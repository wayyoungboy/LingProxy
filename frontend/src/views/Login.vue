<template>
  <div class="login-container">
    <div class="login-form-wrapper">
      <div class="login-logo">
        <img src="/src/assets/vue.svg" alt="Logo" class="logo-img">
        <h1 class="logo-text">LingProxy</h1>
        <p class="logo-desc">智能LLM代理管理系统</p>
      </div>
      
      <el-card class="login-form-card">
        <template #header>
          <div class="login-form-header">
            <h2>用户登录</h2>
          </div>
        </template>
        
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          label-width="80px"
          class="login-form"
        >
          <el-form-item label="用户名" prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入用户名"
              prefix-icon="User"
              size="large"
            ></el-input>
          </el-form-item>
          
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              prefix-icon="Lock"
              show-password
              size="large"
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
              登录
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
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import api from '../api'

const router = useRouter()
const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  try {
    // 验证表单
    await loginFormRef.value.validate()
    
    loading.value = true
    
    // 调用真实登录API
    const response = await api.login({
      username: loginForm.username,
      password: loginForm.password
    })
    
    if (response && response.data) {
      // 存储token和用户信息
      localStorage.setItem('token', response.data.token)
      localStorage.setItem('userInfo', JSON.stringify(response.data.user))
      
      ElMessage.success('登录成功')
      // 跳转到首页
      router.push('/')
    }
  } catch (error) {
    console.error('登录失败:', error)
    const errorMsg = error.response?.data?.error || '登录失败，请检查用户名和密码'
    ElMessage.error(errorMsg)
  } finally {
    loading.value = false
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
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  padding: 12px;
}

.logo-text {
  font-size: 28px;
  font-weight: bold;
  color: #fff;
  margin: 0 0 8px 0;
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