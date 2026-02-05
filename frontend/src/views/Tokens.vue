<template>
  <div class="tokens-container">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="page-title">Token管理</span>
          <el-button type="primary" @click="handleAddToken">
            <el-icon><Plus /></el-icon>
            创建Token
          </el-button>
        </div>
      </template>
      
      <!-- Token列表 -->
      <el-table
        v-loading="loading"
        :data="tokens"
        style="width: 100%"
        border
        stripe
      >
        <el-table-column prop="name" label="Token名称" />
        <el-table-column prop="token" label="Token值" width="200">
          <template #default="scope">
            <el-tag>{{ scope.row.token }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag
              :type="scope.row.status === 'active' ? 'success' : 'danger'"
            >
              {{ scope.row.status === 'active' ? '活跃' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="policy_id" label="策略" width="150">
          <template #default="scope">
            <el-tag v-if="scope.row.policy_id" type="info">
              {{ getPolicyName(scope.row.policy_id) }}
            </el-tag>
            <span v-else style="color: #909399">未配置</span>
          </template>
        </el-table-column>
        <el-table-column prop="last_used_at" label="最后使用时间" width="180">
          <template #default="scope">
            {{ scope.row.last_used_at ? formatDate(scope.row.last_used_at) : '从未使用' }}
          </template>
        </el-table-column>
        <el-table-column prop="expires_at" label="过期时间" width="180">
          <template #default="scope">
            {{ scope.row.expires_at ? formatDate(scope.row.expires_at) : '永不过期' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              @click="handleEditToken(scope.row)"
              style="margin-right: 5px"
            >
              编辑
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="handleResetToken(scope.row.id)"
              style="margin-right: 5px"
            >
              重置
            </el-button>
            <el-button
              type="info"
              size="small"
              @click="handleSetPolicy(scope.row)"
              style="margin-right: 5px"
            >
              策略
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDeleteToken(scope.row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 创建/编辑Token对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
    >
      <el-form
        ref="tokenFormRef"
        :model="tokenForm"
        :rules="tokenRules"
        label-width="100px"
      >
        <el-form-item label="Token名称" prop="name">
          <el-input v-model="tokenForm.name" placeholder="请输入Token名称"></el-input>
        </el-form-item>
        <el-form-item v-if="isAddMode" label="策略" prop="policy_id">
          <el-select
            v-model="tokenForm.policy_id"
            placeholder="请选择策略（必选）"
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="policy in policies"
              :key="policy.id"
              :label="policy.name"
              :value="policy.id"
            >
              <span>{{ policy.name }}</span>
              <span style="color: #909399; margin-left: 10px">({{ policy.type }})</span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="过期时间" prop="expires_at">
          <el-date-picker
            v-model="tokenForm.expires_at"
            type="datetime"
            placeholder="选择过期时间（可选）"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item v-if="!isAddMode" label="状态" prop="status">
          <el-select v-model="tokenForm.status" placeholder="请选择状态">
            <el-option label="活跃" value="active"></el-option>
            <el-option label="禁用" value="inactive"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSaveToken">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- Token显示对话框（创建/重置后显示完整Token） -->
    <el-dialog
      v-model="tokenDisplayVisible"
      title="Token已创建"
      width="600px"
    >
      <el-alert
        type="warning"
        :closable="false"
        style="margin-bottom: 20px"
      >
        <template #title>
          <strong>重要提示：</strong>Token值只显示一次，请妥善保存！
        </template>
      </el-alert>
      <el-form label-width="100px">
        <el-form-item label="Token名称">
          <el-input :value="newToken.name" readonly></el-input>
        </el-form-item>
        <el-form-item label="Token值">
          <el-input
            :value="newToken.token"
            readonly
            ref="tokenInputRef"
          >
            <template #append>
              <el-button @click="copyToken">复制</el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button type="primary" @click="tokenDisplayVisible = false">我已保存</el-button>
      </template>
    </el-dialog>

    <!-- 设置策略对话框 -->
    <el-dialog
      v-model="policyDialogVisible"
      title="设置Token策略"
      width="500px"
    >
      <el-form label-width="100px">
        <el-form-item label="Token名称">
          <el-input :value="selectedToken?.name" readonly></el-input>
        </el-form-item>
        <el-form-item label="选择策略">
          <el-select
            v-model="selectedPolicyId"
            placeholder="请选择策略"
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="policy in policies"
              :key="policy.id"
              :label="policy.name"
              :value="policy.id"
            >
              <span>{{ policy.name }}</span>
              <span style="color: #909399; margin-left: 10px">({{ policy.type }})</span>
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="policyDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="handleRemovePolicy" v-if="selectedToken?.policy_id">
          移除策略
        </el-button>
        <el-button type="primary" @click="handleSavePolicy">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import api from '../api'

const loading = ref(false)
const tokens = ref([])
const policies = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isAddMode = ref(false)
const tokenFormRef = ref(null)
const tokenDisplayVisible = ref(false)
const tokenInputRef = ref(null)
const newToken = ref({})
const policyDialogVisible = ref(false)
const selectedToken = ref(null)
const selectedPolicyId = ref('')

const tokenForm = reactive({
  id: '',
  name: '',
  policy_id: '',
  expires_at: '',
  status: 'active'
})

const tokenRules = {
  name: [
    { required: true, message: '请输入Token名称', trigger: 'blur' }
  ],
  policy_id: [
    { required: true, message: '请选择策略', trigger: 'change' }
  ]
}

// 获取Token列表
const getTokenList = async () => {
  try {
    loading.value = true
    const response = await api.getTokens()
    if (response && response.data) {
      tokens.value = response.data.items || []
    }
  } catch (error) {
    console.error('获取Token列表失败:', error)
    ElMessage.error('获取Token列表失败')
  } finally {
    loading.value = false
  }
}

// 处理添加Token
const handleAddToken = () => {
  isAddMode.value = true
  dialogTitle.value = '创建Token'
  Object.assign(tokenForm, {
    id: '',
    name: '',
    policy_id: '',
    expires_at: '',
    status: 'active'
  })
  // 确保策略列表已加载
  if (policies.value.length === 0) {
    getPolicyList()
  }
  dialogVisible.value = true
}

// 处理编辑Token
const handleEditToken = (token) => {
  isAddMode.value = false
  dialogTitle.value = '编辑Token'
  Object.assign(tokenForm, {
    id: token.id,
    name: token.name,
    expires_at: token.expires_at || '',
    status: token.status
  })
  dialogVisible.value = true
}

// 处理保存Token
const handleSaveToken = async () => {
  try {
    await tokenFormRef.value.validate()
    
    if (isAddMode.value) {
      // 创建Token
      const response = await api.createToken({
        name: tokenForm.name,
        expires_at: tokenForm.expires_at || undefined
      })
      if (response && response.data) {
        newToken.value = response.data
        // 创建成功后立即设置策略
        if (tokenForm.policy_id) {
          try {
            await api.setTokenPolicy(response.data.id, {
              policy_id: tokenForm.policy_id
            })
          } catch (error) {
            console.error('设置策略失败:', error)
            ElMessage.warning('Token创建成功，但设置策略失败，请稍后手动设置')
          }
        }
        dialogVisible.value = false
        tokenDisplayVisible.value = true
        ElMessage.success('Token创建成功')
        getTokenList()
      }
    } else {
      // 更新Token
      await api.updateToken(tokenForm.id, {
        name: tokenForm.name,
        status: tokenForm.status
      })
      ElMessage.success('Token更新成功')
      dialogVisible.value = false
      getTokenList()
    }
  } catch (error) {
    console.error('保存Token失败:', error)
    ElMessage.error('保存Token失败')
  }
}

// 处理删除Token
const handleDeleteToken = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个Token吗？删除后该Token将无法使用。',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'danger'
      }
    )
    
    await api.deleteToken(id)
    ElMessage.success('Token删除成功')
    getTokenList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除Token失败:', error)
      ElMessage.error('删除Token失败')
    }
  }
}

// 处理重置Token
const handleResetToken = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要重置这个Token吗？重置后原Token将立即失效。',
      '重置确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const response = await api.resetToken(id)
    if (response && response.data) {
      newToken.value = response.data
      tokenDisplayVisible.value = true
      ElMessage.success('Token重置成功')
      getTokenList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重置Token失败:', error)
      ElMessage.error('重置Token失败')
    }
  }
}

// 复制Token
const copyToken = () => {
  if (tokenInputRef.value) {
    const input = tokenInputRef.value.$el.querySelector('input')
    if (input) {
      input.select()
      document.execCommand('copy')
      ElMessage.success('Token已复制到剪贴板')
    }
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 获取策略名称
const getPolicyName = (policyId) => {
  const policy = policies.value.find(p => p.id === policyId)
  return policy ? policy.name : policyId
}

// 处理设置策略
const handleSetPolicy = (token) => {
  selectedToken.value = token
  selectedPolicyId.value = token.policy_id || ''
  policyDialogVisible.value = true
}

// 保存策略
const handleSavePolicy = async () => {
  if (!selectedToken.value) return
  
  try {
    if (selectedPolicyId.value) {
      await api.setTokenPolicy(selectedToken.value.id, {
        policy_id: selectedPolicyId.value
      })
      ElMessage.success('策略设置成功')
    } else {
      await api.removeTokenPolicy(selectedToken.value.id)
      ElMessage.success('策略已移除')
    }
    policyDialogVisible.value = false
    getTokenList()
  } catch (error) {
    console.error('设置策略失败:', error)
    ElMessage.error('设置策略失败')
  }
}

// 移除策略
const handleRemovePolicy = async () => {
  if (!selectedToken.value) return
  
  try {
    await api.removeTokenPolicy(selectedToken.value.id)
    ElMessage.success('策略已移除')
    policyDialogVisible.value = false
    getTokenList()
  } catch (error) {
    console.error('移除策略失败:', error)
    ElMessage.error('移除策略失败')
  }
}

// 获取策略列表
const getPolicyList = async () => {
  try {
    const response = await api.getPolicies()
    if (response && response.data) {
      policies.value = response.data.items || []
    }
  } catch (error) {
    console.error('获取策略列表失败:', error)
  }
}

// 组件挂载时获取数据
onMounted(() => {
  getTokenList()
  getPolicyList()
})
</script>

<style scoped>
.tokens-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
}
</style>
