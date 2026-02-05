<template>
  <div class="tokens-container">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="page-title">{{ $t('tokens.title') }}</span>
          <el-button type="primary" @click="handleAddToken">
            <el-icon><Plus /></el-icon>
            {{ $t('tokens.createToken') }}
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
        <el-table-column prop="name" :label="$t('tokens.name')" />
        <el-table-column prop="token" :label="$t('tokens.tokenValue')" width="280">
          <template #default="scope">
            <div style="display: flex; align-items: center; gap: 8px;">
              <el-tag>{{ scope.row.prefix || scope.row.token || 'ling-...' }}</el-tag>
              <el-button
                type="primary"
                size="small"
                text
                @click="handleCopyToken(scope.row)"
                :title="$t('tokens.copyFullToken')"
              >
                <el-icon><DocumentCopy /></el-icon>
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="100">
          <template #default="scope">
            <el-tag
              :type="scope.row.status === 'active' ? 'success' : 'danger'"
            >
              {{ scope.row.status === 'active' ? $t('tokens.active') : $t('tokens.inactive') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="policy_id" :label="$t('tokens.policy')" width="150">
          <template #default="scope">
            <el-tag v-if="scope.row.policy_id" type="info">
              {{ getPolicyName(scope.row.policy_id) }}
            </el-tag>
            <span v-else style="color: #909399">{{ $t('dashboard.notConfigured') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="last_used_at" :label="$t('tokens.lastUsedAt')" width="180">
          <template #default="scope">
            {{ scope.row.last_used_at ? formatDate(scope.row.last_used_at) : $t('dashboard.neverUsed') }}
          </template>
        </el-table-column>
        <el-table-column prop="expires_at" :label="$t('tokens.expiresAt')" width="180">
          <template #default="scope">
            {{ scope.row.expires_at ? formatDate(scope.row.expires_at) : $t('dashboard.neverExpires') }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="200">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              @click="handleEditToken(scope.row)"
              style="margin-right: 5px"
            >
              {{ $t('common.edit') }}
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="handleResetToken(scope.row.id)"
              style="margin-right: 5px"
            >
              {{ $t('tokens.reset') }}
            </el-button>
            <el-button
              type="info"
              size="small"
              @click="handleSetPolicy(scope.row)"
              style="margin-right: 5px"
            >
              {{ $t('tokens.policy') }}
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDeleteToken(scope.row.id)"
            >
              {{ $t('common.delete') }}
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
        <el-form-item :label="$t('tokens.name')" prop="name">
          <el-input v-model="tokenForm.name" :placeholder="$t('tokens.nameRequired')"></el-input>
        </el-form-item>
        <el-form-item v-if="isAddMode" :label="$t('tokens.policy')" prop="policy_id">
          <el-select
            v-model="tokenForm.policy_id"
            :placeholder="$t('tokens.policyRequired')"
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
        <el-form-item :label="$t('tokens.expiresAt')" prop="expires_at">
          <el-date-picker
            v-model="tokenForm.expires_at"
            type="datetime"
            :placeholder="$t('tokens.expiresAtPlaceholder')"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item v-if="!isAddMode" :label="$t('common.status')" prop="status">
          <el-select v-model="tokenForm.status" :placeholder="$t('tokens.selectStatus')">
            <el-option :label="$t('tokens.active')" value="active"></el-option>
            <el-option :label="$t('tokens.inactive')" value="inactive"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
          <el-button type="primary" @click="handleSaveToken">{{ $t('common.confirm') }}</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- Token显示对话框（创建/重置后显示完整Token） -->
    <el-dialog
      v-model="tokenDisplayVisible"
      :title="$t('tokens.tokenCreated')"
      width="600px"
    >
      <el-alert
        type="warning"
        :closable="false"
        style="margin-bottom: 20px"
      >
        <template #title>
          <strong>{{ $t('tokens.importantNotice') }}：</strong>{{ $t('tokens.tokenDisplayOnce') }}
        </template>
      </el-alert>
      <el-form label-width="100px">
        <el-form-item :label="$t('tokens.name')">
          <el-input :value="newToken.name" readonly></el-input>
        </el-form-item>
        <el-form-item :label="$t('tokens.token')">
          <el-input
            :value="newToken.token"
            readonly
            ref="tokenInputRef"
          >
            <template #append>
              <el-button @click="copyToken">{{ $t('tokens.copyToken') }}</el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button type="primary" @click="tokenDisplayVisible = false">{{ $t('tokens.iHaveSaved') }}</el-button>
      </template>
    </el-dialog>

    <!-- 设置策略对话框 -->
    <el-dialog
      v-model="policyDialogVisible"
      :title="$t('tokens.setTokenPolicy')"
      width="500px"
    >
      <el-form label-width="100px">
        <el-form-item :label="$t('tokens.name')">
          <el-input :value="selectedToken?.name" readonly></el-input>
        </el-form-item>
        <el-form-item :label="$t('tokens.selectPolicy')">
          <el-select
            v-model="selectedPolicyId"
            :placeholder="$t('tokens.selectPolicyPlaceholder')"
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
        <el-button @click="policyDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="danger" @click="handleRemovePolicy" v-if="selectedToken?.policy_id">
          {{ $t('tokens.removePolicy') }}
        </el-button>
        <el-button type="primary" @click="handleSavePolicy">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, DocumentCopy } from '@element-plus/icons-vue'
import api from '../api'

const { t } = useI18n()

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

const tokenRules = computed(() => ({
  name: [
    { required: true, message: t('tokens.nameRequired'), trigger: 'blur' }
  ],
  policy_id: [
    { required: true, message: t('tokens.policyRequiredSelect'), trigger: 'change' }
  ]
}))

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
    ElMessage.error(t('tokens.getListFailed'))
  } finally {
    loading.value = false
  }
}

// 处理添加Token
const handleAddToken = () => {
  isAddMode.value = true
  dialogTitle.value = t('tokens.createToken')
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
  dialogTitle.value = t('tokens.editToken')
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
            ElMessage.warning(t('tokens.createSuccessButPolicyFailed'))
          }
        }
        dialogVisible.value = false
        tokenDisplayVisible.value = true
        ElMessage.success(t('tokens.createSuccess'))
        getTokenList()
      }
    } else {
      // 更新Token
      await api.updateToken(tokenForm.id, {
        name: tokenForm.name,
        status: tokenForm.status
      })
      ElMessage.success(t('tokens.updateSuccess'))
      dialogVisible.value = false
      getTokenList()
    }
  } catch (error) {
    console.error('保存Token失败:', error)
    ElMessage.error(t('tokens.saveFailed'))
  }
}

// 处理删除Token
const handleDeleteToken = async (id) => {
  try {
    await ElMessageBox.confirm(
      t('tokens.deleteConfirmMessage'),
      t('tokens.deleteConfirm'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'danger'
      }
    )
    
    await api.deleteToken(id)
    ElMessage.success(t('tokens.deleteSuccess'))
    getTokenList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除Token失败:', error)
      ElMessage.error(t('tokens.deleteFailed'))
    }
  }
}

// 处理重置Token
const handleResetToken = async (id) => {
  try {
    await ElMessageBox.confirm(
      t('tokens.resetConfirmMessage'),
      t('tokens.resetConfirm'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    
    const response = await api.resetToken(id)
    if (response && response.data) {
      newToken.value = response.data
      tokenDisplayVisible.value = true
      ElMessage.success(t('tokens.resetSuccess'))
      getTokenList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重置Token失败:', error)
      ElMessage.error(t('tokens.resetFailed'))
    }
  }
}

// 复制Token（用于创建/重置对话框）
const copyToken = () => {
  if (tokenInputRef.value) {
    const input = tokenInputRef.value.$el.querySelector('input')
    if (input) {
      input.select()
      document.execCommand('copy')
      ElMessage.success(t('tokens.tokenCopied'))
    }
  }
}

// 复制Token（用于列表）
const handleCopyToken = async (token) => {
  try {
    // 如果列表中没有完整token，则从API获取
    let tokenValue = token.token
    if (!tokenValue || tokenValue.includes('...')) {
      const response = await api.getToken(token.id)
      if (response && response.data) {
        tokenValue = response.data.token
      } else {
        ElMessage.error(t('tokens.getTokenFailed'))
        return
      }
    }
    
    // 复制到剪贴板
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(tokenValue)
      ElMessage.success(t('tokens.tokenCopied'))
    } else {
      // 降级方案
      const textArea = document.createElement('textarea')
      textArea.value = tokenValue
      textArea.style.position = 'fixed'
      textArea.style.left = '-999999px'
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
      ElMessage.success(t('tokens.tokenCopied'))
    }
  } catch (error) {
    console.error('复制Token失败:', error)
    ElMessage.error(t('tokens.copyFailed'))
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
      ElMessage.success(t('tokens.policySetSuccess'))
    } else {
      await api.removeTokenPolicy(selectedToken.value.id)
      ElMessage.success(t('tokens.policyRemoved'))
    }
    policyDialogVisible.value = false
    getTokenList()
  } catch (error) {
    console.error('设置策略失败:', error)
    ElMessage.error(t('tokens.policySetFailed'))
  }
}

// 移除策略
const handleRemovePolicy = async () => {
  if (!selectedToken.value) return
  
  try {
    await api.removeTokenPolicy(selectedToken.value.id)
    ElMessage.success(t('tokens.policyRemoved'))
    policyDialogVisible.value = false
    getTokenList()
  } catch (error) {
    console.error('移除策略失败:', error)
    ElMessage.error(t('tokens.policyRemoveFailed'))
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
