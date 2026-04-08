<template>
  <div class="tokens-page">
    <!-- Page Header -->
    <div class="page-header">
      <h1 class="page-title">{{ $t('tokens.title') }}</h1>
      <el-button type="primary" @click="handleAddToken">
        <el-icon><Plus /></el-icon>
        {{ $t('tokens.createToken') }}
      </el-button>
    </div>

    <!-- API Key Table -->
    <el-card class="table-card">
      <el-table v-loading="loading" :data="tokens" style="width: 100%">
        <el-table-column prop="name" :label="$t('tokens.name')" />
        <el-table-column prop="api_key" :label="$t('tokens.tokenValue')" width="280">
          <template #default="scope">
            <div class="token-cell">
              <el-tag size="small">
                {{ scope.row.prefix || scope.row.api_key || 'ling-...' }}
              </el-tag>
              <el-button type="primary" size="small" text @click="handleCopyToken(scope.row)">
                <el-icon><DocumentCopy /></el-icon>
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('common.status')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ scope.row.status === 'active' ? $t('tokens.active') : $t('tokens.inactive') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('tokens.allowedModels')" width="200">
          <template #default="scope">
            <div v-if="scope.row.allowed_models && scope.row.allowed_models.length > 0">
              <el-tag
                v-for="(model, index) in scope.row.allowed_models.slice(0, 2)"
                :key="index"
                size="small"
                class="model-tag"
              >
                {{ model }}
              </el-tag>
              <el-tag v-if="scope.row.allowed_models.length > 2" size="small" type="info">
                +{{ scope.row.allowed_models.length - 2 }}
              </el-tag>
            </div>
            <span v-else class="text-muted">{{ $t('tokens.allModelsAllowed') }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('tokens.policies')" width="250">
          <template #default="scope">
            <div class="policy-tags">
              <el-tag v-if="scope.row.chat_policy_id" size="small" type="success">
                Chat: {{ getPolicyName(scope.row.chat_policy_id) }}
              </el-tag>
              <el-tag v-if="scope.row.embedding_policy_id" size="small" type="warning">
                Embedding: {{ getPolicyName(scope.row.embedding_policy_id) }}
              </el-tag>
              <el-tag v-if="scope.row.rerank_policy_id" size="small" type="info">
                Rerank: {{ getPolicyName(scope.row.rerank_policy_id) }}
              </el-tag>
              <span
                v-if="
                  !scope.row.policy_id &&
                  !scope.row.chat_policy_id &&
                  !scope.row.embedding_policy_id &&
                  !scope.row.rerank_policy_id
                "
                class="text-muted"
              >
                {{ $t('dashboard.notConfigured') }}
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="last_used_at" :label="$t('tokens.lastUsedAt')" width="180">
          <template #default="scope">
            <span class="text-caption">
              {{
                scope.row.last_used_at
                  ? formatDate(scope.row.last_used_at)
                  : $t('dashboard.neverUsed')
              }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="expires_at" :label="$t('tokens.expiresAt')" width="180">
          <template #default="scope">
            <span class="text-caption">
              {{
                scope.row.expires_at
                  ? formatDate(scope.row.expires_at)
                  : $t('dashboard.neverExpires')
              }}
            </span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="200">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEditToken(scope.row)">
              {{ $t('common.edit') }}
            </el-button>
            <el-button type="default" size="small" @click="handleResetToken(scope.row.id)">
              {{ $t('tokens.reset') }}
            </el-button>
            <el-button type="info" size="small" @click="handleSetPolicy(scope.row)">
              {{ $t('tokens.policy') }}
            </el-button>
            <el-button type="danger" size="small" @click="handleDeleteToken(scope.row.id)">
              {{ $t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Create/Edit Dialog -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px">
      <el-form ref="tokenFormRef" :model="tokenForm" :rules="tokenRules" label-width="120px">
        <el-form-item :label="$t('tokens.name')" prop="name">
          <el-input v-model="tokenForm.name" :placeholder="$t('tokens.nameRequired')"></el-input>
        </el-form-item>

        <el-form-item :label="$t('tokens.allowedModels')">
          <el-select
            v-model="tokenForm.allowed_models"
            multiple
            filterable
            :placeholder="$t('tokens.allowedModelsPlaceholder')"
            style="width: 100%"
          >
            <el-option
              v-for="model in availableModels"
              :key="model.model_id || model.id"
              :label="`${model.name} (${model.model_id || model.id})`"
              :value="model.model_id || model.id"
            />
          </el-select>
          <div class="form-hint">{{ $t('tokens.allowedModelsHint') }}</div>
        </el-form-item>

        <el-divider content-position="left">{{ $t('tokens.typeSpecificPolicies') }}</el-divider>

        <el-form-item :label="$t('tokens.chatPolicy')">
          <el-select
            v-model="tokenForm.chat_policy_id"
            :placeholder="$t('tokens.selectPolicyOptional')"
            style="width: 100%"
            filterable
            clearable
          >
            <el-option
              v-for="policy in policies"
              :key="policy.id"
              :label="policy.name"
              :value="policy.id"
            >
              <span>{{ policy.name }}</span>
              <span class="option-hint">({{ policy.type }})</span>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item :label="$t('tokens.embeddingPolicy')">
          <el-select
            v-model="tokenForm.embedding_policy_id"
            :placeholder="$t('tokens.selectPolicyOptional')"
            style="width: 100%"
            filterable
            clearable
          >
            <el-option
              v-for="policy in policies"
              :key="policy.id"
              :label="policy.name"
              :value="policy.id"
            >
              <span>{{ policy.name }}</span>
              <span class="option-hint">({{ policy.type }})</span>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item :label="$t('tokens.rerankPolicy')">
          <el-select
            v-model="tokenForm.rerank_policy_id"
            :placeholder="$t('tokens.selectPolicyOptional')"
            style="width: 100%"
            filterable
            clearable
          >
            <el-option
              v-for="policy in policies"
              :key="policy.id"
              :label="policy.name"
              :value="policy.id"
            >
              <span>{{ policy.name }}</span>
              <span class="option-hint">({{ policy.type }})</span>
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

    <!-- Token Display Dialog -->
    <el-dialog v-model="tokenDisplayVisible" :title="$t('tokens.tokenCreated')" width="600px">
      <el-alert type="warning" :closable="false" style="margin-bottom: 20px">
        <template #title>
          <strong>{{ $t('tokens.importantNotice') }}：</strong>
          {{ $t('tokens.tokenDisplayOnce') }}
        </template>
      </el-alert>
      <el-form label-width="100px">
        <el-form-item :label="$t('tokens.name')">
          <el-input :value="newToken.name" readonly></el-input>
        </el-form-item>
        <el-form-item :label="$t('tokens.token')">
          <el-input :value="newToken.api_key || newToken.token" readonly ref="tokenInputRef">
            <template #append>
              <el-button @click="copyToken">{{ $t('tokens.copyToken') }}</el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button type="primary" @click="tokenDisplayVisible = false">
          {{ $t('tokens.iHaveSaved') }}
        </el-button>
      </template>
    </el-dialog>

    <!-- Policy Dialog -->
    <el-dialog v-model="policyDialogVisible" :title="$t('tokens.setTokenPolicy')" width="500px">
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
              <span class="option-hint">({{ policy.type }})</span>
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
import { formatDate } from '../utils/index'

const { t } = useI18n()

const loading = ref(false)
const tokens = ref([])
const policies = ref([])
const availableModels = ref([])
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
  allowed_models: [],
  chat_policy_id: '',
  embedding_policy_id: '',
  rerank_policy_id: '',
  image_policy_id: '',
  audio_policy_id: '',
  video_policy_id: '',
  expires_at: '',
  status: 'active'
})

const tokenRules = computed(() => ({
  name: [{ required: true, message: t('tokens.nameRequired'), trigger: 'blur' }]
}))

const getTokenList = async () => {
  try {
    loading.value = true
    const response = await api.getAPIKeys()
    if (response && response.data) {
      tokens.value = response.data.items || []
    }
  } catch (error) {
    console.error('获取API Key列表失败:', error)
    ElMessage.error(t('tokens.getListFailed'))
  } finally {
    loading.value = false
  }
}

const handleAddToken = () => {
  isAddMode.value = true
  dialogTitle.value = t('tokens.createToken')
  Object.assign(tokenForm, {
    id: '',
    name: '',
    policy_id: '',
    allowed_models: [],
    chat_policy_id: '',
    embedding_policy_id: '',
    rerank_policy_id: '',
    image_policy_id: '',
    audio_policy_id: '',
    video_policy_id: '',
    expires_at: '',
    status: 'active'
  })
  if (policies.value.length === 0) getPolicyList()
  if (availableModels.value.length === 0) getModelList()
  dialogVisible.value = true
}

const handleEditToken = token => {
  isAddMode.value = false
  dialogTitle.value = t('tokens.editToken')
  Object.assign(tokenForm, {
    id: token.id,
    name: token.name,
    policy_id: token.policy_id || '',
    allowed_models: token.allowed_models || [],
    chat_policy_id: token.chat_policy_id || '',
    embedding_policy_id: token.embedding_policy_id || '',
    rerank_policy_id: token.rerank_policy_id || '',
    image_policy_id: token.image_policy_id || '',
    audio_policy_id: token.audio_policy_id || '',
    video_policy_id: token.video_policy_id || '',
    expires_at: token.expires_at || '',
    status: token.status
  })
  if (policies.value.length === 0) getPolicyList()
  if (availableModels.value.length === 0) getModelList()
  dialogVisible.value = true
}

const handleSaveToken = async () => {
  try {
    await tokenFormRef.value.validate()

    if (isAddMode.value) {
      const response = await api.createAPIKey({
        name: tokenForm.name,
        expires_at: tokenForm.expires_at || undefined,
        allowed_models: tokenForm.allowed_models.length > 0 ? tokenForm.allowed_models : undefined,
        chat_policy_id: tokenForm.chat_policy_id || undefined,
        embedding_policy_id: tokenForm.embedding_policy_id || undefined,
        rerank_policy_id: tokenForm.rerank_policy_id || undefined,
        image_policy_id: tokenForm.image_policy_id || undefined,
        audio_policy_id: tokenForm.audio_policy_id || undefined,
        video_policy_id: tokenForm.video_policy_id || undefined
      })
      if (response && response.data) {
        newToken.value = response.data
        dialogVisible.value = false
        tokenDisplayVisible.value = true
        ElMessage.success(t('tokens.createSuccess'))
        getTokenList()
      }
    } else {
      const updateData = {
        name: tokenForm.name,
        status: tokenForm.status,
        allowed_models: tokenForm.allowed_models || []
      }
      if (tokenForm.chat_policy_id) updateData.chat_policy_id = tokenForm.chat_policy_id
      if (tokenForm.embedding_policy_id)
        updateData.embedding_policy_id = tokenForm.embedding_policy_id
      if (tokenForm.rerank_policy_id) updateData.rerank_policy_id = tokenForm.rerank_policy_id

      await api.updateAPIKey(tokenForm.id, updateData)
      ElMessage.success(t('tokens.updateSuccess'))
      dialogVisible.value = false
      getTokenList()
    }
  } catch (error) {
    console.error('保存API Key失败:', error)
    ElMessage.error(t('tokens.saveFailed'))
  }
}

const handleDeleteToken = async id => {
  try {
    await ElMessageBox.confirm(t('tokens.deleteConfirmMessage'), t('tokens.deleteConfirm'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'danger'
    })
    await api.deleteAPIKey(id)
    ElMessage.success(t('tokens.deleteSuccess'))
    getTokenList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除API Key失败:', error)
      ElMessage.error(t('tokens.deleteFailed'))
    }
  }
}

const handleResetToken = async id => {
  try {
    await ElMessageBox.confirm(t('tokens.resetConfirmMessage'), t('tokens.resetConfirm'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    const response = await api.resetAPIKey(id)
    if (response && response.data) {
      newToken.value = response.data
      tokenDisplayVisible.value = true
      ElMessage.success(t('tokens.resetSuccess'))
      getTokenList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重置API Key失败:', error)
      ElMessage.error(t('tokens.resetFailed'))
    }
  }
}

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

const handleCopyToken = async token => {
  try {
    let tokenValue = token.api_key || token.token
    if (!tokenValue || tokenValue.includes('...')) {
      const response = await api.getAPIKey(token.id)
      if (response && response.data) {
        tokenValue = response.data.api_key || response.data.token
      } else {
        ElMessage.error(t('tokens.getTokenFailed'))
        return
      }
    }
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(tokenValue)
      ElMessage.success(t('tokens.tokenCopied'))
    } else {
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
    console.error('复制API Key失败:', error)
    ElMessage.error(t('tokens.copyFailed'))
  }
}

const getPolicyName = policyId => {
  const policy = policies.value.find(p => p.id === policyId)
  return policy ? policy.name : policyId
}

const handleSetPolicy = token => {
  selectedToken.value = token
  selectedPolicyId.value = token.policy_id || ''
  policyDialogVisible.value = true
}

const handleSavePolicy = async () => {
  if (!selectedToken.value) return
  try {
    if (selectedPolicyId.value) {
      await api.setAPIKeyPolicy(selectedToken.value.id, { policy_id: selectedPolicyId.value })
      ElMessage.success(t('tokens.policySetSuccess'))
    } else {
      await api.removeAPIKeyPolicy(selectedToken.value.id)
      ElMessage.success(t('tokens.policyRemoved'))
    }
    policyDialogVisible.value = false
    getTokenList()
  } catch (error) {
    console.error('设置策略失败:', error)
    ElMessage.error(t('tokens.policySetFailed'))
  }
}

const handleRemovePolicy = async () => {
  if (!selectedToken.value) return
  try {
    await api.removeAPIKeyPolicy(selectedToken.value.id)
    ElMessage.success(t('tokens.policyRemoved'))
    policyDialogVisible.value = false
    getTokenList()
  } catch (error) {
    console.error('移除策略失败:', error)
    ElMessage.error(t('tokens.policyRemoveFailed'))
  }
}

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

const getModelList = async () => {
  try {
    const response = await api.getModels()
    if (response && response.data) {
      const models = Array.isArray(response.data) ? response.data : []
      availableModels.value = models.map(m => ({
        id: m.model_id || m.id,
        model_id: m.model_id || m.id,
        name: m.name || m.model_id || m.id
      }))
      const seen = new Set()
      availableModels.value = availableModels.value.filter(m => {
        if (seen.has(m.model_id)) return false
        seen.add(m.model_id)
        return true
      })
    }
  } catch (error) {
    console.error('获取模型列表失败:', error)
  }
}

onMounted(() => {
  getTokenList()
  getPolicyList()
  getModelList()
})
</script>

<style scoped>
.tokens-page {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-family: var(--font-serif);
  font-size: 32px;
  font-weight: 500;
  line-height: 1.1;
  color: var(--claude-text-primary);
}

.table-card {
  margin-bottom: 24px;
}

.token-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.model-tag {
  margin-right: 4px;
  margin-bottom: 4px;
}

.policy-tags {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.text-muted {
  color: var(--claude-text-tertiary);
  font-size: 12px;
}

.text-caption {
  font-size: 14px;
  color: var(--claude-text-secondary);
}

.form-hint {
  font-size: 12px;
  color: var(--claude-text-tertiary);
  margin-top: 4px;
}

.option-hint {
  color: var(--claude-text-tertiary);
  margin-left: 10px;
}

.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
