<template>
  <div class="policies-page">
    <div class="page-header">
      <h1 class="page-title">{{ $t('policies.title') }}</h1>
      <el-button type="primary" @click="handleAddPolicy">
        <el-icon><Plus /></el-icon>
        {{ $t('policies.addPolicy') }}
      </el-button>
    </div>

    <!-- 策略列表 -->
    <div class="table-section">
      <el-table v-loading="loading" :data="policies" style="width: 100%">
        <el-table-column prop="name" :label="$t('policies.name')" width="200">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
            <el-tag v-if="scope.row.builtin" type="info" size="small" style="margin-left: 8px">
              {{ $t('policies.builtin') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="$t('policies.type')" width="150">
          <template #default="scope">
            <el-tag>{{ getTypeLabel(scope.row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" :label="$t('common.status')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.enabled ? 'success' : 'danger'">
              {{ scope.row.enabled ? $t('policies.enabled') : $t('policies.disabled') }}
            </el-tag>
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
              @click="handleEditPolicy(scope.row)"
              style="margin-right: 5px"
            >
              {{ $t('common.edit') }}
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDeletePolicy(scope.row.id)"
              :disabled="scope.row.builtin"
            >
              {{ $t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 创建/编辑策略对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px">
      <el-form ref="policyFormRef" :model="policyForm" :rules="policyRules" label-width="120px">
        <el-form-item :label="$t('policies.name')" prop="name">
          <el-input v-model="policyForm.name" :placeholder="$t('policies.nameRequired')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('policies.template')" prop="template_id">
          <el-select
            v-model="policyForm.template_id"
            :placeholder="$t('policies.selectTemplate')"
            style="width: 100%"
            @change="handleTemplateChange"
            :disabled="!isAddMode"
          >
            <el-option
              v-for="template in templates"
              :key="template.id"
              :label="template.name"
              :value="template.id"
            >
              <div>
                <div>{{ template.name }}</div>
                <div style="color: #909399; font-size: 12px">{{ template.description }}</div>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item v-if="selectedTemplate" :label="$t('policies.type')">
          <el-tag>{{ getTypeLabel(selectedTemplate.type) }}</el-tag>
        </el-form-item>
        <el-form-item v-if="selectedTemplate" :label="$t('policies.parameters')">
          <div v-if="selectedTemplate.type === 'random'" class="policy-params">
            <el-form-item :label="$t('policies.resourcePool')">
              <div class="resource-pool-info">
                <el-alert
                  v-if="
                    !policyForm.parameters.resources ||
                    !Array.isArray(policyForm.parameters.resources) ||
                    policyForm.parameters.resources.length === 0
                  "
                  type="info"
                  :closable="false"
                  show-icon
                >
                  <template #title>
                    <span>{{ $t('policies.resourcePoolNotConfigured') }}</span>
                  </template>
                </el-alert>
                <div v-else class="resource-pool-list">
                  <div
                    v-for="(resourceId, index) in policyForm.parameters.resources"
                    :key="`${resourceId}-${index}`"
                    class="resource-pool-item"
                  >
                    <el-tag
                      :type="getResourceStatusType(resourceId)"
                      closable
                      @close="removeRandomResource(index)"
                    >
                      {{ getResourceName(resourceId) }}
                    </el-tag>
                  </div>
                </div>
              </div>
              <div style="margin-top: 10px">
                <el-select
                  v-model="newRandomResourceId"
                  :placeholder="$t('policies.selectResourceToAdd')"
                  filterable
                  clearable
                  style="width: 300px; margin-right: 10px"
                  @change="addRandomResource"
                >
                  <el-option
                    v-for="resource in availableRandomResources"
                    :key="resource.id"
                    :label="resource.name"
                    :value="resource.id"
                  >
                    <span>{{ resource.name }}</span>
                    <el-tag
                      :type="resource.status === 'active' ? 'success' : 'danger'"
                      size="small"
                      style="margin-left: 8px"
                    >
                      {{
                        resource.status === 'active'
                          ? $t('llmResources.active')
                          : $t('llmResources.inactive')
                      }}
                    </el-tag>
                  </el-option>
                </el-select>
                <el-button
                  v-if="
                    policyForm.parameters.resources && policyForm.parameters.resources.length > 0
                  "
                  type="danger"
                  size="small"
                  @click="clearRandomResources"
                >
                  {{ $t('policies.clearResourcePool') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="$t('policies.filterActiveOnly')">
              <el-switch v-model="policyForm.parameters.filter_by_status" />
              <span style="margin-left: 10px; color: #909399; font-size: 12px">
                {{ $t('policies.filterActiveOnlyDesc') }}
              </span>
            </el-form-item>
          </div>

          <div v-else-if="selectedTemplate.type === 'round_robin'" class="policy-params">
            <el-form-item :label="$t('policies.resourceList')" required>
              <el-select
                v-model="policyForm.parameters.resources"
                multiple
                :placeholder="$t('policies.selectResourcesOrder')"
                style="width: 100%"
              >
                <el-option
                  v-for="resource in llmResources"
                  :key="resource.id"
                  :label="resource.name"
                  :value="resource.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('policies.filterActiveOnly')">
              <el-switch v-model="policyForm.parameters.filter_by_status" />
            </el-form-item>
          </div>

          <div v-else-if="selectedTemplate.type === 'weighted'" class="policy-params">
            <el-form-item :label="$t('policies.weightedConfig')" required>
              <div
                v-for="(item, index) in policyForm.parameters.resources"
                :key="index"
                style="margin-bottom: 10px"
              >
                <el-select
                  v-model="item.id"
                  :placeholder="$t('policies.selectResource')"
                  style="width: 200px; margin-right: 10px"
                >
                  <el-option
                    v-for="resource in llmResources"
                    :key="resource.id"
                    :label="resource.name"
                    :value="resource.id"
                  />
                </el-select>
                <el-input-number
                  v-model="item.weight"
                  :min="1"
                  :max="100"
                  :placeholder="$t('policies.weight')"
                  style="width: 120px"
                />
                <el-button
                  type="danger"
                  size="small"
                  @click="removeWeightedResource(index)"
                  style="margin-left: 10px"
                >
                  {{ $t('common.delete') }}
                </el-button>
              </div>
              <el-button type="primary" size="small" @click="addWeightedResource">
                {{ $t('policies.addResource') }}
              </el-button>
            </el-form-item>
          </div>

          <div v-else-if="selectedTemplate.type === 'model_match'" class="policy-params">
            <el-form-item :label="$t('policies.modelMatchRules')" required>
              <div
                v-for="(mapping, index) in policyForm.parameters.mappings"
                :key="index"
                style="margin-bottom: 10px"
              >
                <el-input
                  v-model="mapping.model_pattern"
                  :placeholder="$t('policies.modelPatternPlaceholder')"
                  style="width: 200px; margin-right: 10px"
                />
                <el-select
                  v-model="mapping.resource_id"
                  :placeholder="$t('policies.selectResource')"
                  style="width: 200px; margin-right: 10px"
                >
                  <el-option
                    v-for="resource in llmResources"
                    :key="resource.id"
                    :label="resource.name"
                    :value="resource.id"
                  />
                </el-select>
                <el-button type="danger" size="small" @click="removeModelMapping(index)">
                  {{ $t('common.delete') }}
                </el-button>
              </div>
              <el-button type="primary" size="small" @click="addModelMapping">
                {{ $t('policies.addRule') }}
              </el-button>
            </el-form-item>
            <el-form-item :label="$t('policies.defaultResource')">
              <el-select
                v-model="policyForm.parameters.default_resource_id"
                :placeholder="$t('policies.selectDefaultResource')"
                style="width: 100%"
              >
                <el-option
                  v-for="resource in llmResources"
                  :key="resource.id"
                  :label="resource.name"
                  :value="resource.id"
                />
              </el-select>
            </el-form-item>
          </div>

          <div v-else-if="selectedTemplate.type === 'regex_match'" class="policy-params">
            <el-form-item :label="$t('policies.regexMatchRules')" required>
              <div
                v-for="(rule, index) in policyForm.parameters.rules"
                :key="index"
                style="margin-bottom: 10px"
              >
                <el-input
                  v-model="rule.pattern"
                  :placeholder="$t('policies.regexPatternPlaceholder')"
                  style="width: 200px; margin-right: 10px"
                />
                <el-select
                  v-model="rule.resource_id"
                  :placeholder="$t('policies.selectResource')"
                  style="width: 200px; margin-right: 10px"
                >
                  <el-option
                    v-for="resource in llmResources"
                    :key="resource.id"
                    :label="resource.name"
                    :value="resource.id"
                  />
                </el-select>
                <el-button type="danger" size="small" @click="removeRegexRule(index)">
                  {{ $t('common.delete') }}
                </el-button>
              </div>
              <el-button type="primary" size="small" @click="addRegexRule">
                {{ $t('policies.addRule') }}
              </el-button>
            </el-form-item>
            <el-form-item :label="$t('policies.defaultResource')">
              <el-select
                v-model="policyForm.parameters.default_resource_id"
                :placeholder="$t('policies.selectDefaultResource')"
                style="width: 100%"
              >
                <el-option
                  v-for="resource in llmResources"
                  :key="resource.id"
                  :label="resource.name"
                  :value="resource.id"
                />
              </el-select>
            </el-form-item>
          </div>

          <div v-else-if="selectedTemplate.type === 'regex_model_match'" class="policy-params">
            <el-form-item :label="$t('policies.resourcePool')">
              <div class="resource-pool-info">
                <el-alert
                  v-if="
                    !policyForm.parameters.resources ||
                    !Array.isArray(policyForm.parameters.resources) ||
                    policyForm.parameters.resources.length === 0
                  "
                  type="info"
                  :closable="false"
                  show-icon
                >
                  <template #title>
                    <span>{{ $t('policies.resourcePoolNotConfigured') }}</span>
                  </template>
                </el-alert>
                <div v-else class="resource-pool-list">
                  <div
                    v-for="(resourceId, index) in policyForm.parameters.resources"
                    :key="`${resourceId}-${index}`"
                    class="resource-pool-item"
                  >
                    <el-tag
                      :type="getResourceStatusType(resourceId)"
                      closable
                      @close="removeRandomResource(index)"
                    >
                      {{ getResourceName(resourceId) }}
                    </el-tag>
                  </div>
                </div>
              </div>
              <div style="margin-top: 10px">
                <el-select
                  v-model="newRandomResourceId"
                  :placeholder="$t('policies.selectResourceToAdd')"
                  filterable
                  clearable
                  style="width: 300px; margin-right: 10px"
                  @change="addRandomResource"
                >
                  <el-option
                    v-for="resource in availableRandomResources"
                    :key="resource.id"
                    :label="resource.name"
                    :value="resource.id"
                  >
                    <span>{{ resource.name }}</span>
                    <el-tag
                      :type="resource.status === 'active' ? 'success' : 'danger'"
                      size="small"
                      style="margin-left: 8px"
                    >
                      {{
                        resource.status === 'active'
                          ? $t('llmResources.active')
                          : $t('llmResources.inactive')
                      }}
                    </el-tag>
                  </el-option>
                </el-select>
                <el-button
                  v-if="
                    policyForm.parameters.resources && policyForm.parameters.resources.length > 0
                  "
                  type="danger"
                  size="small"
                  @click="clearRandomResources"
                >
                  {{ $t('policies.clearResourcePool') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="$t('policies.filterByStatus')">
              <el-switch v-model="policyForm.parameters.filter_by_status" />
            </el-form-item>
            <el-alert
              :title="$t('policies.regexModelMatchHint')"
              type="info"
              :closable="false"
              show-icon
              style="margin-top: 10px"
            />
          </div>

          <div v-else-if="selectedTemplate.type === 'priority'" class="policy-params">
            <el-form-item :label="$t('policies.priorityConfig')" required>
              <div
                v-for="(item, index) in policyForm.parameters.resources"
                :key="index"
                style="margin-bottom: 10px"
              >
                <el-select
                  v-model="item.id"
                  :placeholder="$t('policies.selectResource')"
                  style="width: 200px; margin-right: 10px"
                >
                  <el-option
                    v-for="resource in llmResources"
                    :key="resource.id"
                    :label="resource.name"
                    :value="resource.id"
                  />
                </el-select>
                <el-input-number
                  v-model="item.priority"
                  :min="1"
                  :placeholder="$t('policies.priorityPlaceholder')"
                  style="width: 200px"
                />
                <el-button
                  type="danger"
                  size="small"
                  @click="removePriorityResource(index)"
                  style="margin-left: 10px"
                >
                  {{ $t('common.delete') }}
                </el-button>
              </div>
              <el-button type="primary" size="small" @click="addPriorityResource">
                {{ $t('policies.addResource') }}
              </el-button>
            </el-form-item>
            <el-form-item :label="$t('policies.enableFallback')">
              <el-switch v-model="policyForm.parameters.fallback_enabled" />
            </el-form-item>
          </div>

          <div v-else-if="selectedTemplate.type === 'failover'" class="policy-params">
            <el-form-item :label="$t('policies.primaryResource')" required>
              <el-select
                v-model="policyForm.parameters.primary_resource_id"
                :placeholder="$t('policies.selectPrimaryResource')"
                style="width: 100%"
              >
                <el-option
                  v-for="resource in llmResources"
                  :key="resource.id"
                  :label="resource.name"
                  :value="resource.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('policies.fallbackResources')" required>
              <el-select
                v-model="policyForm.parameters.fallback_resources"
                multiple
                :placeholder="$t('policies.selectFallbackResources')"
                style="width: 100%"
              >
                <el-option
                  v-for="resource in llmResources"
                  :key="resource.id"
                  :label="resource.name"
                  :value="resource.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('policies.enableHealthCheck')">
              <el-switch v-model="policyForm.parameters.health_check_enabled" />
            </el-form-item>
            <el-form-item :label="$t('policies.healthCheckInterval')">
              <el-input-number
                v-model="policyForm.parameters.health_check_interval"
                :min="10"
                :max="300"
              />
            </el-form-item>
          </div>
        </el-form-item>
        <el-form-item v-if="!isAddMode" :label="$t('common.status')" prop="enabled">
          <el-switch v-model="policyForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
          <el-button type="primary" @click="handleSavePolicy">{{ $t('common.confirm') }}</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import api from '../api'

const { t } = useI18n()

const loading = ref(false)
const policies = ref([])
const templates = ref([])
const llmResources = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isAddMode = ref(false)
const policyFormRef = ref(null)
const selectedTemplate = ref(null)
const newRandomResourceId = ref('')

const policyForm = reactive({
  id: '',
  name: '',
  template_id: '',
  parameters: {},
  enabled: true
})

const policyRules = computed(() => ({
  name: [{ required: true, message: t('policies.nameRequired'), trigger: 'blur' }],
  template_id: [{ required: true, message: t('policies.templateRequired'), trigger: 'change' }]
}))

// 获取策略类型标签
const getTypeLabel = type => {
  const labels = {
    random: t('policies.randomSelect'),
    round_robin: t('policies.roundRobinLB'),
    weighted: t('policies.weightedLB'),
    model_match: t('policies.modelMatch'),
    regex_match: t('policies.regexMatch'),
    regex_model_match: t('policies.regexModelMatch'),
    priority: t('policies.priorityPolicy'),
    failover: t('policies.failover')
  }
  return labels[type] || type
}

// 随机策略资源池管理
const availableRandomResources = computed(() => {
  console.log('Computing availableRandomResources:', {
    llmResourcesCount: llmResources.value.length,
    resourcesPool: policyForm.parameters.resources,
    isArray: Array.isArray(policyForm.parameters.resources)
  })

  if (!llmResources.value || llmResources.value.length === 0) {
    console.warn('No LLM resources available')
    return []
  }

  if (!policyForm.parameters.resources || !Array.isArray(policyForm.parameters.resources)) {
    return llmResources.value
  }

  const resourcePoolIds = new Set(policyForm.parameters.resources)
  const available = llmResources.value.filter(resource => !resourcePoolIds.has(resource.id))
  console.log('Available resources after filtering:', available.length, available)
  return available
})

const addRandomResource = () => {
  if (!newRandomResourceId.value) {
    return
  }
  // 确保 resources 是数组
  if (!policyForm.parameters.resources || !Array.isArray(policyForm.parameters.resources)) {
    policyForm.parameters.resources = []
  }
  // 检查资源是否已在资源池中
  if (policyForm.parameters.resources.includes(newRandomResourceId.value)) {
    ElMessage.warning(t('policies.resourceAlreadyInPool'))
    newRandomResourceId.value = ''
    return
  }
  // 使用扩展运算符创建新数组以确保响应式更新
  policyForm.parameters.resources = [...policyForm.parameters.resources, newRandomResourceId.value]
  newRandomResourceId.value = ''
  console.log('Resource added to pool:', policyForm.parameters.resources)
}

const removeRandomResource = index => {
  if (
    policyForm.parameters.resources &&
    Array.isArray(policyForm.parameters.resources) &&
    policyForm.parameters.resources.length > index
  ) {
    // 使用扩展运算符创建新数组以确保响应式更新
    policyForm.parameters.resources = policyForm.parameters.resources.filter((_, i) => i !== index)
    console.log('Resource removed from pool:', policyForm.parameters.resources)
  }
}

const clearRandomResources = () => {
  ElMessageBox.confirm(t('policies.clearResourcePoolConfirm'), t('common.info'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning'
  })
    .then(() => {
      policyForm.parameters.resources = []
      ElMessage.success(t('policies.resourcePoolCleared'))
    })
    .catch(() => {})
}

const getResourceName = resourceId => {
  const resource = llmResources.value.find(r => r.id === resourceId)
  return resource ? resource.name : resourceId
}

const getResourceStatusType = resourceId => {
  const resource = llmResources.value.find(r => r.id === resourceId)
  if (!resource) {
    return 'info'
  }
  return resource.status === 'active' ? 'success' : 'danger'
}

// 获取策略列表
const getPolicyList = async () => {
  try {
    loading.value = true
    const response = await api.getPolicies()
    if (response && response.data) {
      policies.value = response.data.items || []
    }
  } catch (error) {
    console.error('获取策略列表失败:', error)
    ElMessage.error(t('policies.getListFailed'))
  } finally {
    loading.value = false
  }
}

// 获取模板列表
const getTemplateList = async () => {
  try {
    const response = await api.getPolicyTemplates()
    if (response && response.data) {
      templates.value = response.data.items || []
    }
  } catch (error) {
    console.error('获取策略模板列表失败:', error)
  }
}

// 获取LLM资源列表
const getLLMResourceList = async () => {
  try {
    const response = await api.getLLMResources({ page: 1, page_size: 100 })
    console.log('LLM Resources API response:', response)

    // API响应拦截器已经返回了 response.data，后端返回的是 {data: resources}
    // 所以这里 response 就是 {data: resources}，其中 resources 是数组
    if (response && response.data) {
      // 后端返回的data可能是数组，也可能是分页对象
      if (Array.isArray(response.data)) {
        // 直接是数组格式
        llmResources.value = response.data
      } else if (response.data.items && Array.isArray(response.data.items)) {
        // 分页格式 {data: {items: [...]}}
        llmResources.value = response.data.items
      } else {
        llmResources.value = []
        console.warn('Unexpected response format:', response)
      }
      console.log(
        'LLM Resources loaded:',
        llmResources.value.length,
        'resources:',
        llmResources.value.map(r => ({ id: r.id, name: r.name }))
      )
    } else {
      llmResources.value = []
      console.warn('No response from LLM Resources API or no data field')
    }
  } catch (error) {
    console.error('获取LLM资源列表失败:', error)
    ElMessage.error(
      t('llmResources.getListFailed') + ': ' + (error.response?.data?.error || error.message)
    )
    llmResources.value = []
  }
}

// 处理模板变化
const handleTemplateChange = templateId => {
  const template = templates.value.find(t => t.id === templateId)
  if (template) {
    selectedTemplate.value = template
    // 初始化参数
    initParameters(template)
  }
}

// 初始化参数
const initParameters = template => {
  const defaultParams = template.default_parameters ? JSON.parse(template.default_parameters) : {}

  switch (template.type) {
    case 'random':
      policyForm.parameters = {
        resources: Array.isArray(defaultParams.resources) ? [...defaultParams.resources] : [],
        filter_by_status: defaultParams.filter_by_status !== false
      }
      newRandomResourceId.value = ''
      console.log('Random policy parameters initialized:', policyForm.parameters)
      break
    case 'round_robin':
      policyForm.parameters = {
        resources: defaultParams.resources || [],
        filter_by_status: defaultParams.filter_by_status !== false
      }
      break
    case 'weighted':
      policyForm.parameters = {
        resources: defaultParams.resources || [],
        filter_by_status: defaultParams.filter_by_status !== false
      }
      break
    case 'model_match':
      policyForm.parameters = {
        mappings: defaultParams.mappings || [],
        default_resource_id: defaultParams.default_resource_id || ''
      }
      break
    case 'regex_match':
      policyForm.parameters = {
        rules: defaultParams.rules || [],
        default_resource_id: defaultParams.default_resource_id || ''
      }
      break
    case 'regex_model_match':
      policyForm.parameters = {
        resources: defaultParams.resources || [],
        filter_by_status: defaultParams.filter_by_status !== false
      }
      break
    case 'priority':
      policyForm.parameters = {
        resources: defaultParams.resources || [],
        fallback_enabled: defaultParams.fallback_enabled !== false
      }
      break
    case 'failover':
      policyForm.parameters = {
        primary_resource_id: defaultParams.primary_resource_id || '',
        fallback_resources: defaultParams.fallback_resources || [],
        health_check_enabled: defaultParams.health_check_enabled !== false,
        health_check_interval: defaultParams.health_check_interval || 30
      }
      break
    default:
      policyForm.parameters = defaultParams
  }
}

// 处理添加策略
const handleAddPolicy = async () => {
  isAddMode.value = true
  dialogTitle.value = t('policies.createPolicy')
  Object.assign(policyForm, {
    id: '',
    name: '',
    template_id: '',
    parameters: {},
    enabled: true
  })
  selectedTemplate.value = null
  newRandomResourceId.value = ''

  // 确保资源列表已加载
  if (!llmResources.value || llmResources.value.length === 0) {
    await getLLMResourceList()
  }

  dialogVisible.value = true
}

// 处理编辑策略
const handleEditPolicy = async policy => {
  isAddMode.value = false
  dialogTitle.value = t('policies.editPolicy')

  // 确保资源列表已加载
  if (!llmResources.value || llmResources.value.length === 0) {
    await getLLMResourceList()
  }

  // 获取策略详情
  try {
    const response = await api.getPolicy(policy.id)
    if (response && response.data) {
      // 深度复制 parameters 以确保响应式
      const parameters = response.data.parameters
        ? JSON.parse(JSON.stringify(response.data.parameters))
        : {}

      Object.assign(policyForm, {
        id: response.data.id,
        name: response.data.name,
        template_id: response.data.template_id,
        parameters: parameters,
        enabled: response.data.enabled
      })

      // 设置选中的模板
      const template = templates.value.find(t => t.id === response.data.template_id)
      if (template) {
        selectedTemplate.value = template
        // 如果是随机策略，确保 resources 是数组格式
        if (template.type === 'random') {
          // 确保 resources 字段存在且是数组
          if (!policyForm.parameters.resources || !Array.isArray(policyForm.parameters.resources)) {
            policyForm.parameters.resources = []
          } else {
            // 创建新数组以确保响应式
            policyForm.parameters.resources = [...policyForm.parameters.resources]
          }
          console.log('Random policy resources initialized:', policyForm.parameters.resources)
        }
        newRandomResourceId.value = ''
      }

      dialogVisible.value = true
    }
  } catch (error) {
    console.error('获取策略详情失败:', error)
    ElMessage.error(t('policies.getDetailFailed'))
  }
}

// 处理保存策略
const handleSavePolicy = async () => {
  try {
    await policyFormRef.value.validate()

    if (isAddMode.value) {
      // 创建策略
      await api.createPolicy({
        name: policyForm.name,
        template_id: policyForm.template_id,
        parameters: policyForm.parameters
      })
      ElMessage.success(t('policies.createSuccess'))
    } else {
      // 更新策略
      await api.updatePolicy(policyForm.id, {
        name: policyForm.name,
        parameters: policyForm.parameters,
        enabled: policyForm.enabled
      })
      ElMessage.success(t('policies.updateSuccess'))
    }

    dialogVisible.value = false
    getPolicyList()
  } catch (error) {
    console.error('保存策略失败:', error)
    ElMessage.error(t('policies.saveFailed'))
  }
}

// 处理删除策略
const handleDeletePolicy = async id => {
  try {
    // 检查是否为内置策略
    const policy = policies.value.find(p => p.id === id)
    if (policy && policy.builtin) {
      ElMessage.warning(t('policies.builtinCannotDelete'))
      return
    }

    await ElMessageBox.confirm(t('policies.deleteConfirmMessage'), t('policies.deleteConfirm'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'danger'
    })

    await api.deletePolicy(id)
    ElMessage.success(t('policies.deleteSuccess'))
    getPolicyList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除策略失败:', error)
      const errorMsg = error.response?.data?.error || t('policies.deleteFailed')
      ElMessage.error(errorMsg)
    }
  }
}

// 加权策略资源管理
const addWeightedResource = () => {
  if (!policyForm.parameters.resources) {
    policyForm.parameters.resources = []
  }
  policyForm.parameters.resources.push({ id: '', weight: 1 })
}

const removeWeightedResource = index => {
  policyForm.parameters.resources.splice(index, 1)
}

// 模型匹配规则管理
const addModelMapping = () => {
  if (!policyForm.parameters.mappings) {
    policyForm.parameters.mappings = []
  }
  policyForm.parameters.mappings.push({ model_pattern: '', resource_id: '' })
}

const removeModelMapping = index => {
  policyForm.parameters.mappings.splice(index, 1)
}

// 正则匹配规则管理
const addRegexRule = () => {
  if (!policyForm.parameters.rules) {
    policyForm.parameters.rules = []
  }
  policyForm.parameters.rules.push({ pattern: '', resource_id: '' })
}

const removeRegexRule = index => {
  policyForm.parameters.rules.splice(index, 1)
}

// 正则模型匹配模式管理

// 优先级策略资源管理
const addPriorityResource = () => {
  if (!policyForm.parameters.resources) {
    policyForm.parameters.resources = []
  }
  policyForm.parameters.resources.push({ id: '', priority: 1 })
}

const removePriorityResource = index => {
  policyForm.parameters.resources.splice(index, 1)
}

// 格式化日期
const formatDate = dateString => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 组件挂载时获取数据
onMounted(() => {
  getPolicyList()
  getTemplateList()
  getLLMResourceList()
})
</script>

<style scoped>
.policies-page {
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

.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.table-section {
  margin-top: 20px;
}

.policy-params {
  border: 1px solid var(--claude-border-cream);
  border-radius: var(--radius-comfortable);
  padding: 16px;
  background-color: #f8f7f3;
}

.resource-pool-info {
  margin-bottom: 10px;
}

.resource-pool-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 10px;
  background-color: var(--claude-ivory);
  border-radius: var(--radius-comfortable);
  min-height: 40px;
}

.resource-pool-item {
  display: inline-block;
}
</style>
