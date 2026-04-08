<template>
  <div class="models-page">
    <!-- Page Header -->
    <div class="page-header">
      <h1 class="page-title">{{ $t('models.title') }}</h1>
      <el-button type="primary" @click="handleAddModel">
        <el-icon><Plus /></el-icon>
        {{ $t('models.addModel') }}
      </el-button>
    </div>

    <!-- Search and Filter - Claude style -->
    <el-card class="filter-card">
      <div class="search-filter">
        <el-input
          v-model="searchQuery"
          :placeholder="$t('models.searchModelName')"
          prefix-icon="Search"
          style="width: 240px"
        ></el-input>
        <el-select v-model="typeFilter" :placeholder="$t('models.filterType')" style="width: 140px">
          <el-option :label="$t('models.all')" value=""></el-option>
          <el-option :label="$t('models.chat')" value="chat"></el-option>
          <el-option :label="$t('models.completion')" value="completion"></el-option>
          <el-option :label="$t('models.embedding')" value="embedding"></el-option>
          <el-option :label="$t('models.image')" value="image"></el-option>
        </el-select>
        <el-select
          v-model="statusFilter"
          :placeholder="$t('models.filterStatus')"
          style="width: 140px"
        >
          <el-option :label="$t('models.all')" value=""></el-option>
          <el-option :label="$t('models.active')" value="active"></el-option>
          <el-option :label="$t('models.inactive')" value="inactive"></el-option>
        </el-select>
      </div>
    </el-card>

    <!-- Models Table - Claude style -->
    <el-card class="table-card">
      <el-table v-loading="loading" :data="filteredModels" style="width: 100%">
        <el-table-column prop="id" :label="$t('models.id')" width="180" />
        <el-table-column prop="name" :label="$t('models.name')" />
        <el-table-column prop="model_id" :label="$t('models.modelId')" width="150">
          <template #default="scope">
            <el-tag type="info" size="small">{{ scope.row.model_id }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="$t('models.type')" width="100">
          <template #default="scope">
            <el-tag size="small">{{ scope.row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="category" :label="$t('models.category')" width="100">
          <template #default="scope">
            <el-tag type="success" size="small" v-if="scope.row.category">
              {{ scope.row.category }}
            </el-tag>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="llm_resource_id" :label="$t('models.llmResource')" width="180">
          <template #default="scope">
            {{ getLLMResourceName(scope.row.llm_resource_id) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('models.status')" width="100">
          <template #default="scope">
            <el-tag
              :type="
                scope.row.status === 'active'
                  ? 'success'
                  : scope.row.status === 'deprecated'
                    ? 'warning'
                    : 'danger'
              "
              size="small"
            >
              {{
                scope.row.status === 'active'
                  ? $t('models.active')
                  : scope.row.status === 'deprecated'
                    ? $t('models.deprecated')
                    : $t('models.inactive')
              }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('models.createdAt')" width="180">
          <template #default="scope">
            <span class="text-caption">{{ formatDate(scope.row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('models.actions')" width="240">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEditModel(scope.row)">
              {{ $t('models.edit') }}
            </el-button>
            <el-button type="default" size="small" @click="handleViewPricing(scope.row.id)">
              {{ $t('models.pricing') }}
            </el-button>
            <el-button type="danger" size="small" @click="handleDeleteModel(scope.row.id)">
              {{ $t('models.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- Add/Edit Dialog - Claude style -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="modelFormRef" :model="modelForm" :rules="modelRules" label-width="100px">
        <el-form-item :label="$t('models.modelName')" prop="name">
          <el-input v-model="modelForm.name" :placeholder="$t('models.modelNamePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('models.modelId')" prop="model_id">
          <el-input v-model="modelForm.model_id" :placeholder="$t('models.modelIdPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('models.modelType')" prop="type">
          <el-select v-model="modelForm.type" :placeholder="$t('models.selectModelType')">
            <el-option :label="$t('models.chat')" value="chat"></el-option>
            <el-option :label="$t('models.completion')" value="completion"></el-option>
            <el-option :label="$t('models.embedding')" value="embedding"></el-option>
            <el-option :label="$t('models.image')" value="image"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('models.modelCategory')" prop="category">
          <el-select v-model="modelForm.category" :placeholder="$t('models.selectModelCategory')">
            <el-option :label="$t('models.gpt')" value="gpt"></el-option>
            <el-option :label="$t('models.claude')" value="claude"></el-option>
            <el-option :label="$t('models.gemini')" value="gemini"></el-option>
            <el-option :label="$t('models.llama')" value="llama"></el-option>
            <el-option :label="$t('models.custom')" value="custom"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('models.llmResource')" prop="llm_resource_id">
          <el-select
            v-model="modelForm.llm_resource_id"
            :placeholder="$t('models.selectLLMResource')"
          >
            <el-option
              v-for="resource in llmResources"
              :key="resource.id"
              :label="resource.name"
              :value="resource.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('models.status')" prop="status">
          <el-select v-model="modelForm.status" :placeholder="$t('models.selectStatus')">
            <el-option :label="$t('models.active')" value="active"></el-option>
            <el-option :label="$t('models.inactive')" value="inactive"></el-option>
            <el-option :label="$t('models.deprecated')" value="deprecated"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('models.description')" prop="description">
          <el-input
            v-model="modelForm.description"
            :placeholder="$t('models.descriptionPlaceholder')"
            type="textarea"
            rows="3"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('models.cancel') }}</el-button>
          <el-button type="primary" @click="handleSaveModel">{{ $t('models.confirm') }}</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- Pricing Dialog - Claude style -->
    <el-dialog v-model="pricingVisible" :title="$t('models.pricingInfo')" width="500px">
      <el-card v-if="modelPricing" class="pricing-card">
        <el-descriptions :column="1" border>
          <el-descriptions-item
            :label="$t('models.inputPrice')"
            v-if="modelPricing.input_token_price"
          >
            <span class="price-value">${{ modelPricing.input_token_price }} / 1K tokens</span>
          </el-descriptions-item>
          <el-descriptions-item
            :label="$t('models.outputPrice')"
            v-if="modelPricing.output_token_price"
          >
            <span class="price-value">${{ modelPricing.output_token_price }} / 1K tokens</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('models.imagePrice')" v-if="modelPricing.image_price">
            <span class="price-value">${{ modelPricing.image_price }} / image</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('models.audioPrice')" v-if="modelPricing.audio_price">
            <span class="price-value">${{ modelPricing.audio_price }} / minute</span>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('models.currency')" v-if="modelPricing.currency">
            {{ modelPricing.currency }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
      <el-empty v-else :description="$t('models.noPricingInfo')" />
      <template #footer>
        <el-button @click="pricingVisible = false">{{ $t('models.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, PriceTag } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import api from '../api'

const { t } = useI18n()

const loading = ref(false)
const models = ref([])
const llmResources = ref([])
const searchQuery = ref('')
const typeFilter = ref('')
const statusFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const dialogVisible = ref(false)
const dialogTitle = ref('')
const isAddMode = ref(false)
const modelFormRef = ref(null)
const modelForm = reactive({
  id: '',
  name: '',
  model_id: '',
  type: 'chat',
  category: '',
  llm_resource_id: '',
  status: 'active',
  description: '',
  pricing: '',
  limits: '',
  parameters: '',
  features: ''
})

const pricingVisible = ref(false)
const modelPricing = ref(null)

const modelRules = computed(() => ({
  name: [{ required: true, message: t('models.nameRequired'), trigger: 'blur' }],
  model_id: [{ required: true, message: t('models.modelIdRequired'), trigger: 'blur' }],
  type: [{ required: true, message: t('models.typeRequired'), trigger: 'change' }],
  llm_resource_id: [
    { required: true, message: t('models.llmResourceRequired'), trigger: 'change' }
  ],
  status: [{ required: true, message: t('models.statusRequired'), trigger: 'change' }]
}))

const getModelList = async () => {
  try {
    loading.value = true
    const response = await api.getModels()
    if (response && response.data) {
      models.value = Array.isArray(response.data) ? response.data : []
      total.value = models.value.length
    }
  } catch (error) {
    console.error('获取模型列表失败:', error)
    ElMessage.error(t('models.getListFailed'))
  } finally {
    loading.value = false
  }
}

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

const filteredModels = computed(() => {
  let result = models.value

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(model => model.name.toLowerCase().includes(query))
  }

  if (typeFilter.value) {
    result = result.filter(model => model.type === typeFilter.value)
  }

  if (statusFilter.value) {
    result = result.filter(model => model.status === statusFilter.value)
  }

  const startIndex = (currentPage.value - 1) * pageSize.value
  const endIndex = startIndex + pageSize.value
  return result.slice(startIndex, endIndex)
})

const getLLMResourceName = resourceId => {
  const resource = llmResources.value.find(r => r.id === resourceId)
  return resource ? resource.name : t('models.unknown')
}

const handleAddModel = () => {
  isAddMode.value = true
  dialogTitle.value = t('models.addModel')
  Object.assign(modelForm, {
    id: '',
    name: '',
    model_id: '',
    type: 'chat',
    category: '',
    llm_resource_id: '',
    status: 'active',
    description: '',
    pricing: '',
    limits: '',
    parameters: '',
    features: ''
  })
  dialogVisible.value = true
}

const handleEditModel = model => {
  isAddMode.value = false
  dialogTitle.value = t('models.editModel')
  Object.assign(modelForm, model)
  dialogVisible.value = true
}

const handleSaveModel = async () => {
  try {
    await modelFormRef.value.validate()

    if (isAddMode.value) {
      await api.createModel(modelForm)
      ElMessage.success(t('models.createSuccess'))
    } else {
      await api.updateModel(modelForm.id, modelForm)
      ElMessage.success(t('models.updateSuccess'))
    }

    dialogVisible.value = false
    getModelList()
  } catch (error) {
    console.error('保存模型失败:', error)
    ElMessage.error(t('models.saveFailed'))
  }
}

const handleDeleteModel = async id => {
  try {
    await ElMessageBox.confirm(t('models.deleteConfirmMessage'), t('models.deleteConfirm'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'danger'
    })

    await api.deleteModel(id)
    ElMessage.success(t('models.deleteSuccess'))
    getModelList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除模型失败:', error)
      ElMessage.error(t('models.deleteFailed'))
    }
  }
}

const handleViewPricing = async id => {
  try {
    const response = await api.getModelPricing(id)
    if (response && response.data) {
      if (typeof response.data === 'string') {
        try {
          modelPricing.value = JSON.parse(response.data)
        } catch (e) {
          modelPricing.value = { raw: response.data }
        }
      } else {
        modelPricing.value = response.data
      }
    } else {
      modelPricing.value = null
    }
    pricingVisible.value = true
  } catch (error) {
    console.error('获取模型价格失败:', error)
    ElMessage.error(t('models.getPricingFailed'))
  }
}

const handleSizeChange = size => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = current => {
  currentPage.value = current
}

const formatDate = dateString => {
  const date = new Date(dateString)
  return date.toLocaleString()
}

onMounted(() => {
  getModelList()
  getLLMResourceList()
})
</script>

<style scoped>
.models-page {
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

/* Claude Style Page Header */
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

/* Claude Style Filter Card */
.filter-card {
  margin-bottom: 16px;
}

.search-filter {
  display: flex;
  gap: 12px;
  align-items: center;
}

/* Claude Style Table Card */
.table-card {
  margin-bottom: 24px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* Claude Style Pricing */
.pricing-card {
  border: none !important;
}

.price-value {
  font-family: var(--font-serif);
  font-size: 16px;
  font-weight: 500;
  color: var(--claude-terracotta);
}

/* Responsive */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }

  .search-filter {
    flex-direction: column;
    width: 100%;
  }

  .search-filter .el-input,
  .search-filter .el-select {
    width: 100% !important;
  }
}
</style>
