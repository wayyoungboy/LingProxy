<template>
  <div class="models">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('models.title') }}</span>
          <el-button type="primary" @click="handleAddModel">
            <el-icon><Plus /></el-icon>
            {{ $t('models.addModel') }}
          </el-button>
        </div>
      </template>
      
      <!-- 搜索和筛选 -->
      <div class="search-filter">
        <el-input
          v-model="searchQuery"
          :placeholder="$t('models.searchModelName')"
          prefix-icon="Search"
          style="width: 240px; margin-right: 10px"
        ></el-input>
        <el-select
          v-model="typeFilter"
          :placeholder="$t('models.filterType')"
          style="width: 120px; margin-right: 10px"
        >
          <el-option :label="$t('models.all')" value=""></el-option>
          <el-option :label="$t('models.chat')" value="chat"></el-option>
          <el-option :label="$t('models.completion')" value="completion"></el-option>
          <el-option :label="$t('models.embedding')" value="embedding"></el-option>
          <el-option :label="$t('models.image')" value="image"></el-option>
        </el-select>
        <el-select
          v-model="statusFilter"
          :placeholder="$t('models.filterStatus')"
          style="width: 120px"
        >
          <el-option :label="$t('models.all')" value=""></el-option>
          <el-option :label="$t('models.active')" value="active"></el-option>
          <el-option :label="$t('models.inactive')" value="inactive"></el-option>
        </el-select>
      </div>
      
      <!-- 模型列表 -->
      <el-table
        v-loading="loading"
        :data="filteredModels"
        style="width: 100%; margin-top: 20px"
        border
      >
        <el-table-column prop="id" :label="$t('models.id')" width="180" />
        <el-table-column prop="name" :label="$t('models.name')" />
        <el-table-column prop="model_id" :label="$t('models.modelId')" width="150">
          <template #default="scope">
            <el-tag type="info">{{ scope.row.model_id }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="$t('models.type')" width="100">
          <template #default="scope">
            <el-tag>{{ scope.row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="category" :label="$t('models.category')" width="100">
          <template #default="scope">
            <el-tag type="success" v-if="scope.row.category">{{ scope.row.category }}</el-tag>
            <span v-else>-</span>
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
              :type="scope.row.status === 'active' ? 'success' : scope.row.status === 'deprecated' ? 'warning' : 'danger'"
            >
              {{ scope.row.status === 'active' ? $t('models.active') : scope.row.status === 'deprecated' ? $t('models.deprecated') : $t('models.inactive') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('models.createdAt')" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('models.actions')" width="240">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              @click="handleEditModel(scope.row)"
              style="margin-right: 5px"
            >
              {{ $t('models.edit') }}
            </el-button>
            <el-button
              type="info"
              size="small"
              @click="handleViewPricing(scope.row.id)"
              style="margin-right: 5px"
            >
              {{ $t('models.pricing') }}
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDeleteModel(scope.row.id)"
            >
              {{ $t('models.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
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
    
    <!-- 添加/编辑模型对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
    >
      <el-form
        ref="modelFormRef"
        :model="modelForm"
        :rules="modelRules"
        label-width="100px"
      >
        <el-form-item :label="$t('models.modelName')" prop="name">
          <el-input v-model="modelForm.name" :placeholder="$t('models.modelNamePlaceholder')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('models.modelId')" prop="model_id">
          <el-input v-model="modelForm.model_id" :placeholder="$t('models.modelIdPlaceholder')"></el-input>
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
          <el-select v-model="modelForm.llm_resource_id" :placeholder="$t('models.selectLLMResource')">
            <el-option
              v-for="resource in llmResources"
              :key="resource.id"
              :label="resource.name"
              :value="resource.id"
            ></el-option>
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
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('models.cancel') }}</el-button>
          <el-button type="primary" @click="handleSaveModel">{{ $t('models.confirm') }}</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 模型价格对话框 -->
    <el-dialog
      v-model="pricingVisible"
      :title="$t('models.pricingInfo')"
      width="500px"
    >
      <el-card v-if="modelPricing" class="pricing-card">
        <el-descriptions :column="1" border>
          <el-descriptions-item :label="$t('models.inputPrice')" v-if="modelPricing.input_token_price">
            ${{ modelPricing.input_token_price }} / 1K tokens
          </el-descriptions-item>
          <el-descriptions-item :label="$t('models.outputPrice')" v-if="modelPricing.output_token_price">
            ${{ modelPricing.output_token_price }} / 1K tokens
          </el-descriptions-item>
          <el-descriptions-item :label="$t('models.imagePrice')" v-if="modelPricing.image_price">
            ${{ modelPricing.image_price }} / image
          </el-descriptions-item>
          <el-descriptions-item :label="$t('models.audioPrice')" v-if="modelPricing.audio_price">
            ${{ modelPricing.audio_price }} / minute
          </el-descriptions-item>
          <el-descriptions-item :label="$t('models.currency')" v-if="modelPricing.currency">
            {{ modelPricing.currency }}
          </el-descriptions-item>
          <el-descriptions-item :label="$t('models.rawData')" v-if="modelPricing.raw">
            <pre>{{ modelPricing.raw }}</pre>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
      <el-empty v-else :description="$t('models.noPricingInfo')"></el-empty>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="pricingVisible = false">{{ $t('models.close') }}</el-button>
        </span>
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

// 对话框相关
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

// 价格对话框相关
const pricingVisible = ref(false)
const modelPricing = ref(null)

// 表单验证规则
const modelRules = computed(() => ({
  name: [
    { required: true, message: t('models.nameRequired'), trigger: 'blur' }
  ],
  model_id: [
    { required: true, message: t('models.modelIdRequired'), trigger: 'blur' }
  ],
  type: [
    { required: true, message: t('models.typeRequired'), trigger: 'change' }
  ],
  llm_resource_id: [
    { required: true, message: t('models.llmResourceRequired'), trigger: 'change' }
  ],
  status: [
    { required: true, message: t('models.statusRequired'), trigger: 'change' }
  ]
}))

// 获取模型列表
const getModelList = async () => {
  try {
    loading.value = true
    const response = await api.getModels()
    if (response && response.data) {
      // 后端返回的是数组，不是分页对象
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

// 过滤模型
const filteredModels = computed(() => {
  let result = models.value
  
  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(model => 
      model.name.toLowerCase().includes(query)
    )
  }
  
  // 类型过滤
  if (typeFilter.value) {
    result = result.filter(model => model.type === typeFilter.value)
  }
  
  // 状态过滤
  if (statusFilter.value) {
    result = result.filter(model => model.status === statusFilter.value)
  }
  
  // 分页
  const startIndex = (currentPage.value - 1) * pageSize.value
  const endIndex = startIndex + pageSize.value
  return result.slice(startIndex, endIndex)
})

// 获取LLM资源名称
const getLLMResourceName = (resourceId) => {
  const resource = llmResources.value.find(r => r.id === resourceId)
  return resource ? resource.name : t('models.unknown')
}

// 处理添加模型
const handleAddModel = () => {
  isAddMode.value = true
  dialogTitle.value = t('models.addModel')
  // 重置表单
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

// 处理编辑模型
const handleEditModel = (model) => {
  isAddMode.value = false
  dialogTitle.value = t('models.editModel')
  // 填充表单
  Object.assign(modelForm, model)
  dialogVisible.value = true
}

// 处理保存模型
const handleSaveModel = async () => {
  try {
    // 验证表单
    await modelFormRef.value.validate()
    
    if (isAddMode.value) {
      // 创建模型
      await api.createModel(modelForm)
      ElMessage.success(t('models.createSuccess'))
    } else {
      // 更新模型
      await api.updateModel(modelForm.id, modelForm)
      ElMessage.success(t('models.updateSuccess'))
    }
    
    // 关闭对话框
    dialogVisible.value = false
    // 重新获取模型列表
    getModelList()
  } catch (error) {
    console.error('保存模型失败:', error)
    ElMessage.error(t('models.saveFailed'))
  }
}

// 处理删除模型
const handleDeleteModel = async (id) => {
  try {
    await ElMessageBox.confirm(
      t('models.deleteConfirmMessage'),
      t('models.deleteConfirm'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'danger'
      }
    )
    
    await api.deleteModel(id)
    ElMessage.success(t('models.deleteSuccess'))
    // 重新获取模型列表
    getModelList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除模型失败:', error)
      ElMessage.error(t('models.deleteFailed'))
    }
  }
}

// 处理查看价格
const handleViewPricing = async (id) => {
  try {
    const response = await api.getModelPricing(id)
    if (response && response.data) {
      // 如果返回的是对象，直接使用；如果是字符串，尝试解析
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

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = (current) => {
  currentPage.value = current
}

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString()
}

// 组件挂载时获取数据
onMounted(() => {
  getModelList()
  getLLMResourceList()
})
</script>

<style scoped>
.models {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-filter {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.dialog-footer {
  width: 100%;
  display: flex;
  justify-content: flex-end;
}

.pricing-card {
  margin-bottom: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .search-filter {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .el-input,
  .el-select {
    width: 100% !important;
  }
}
</style>