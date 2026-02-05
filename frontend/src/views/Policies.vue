<template>
  <div class="policies-container">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="page-title">策略管理</span>
          <el-button type="primary" @click="handleAddPolicy">
            <el-icon><Plus /></el-icon>
            创建策略
          </el-button>
        </div>
      </template>
      
      <!-- 策略列表 -->
      <el-table
        v-loading="loading"
        :data="policies"
        style="width: 100%"
        border
        stripe
      >
        <el-table-column prop="name" label="策略名称" width="200">
          <template #default="scope">
            <span>{{ scope.row.name }}</span>
            <el-tag v-if="scope.row.builtin" type="info" size="small" style="margin-left: 8px">内置</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="策略类型" width="150">
          <template #default="scope">
            <el-tag>{{ getTypeLabel(scope.row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="enabled" label="状态" width="100">
          <template #default="scope">
            <el-tag
              :type="scope.row.enabled ? 'success' : 'danger'"
            >
              {{ scope.row.enabled ? '启用' : '禁用' }}
            </el-tag>
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
              @click="handleEditPolicy(scope.row)"
              style="margin-right: 5px"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDeletePolicy(scope.row.id)"
              :disabled="scope.row.builtin"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 创建/编辑策略对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="700px"
    >
      <el-form
        ref="policyFormRef"
        :model="policyForm"
        :rules="policyRules"
        label-width="120px"
      >
        <el-form-item label="策略名称" prop="name">
          <el-input v-model="policyForm.name" placeholder="请输入策略名称"></el-input>
        </el-form-item>
        <el-form-item label="策略模板" prop="template_id">
          <el-select
            v-model="policyForm.template_id"
            placeholder="请选择策略模板"
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
        <el-form-item v-if="selectedTemplate" label="策略类型">
          <el-tag>{{ getTypeLabel(selectedTemplate.type) }}</el-tag>
        </el-form-item>
        <el-form-item v-if="selectedTemplate" label="参数配置">
          <div v-if="selectedTemplate.type === 'random'" class="policy-params">
            <el-form-item label="资源列表（可选）">
              <el-select
                v-model="policyForm.parameters.resources"
                multiple
                placeholder="选择资源（留空则使用所有可用资源）"
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
            <el-form-item label="只选择活跃资源">
              <el-switch v-model="policyForm.parameters.filter_by_status" />
            </el-form-item>
          </div>
          
          <div v-else-if="selectedTemplate.type === 'round_robin'" class="policy-params">
            <el-form-item label="资源列表" required>
              <el-select
                v-model="policyForm.parameters.resources"
                multiple
                placeholder="选择资源（按顺序）"
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
            <el-form-item label="只选择活跃资源">
              <el-switch v-model="policyForm.parameters.filter_by_status" />
            </el-form-item>
          </div>
          
          <div v-else-if="selectedTemplate.type === 'weighted'" class="policy-params">
            <el-form-item label="资源权重配置" required>
              <div
                v-for="(item, index) in policyForm.parameters.resources"
                :key="index"
                style="margin-bottom: 10px"
              >
                <el-select
                  v-model="item.id"
                  placeholder="选择资源"
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
                  placeholder="权重"
                  style="width: 120px"
                />
                <el-button
                  type="danger"
                  size="small"
                  @click="removeWeightedResource(index)"
                  style="margin-left: 10px"
                >
                  删除
                </el-button>
              </div>
              <el-button type="primary" size="small" @click="addWeightedResource">
                添加资源
              </el-button>
            </el-form-item>
          </div>
          
          <div v-else-if="selectedTemplate.type === 'model_match'" class="policy-params">
            <el-form-item label="模型匹配规则" required>
              <div
                v-for="(mapping, index) in policyForm.parameters.mappings"
                :key="index"
                style="margin-bottom: 10px"
              >
                <el-input
                  v-model="mapping.model_pattern"
                  placeholder="模型名模式（如 gpt-*）"
                  style="width: 200px; margin-right: 10px"
                />
                <el-select
                  v-model="mapping.resource_id"
                  placeholder="选择资源"
                  style="width: 200px; margin-right: 10px"
                >
                  <el-option
                    v-for="resource in llmResources"
                    :key="resource.id"
                    :label="resource.name"
                    :value="resource.id"
                  />
                </el-select>
                <el-button
                  type="danger"
                  size="small"
                  @click="removeModelMapping(index)"
                >
                  删除
                </el-button>
              </div>
              <el-button type="primary" size="small" @click="addModelMapping">
                添加规则
              </el-button>
            </el-form-item>
            <el-form-item label="默认资源">
              <el-select
                v-model="policyForm.parameters.default_resource_id"
                placeholder="选择默认资源（可选）"
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
            <el-form-item label="正则匹配规则" required>
              <div
                v-for="(rule, index) in policyForm.parameters.rules"
                :key="index"
                style="margin-bottom: 10px"
              >
                <el-input
                  v-model="rule.pattern"
                  placeholder="正则表达式（如 ^gpt-）"
                  style="width: 200px; margin-right: 10px"
                />
                <el-select
                  v-model="rule.resource_id"
                  placeholder="选择资源"
                  style="width: 200px; margin-right: 10px"
                >
                  <el-option
                    v-for="resource in llmResources"
                    :key="resource.id"
                    :label="resource.name"
                    :value="resource.id"
                  />
                </el-select>
                <el-button
                  type="danger"
                  size="small"
                  @click="removeRegexRule(index)"
                >
                  删除
                </el-button>
              </div>
              <el-button type="primary" size="small" @click="addRegexRule">
                添加规则
              </el-button>
            </el-form-item>
            <el-form-item label="默认资源">
              <el-select
                v-model="policyForm.parameters.default_resource_id"
                placeholder="选择默认资源（可选）"
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
          
          <div v-else-if="selectedTemplate.type === 'priority'" class="policy-params">
            <el-form-item label="资源优先级配置" required>
              <div
                v-for="(item, index) in policyForm.parameters.resources"
                :key="index"
                style="margin-bottom: 10px"
              >
                <el-select
                  v-model="item.id"
                  placeholder="选择资源"
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
                  placeholder="优先级（数字越小优先级越高）"
                  style="width: 200px"
                />
                <el-button
                  type="danger"
                  size="small"
                  @click="removePriorityResource(index)"
                  style="margin-left: 10px"
                >
                  删除
                </el-button>
              </div>
              <el-button type="primary" size="small" @click="addPriorityResource">
                添加资源
              </el-button>
            </el-form-item>
            <el-form-item label="启用降级">
              <el-switch v-model="policyForm.parameters.fallback_enabled" />
            </el-form-item>
          </div>
          
          <div v-else-if="selectedTemplate.type === 'failover'" class="policy-params">
            <el-form-item label="主资源" required>
              <el-select
                v-model="policyForm.parameters.primary_resource_id"
                placeholder="选择主资源"
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
            <el-form-item label="备用资源" required>
              <el-select
                v-model="policyForm.parameters.fallback_resources"
                multiple
                placeholder="选择备用资源（按顺序）"
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
            <el-form-item label="启用健康检查">
              <el-switch v-model="policyForm.parameters.health_check_enabled" />
            </el-form-item>
            <el-form-item label="健康检查间隔（秒）">
              <el-input-number
                v-model="policyForm.parameters.health_check_interval"
                :min="10"
                :max="300"
              />
            </el-form-item>
          </div>
        </el-form-item>
        <el-form-item v-if="!isAddMode" label="状态" prop="enabled">
          <el-switch v-model="policyForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSavePolicy">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import api from '../api'

const loading = ref(false)
const policies = ref([])
const templates = ref([])
const llmResources = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isAddMode = ref(false)
const policyFormRef = ref(null)
const selectedTemplate = ref(null)

const policyForm = reactive({
  id: '',
  name: '',
  template_id: '',
  parameters: {},
  enabled: true
})

const policyRules = {
  name: [
    { required: true, message: '请输入策略名称', trigger: 'blur' }
  ],
  template_id: [
    { required: true, message: '请选择策略模板', trigger: 'change' }
  ]
}

// 策略类型标签
const typeLabels = {
  random: '随机选择',
  round_robin: '轮询负载均衡',
  weighted: '加权负载均衡',
  model_match: '模型名匹配',
  regex_match: '正则匹配',
  priority: '优先级策略',
  failover: '故障转移'
}

const getTypeLabel = (type) => {
  return typeLabels[type] || type
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
    ElMessage.error('获取策略列表失败')
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
    if (response && response.data) {
      llmResources.value = response.data.items || []
    }
  } catch (error) {
    console.error('获取LLM资源列表失败:', error)
  }
}

// 处理模板变化
const handleTemplateChange = (templateId) => {
  const template = templates.value.find(t => t.id === templateId)
  if (template) {
    selectedTemplate.value = template
    // 初始化参数
    initParameters(template)
  }
}

// 初始化参数
const initParameters = (template) => {
  const defaultParams = template.default_parameters ? JSON.parse(template.default_parameters) : {}
  
  switch (template.type) {
    case 'random':
      policyForm.parameters = {
        resources: defaultParams.resources || [],
        filter_by_status: defaultParams.filter_by_status !== false
      }
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
const handleAddPolicy = () => {
  isAddMode.value = true
  dialogTitle.value = '创建策略'
  Object.assign(policyForm, {
    id: '',
    name: '',
    template_id: '',
    parameters: {},
    enabled: true
  })
  selectedTemplate.value = null
  dialogVisible.value = true
}

// 处理编辑策略
const handleEditPolicy = async (policy) => {
  isAddMode.value = false
  dialogTitle.value = '编辑策略'
  
  // 获取策略详情
  try {
    const response = await api.getPolicy(policy.id)
    if (response && response.data) {
      Object.assign(policyForm, {
        id: response.data.id,
        name: response.data.name,
        template_id: response.data.template_id,
        parameters: response.data.parameters || {},
        enabled: response.data.enabled
      })
      
      // 设置选中的模板
      const template = templates.value.find(t => t.id === response.data.template_id)
      if (template) {
        selectedTemplate.value = template
      }
      
      dialogVisible.value = true
    }
  } catch (error) {
    console.error('获取策略详情失败:', error)
    ElMessage.error('获取策略详情失败')
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
      ElMessage.success('策略创建成功')
    } else {
      // 更新策略
      await api.updatePolicy(policyForm.id, {
        name: policyForm.name,
        parameters: policyForm.parameters,
        enabled: policyForm.enabled
      })
      ElMessage.success('策略更新成功')
    }
    
    dialogVisible.value = false
    getPolicyList()
  } catch (error) {
    console.error('保存策略失败:', error)
    ElMessage.error('保存策略失败')
  }
}

// 处理删除策略
const handleDeletePolicy = async (id) => {
  try {
    // 检查是否为内置策略
    const policy = policies.value.find(p => p.id === id)
    if (policy && policy.builtin) {
      ElMessage.warning('内置策略不允许删除')
      return
    }

    await ElMessageBox.confirm(
      '确定要删除这个策略吗？删除后使用该策略的Token将使用默认策略。',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'danger'
      }
    )
    
    await api.deletePolicy(id)
    ElMessage.success('策略删除成功')
    getPolicyList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除策略失败:', error)
      const errorMsg = error.response?.data?.error || '删除策略失败'
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

const removeWeightedResource = (index) => {
  policyForm.parameters.resources.splice(index, 1)
}

// 模型匹配规则管理
const addModelMapping = () => {
  if (!policyForm.parameters.mappings) {
    policyForm.parameters.mappings = []
  }
  policyForm.parameters.mappings.push({ model_pattern: '', resource_id: '' })
}

const removeModelMapping = (index) => {
  policyForm.parameters.mappings.splice(index, 1)
}

// 正则匹配规则管理
const addRegexRule = () => {
  if (!policyForm.parameters.rules) {
    policyForm.parameters.rules = []
  }
  policyForm.parameters.rules.push({ pattern: '', resource_id: '' })
}

const removeRegexRule = (index) => {
  policyForm.parameters.rules.splice(index, 1)
}

// 优先级策略资源管理
const addPriorityResource = () => {
  if (!policyForm.parameters.resources) {
    policyForm.parameters.resources = []
  }
  policyForm.parameters.resources.push({ id: '', priority: 1 })
}

const removePriorityResource = (index) => {
  policyForm.parameters.resources.splice(index, 1)
}

// 格式化日期
const formatDate = (dateString) => {
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
.policies-container {
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

.policy-params {
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 15px;
  background-color: #f5f7fa;
}
</style>
