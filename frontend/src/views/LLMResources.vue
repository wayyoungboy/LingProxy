<template>
  <div class="llm-resources-page">
    <div class="page-header">
      <h1 class="page-title">{{ $t('llmResources.title') }}</h1>
      <div class="header-actions">
        <el-button type="default" @click="handleDownloadTemplate">
          <el-icon><Download /></el-icon>
          {{ $t('llmResources.downloadTemplate') }}
        </el-button>
        <el-button type="primary" @click="handleAddResource">
          <el-icon><Plus /></el-icon>
          {{ $t('llmResources.addResource') }}
        </el-button>
      </div>
    </div>

    <el-card class="filter-card">
      <!-- 搜索和筛选 -->
      <div class="search-filter">
        <el-input
          v-model="searchQuery"
          :placeholder="$t('llmResources.searchPlaceholder')"
          :prefix-icon="Search"
          style="width: 300px; margin-right: 10px"
        ></el-input>
        <el-select
          v-model="typeFilter"
          :placeholder="$t('llmResources.filterType')"
          style="width: 140px; margin-right: 10px"
        >
          <el-option :label="$t('common.all')" value=""></el-option>
          <el-option :label="$t('llmResources.typeChat')" value="chat"></el-option>
          <el-option :label="$t('llmResources.typeImage')" value="image"></el-option>
          <el-option :label="$t('llmResources.typeEmbedding')" value="embedding"></el-option>
          <el-option :label="$t('llmResources.typeRerank')" value="rerank"></el-option>
          <el-option :label="$t('llmResources.typeAudio')" value="audio"></el-option>
          <el-option :label="$t('llmResources.typeVideo')" value="video"></el-option>
        </el-select>
        <el-select
          v-model="statusFilter"
          :placeholder="$t('llmResources.filterStatus')"
          style="width: 120px"
        >
          <el-option :label="$t('common.all')" value=""></el-option>
          <el-option :label="$t('llmResources.active')" value="active"></el-option>
          <el-option :label="$t('llmResources.inactive')" value="inactive"></el-option>
        </el-select>
      </div>

      <!-- LLM资源列表 -->
      <el-table
        v-loading="loading"
        :data="filteredResources"
        style="width: 100%; margin-top: 20px"
        border
        stripe
      >
        <el-table-column prop="id" label="ID" width="180" />
        <el-table-column prop="name" :label="$t('llmResources.name')" />
        <el-table-column prop="type" :label="$t('llmResources.modelType')" width="120">
          <template #default="scope">
            <el-tag>{{ getTypeLabel(scope.row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="driver" :label="$t('llmResources.driver')" width="100">
          <template #default="scope">
            <el-tag type="info">{{ getDriverLabel(scope.row.driver) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="model" :label="$t('llmResources.modelId')" width="150">
          <template #default="scope">
            <el-tag type="warning" v-if="scope.row.model">{{ scope.row.model }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="base_url" :label="$t('llmResources.baseUrl')" />
        <el-table-column prop="status" :label="$t('common.status')">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'danger'">
              {{
                scope.row.status === 'active'
                  ? $t('llmResources.active')
                  : $t('llmResources.inactive')
              }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="test_status" :label="$t('llmResources.testStatus')" width="120">
          <template #default="scope">
            <el-tag
              :type="
                scope.row.test_status === 'passed'
                  ? 'success'
                  : scope.row.test_status === 'failed'
                    ? 'danger'
                    : 'info'
              "
            >
              {{
                scope.row.test_status === 'passed'
                  ? $t('llmResources.testPassed')
                  : scope.row.test_status === 'failed'
                    ? $t('llmResources.testFailed')
                    : $t('llmResources.testNotTested')
              }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="300">
          <template #default="scope">
            <el-button
              type="success"
              size="small"
              @click="handleTestResource(scope.row)"
              :loading="testingResources[scope.row.id]"
              style="margin-right: 5px"
            >
              {{ $t('llmResources.test') }}
            </el-button>
            <el-button
              type="primary"
              size="small"
              @click="handleEditResource(scope.row)"
              style="margin-right: 5px"
            >
              {{ $t('common.edit') }}
            </el-button>
            <el-button type="danger" size="small" @click="handleDeleteResource(scope.row.id)">
              {{ $t('common.delete') }}
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

    <!-- 添加/编辑LLM资源对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form
        ref="resourceFormRef"
        :model="resourceForm"
        :rules="resourceRules"
        label-width="100px"
      >
        <el-form-item :label="$t('llmResources.name')" prop="name">
          <el-input
            v-model="resourceForm.name"
            :placeholder="$t('llmResources.nameRequired')"
          ></el-input>
        </el-form-item>
        <el-form-item :label="$t('llmResources.modelType')" prop="type">
          <el-select
            v-model="resourceForm.type"
            :placeholder="$t('llmResources.selectType')"
            @change="handleTypeChange"
          >
            <el-option :label="$t('llmResources.typeChat')" value="chat"></el-option>
            <el-option :label="$t('llmResources.typeImage')" value="image"></el-option>
            <el-option :label="$t('llmResources.typeEmbedding')" value="embedding"></el-option>
            <el-option :label="$t('llmResources.typeRerank')" value="rerank"></el-option>
            <el-option :label="$t('llmResources.typeAudio')" value="audio"></el-option>
            <el-option :label="$t('llmResources.typeVideo')" value="video"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('llmResources.driver')" prop="driver">
          <el-select
            v-model="resourceForm.driver"
            :placeholder="$t('llmResources.selectDriver')"
            @change="handleDriverChange"
          >
            <el-option label="OpenAI" value="openai"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('llmResources.modelId')" prop="model">
          <el-input
            v-model="resourceForm.model"
            :placeholder="$t('llmResources.modelIdPlaceholder')"
          ></el-input>
        </el-form-item>
        <el-form-item :label="$t('llmResources.baseUrl')" prop="base_url">
          <el-input
            v-model="resourceForm.base_url"
            :placeholder="$t('llmResources.baseUrlPlaceholder')"
          ></el-input>
        </el-form-item>
        <el-form-item :label="$t('llmResources.apiKey')" prop="api_key">
          <el-input
            v-model="resourceForm.api_key"
            :placeholder="$t('llmResources.apiKeyRequired')"
            :type="apiKeyVisible ? 'text' : 'password'"
          >
            <template #append>
              <el-button @click="apiKeyVisible = !apiKeyVisible" type="text">
                {{ apiKeyVisible ? $t('common.hide') : $t('common.show') }}
              </el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('common.status')" prop="status">
          <el-select v-model="resourceForm.status" :placeholder="$t('llmResources.selectStatus')">
            <el-option :label="$t('llmResources.active')" value="active"></el-option>
            <el-option :label="$t('llmResources.inactive')" value="inactive"></el-option>
          </el-select>
        </el-form-item>

        <!-- Embedding 类型特定配置 -->
        <template v-if="resourceForm.type === 'embedding'">
          <el-divider content-position="left">
            {{ $t('llmResources.typeEmbedding') }} 配置
          </el-divider>
          <el-form-item :label="$t('llmResources.supportedDimensions')">
            <el-input
              v-model="resourceForm.embedding_config.supported_dimensions"
              :placeholder="$t('llmResources.supportedDimensionsPlaceholder')"
            >
              <template #append>
                <el-button @click="parseDimensions">解析</el-button>
              </template>
            </el-input>
            <div style="font-size: 12px; color: #909399; margin-top: 4px">例如：512,1024,1536</div>
          </el-form-item>
          <el-form-item :label="$t('llmResources.defaultDimension')">
            <el-input-number
              v-model="resourceForm.embedding_config.default_dimension"
              :placeholder="$t('llmResources.defaultDimensionPlaceholder')"
              :min="1"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item :label="$t('llmResources.normalize')">
            <el-switch v-model="resourceForm.embedding_config.normalize" />
          </el-form-item>
          <el-form-item :label="$t('llmResources.maxTokens')">
            <el-input-number
              v-model="resourceForm.embedding_config.max_input_tokens"
              :min="1"
              style="width: 100%"
            />
          </el-form-item>
        </template>

        <!-- Rerank 类型特定配置 -->
        <template v-if="resourceForm.type === 'rerank'">
          <el-divider content-position="left">{{ $t('llmResources.typeRerank') }} 配置</el-divider>
          <el-form-item :label="$t('llmResources.defaultTopN')">
            <el-input-number
              v-model="resourceForm.rerank_config.default_top_n"
              :placeholder="$t('llmResources.defaultTopNPlaceholder')"
              :min="1"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item :label="$t('llmResources.maxDocuments')">
            <el-input-number
              v-model="resourceForm.rerank_config.max_documents"
              :min="1"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item :label="$t('llmResources.maxQueryLength')">
            <el-input-number
              v-model="resourceForm.rerank_config.max_query_length"
              :min="1"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item :label="$t('llmResources.maxDocumentLength')">
            <el-input-number
              v-model="resourceForm.rerank_config.max_document_length"
              :min="1"
              style="width: 100%"
            />
          </el-form-item>
        </template>

        <!-- Chat 类型特定配置 -->
        <template v-if="resourceForm.type === 'chat'">
          <el-divider content-position="left">{{ $t('llmResources.typeChat') }} 配置</el-divider>
          <el-form-item :label="$t('llmResources.maxTokens')">
            <el-input-number
              v-model="resourceForm.chat_config.max_tokens"
              :min="1"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item :label="$t('llmResources.contextWindow')">
            <el-input-number
              v-model="resourceForm.chat_config.context_window"
              :min="1"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item :label="$t('llmResources.supportsStreaming')">
            <el-switch v-model="resourceForm.chat_config.supports_streaming" />
          </el-form-item>
          <el-form-item :label="$t('llmResources.supportsFunctionCalling')">
            <el-switch v-model="resourceForm.chat_config.supports_function_calling" />
          </el-form-item>
          <el-form-item :label="$t('llmResources.supportsVision')">
            <el-switch v-model="resourceForm.chat_config.supports_vision" />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
          <el-button type="primary" @click="handleSaveResource">
            {{ $t('common.confirm') }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- JSON 批量导入对话框 -->
    <el-dialog
      v-model="jsonImportDialogVisible"
      :title="$t('llmResources.jsonImportTitle')"
      width="720px"
    >
      <el-alert
        type="info"
        :closable="false"
        show-icon
        :title="$t('llmResources.jsonImportNotice')"
      />

      <div class="json-import-tabs">
        <el-radio-group v-model="jsonImportMode" size="small">
          <el-radio-button label="array">{{ $t('llmResources.arrayFormat') }}</el-radio-button>
          <el-radio-button label="codingPlan">
            {{ $t('llmResources.codingPlanFormat') }}
          </el-radio-button>
        </el-radio-group>
      </div>

      <div class="json-import-example">
        <div class="json-import-example-header">
          <span>{{ $t('llmResources.jsonExample') }}</span>
          <el-button type="primary" text size="small" @click="fillJsonExample">
            {{ $t('llmResources.fillExample') }}
          </el-button>
        </div>
        <pre class="json-example-block">{{ jsonImportExample }}</pre>
      </div>

      <el-form label-position="top">
        <el-form-item :label="$t('llmResources.jsonDataLabel')">
          <el-input
            v-model="jsonImportText"
            type="textarea"
            :rows="12"
            :placeholder="$t('llmResources.jsonDataPlaceholder')"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="jsonImportDialogVisible = false">{{ $t('common.cancel') }}</el-button>
          <el-button type="primary" :loading="jsonImportLoading" @click="handleJsonImport">
            {{ $t('llmResources.startImport') }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Upload, Download } from '@element-plus/icons-vue'
import api from '../api'
import { formatDate } from '../utils/index'

const { t } = useI18n()
const loading = ref(false)
const resources = ref([])
const searchQuery = ref('')
const typeFilter = ref('')
const statusFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const uploadRef = ref(null)

// 测试相关
const testingResources = ref({}) // 记录正在测试的资源ID

// JSON 导入相关
const jsonImportDialogVisible = ref(false)
const jsonImportText = ref('')
const jsonImportLoading = ref(false)
const jsonImportMode = ref('array') // 'array' 或 'codingPlan'

// 数组格式示例
const arrayFormatExample = `[
  {
    "name": "OpenAI GPT-4",
    "type": "chat",
    "driver": "openai",
    "model": "gpt-4",
    "base_url": "https://api.openai.com/v1",
    "api_key": "sk-xxxxxxxxxxxxx",
    "status": "active"
  },
  {
    "name": "OpenAI GPT-3.5",
    "type": "chat",
    "driver": "openai",
    "model": "gpt-3.5-turbo",
    "base_url": "https://api.openai.com/v1",
    "api_key": "sk-yyyyyyyyyyyyy",
    "status": "active"
  }
]`

// Coding Plan 格式示例
const codingPlanFormatExample = `{
  "base_url": "https://dashscope.aliyuncs.com/compatible-mode/v1",
  "api_key": "sk-xxxxxxxxxxxxx",
  "models": [
    {
      "name": "qwen-turbo",
      "type": "chat",
      "display_name": "通义千问 Turbo"
    },
    {
      "name": "qwen-plus",
      "types": ["chat", "embedding"],
      "display_name": "通义千问 Plus"
    }
  ]
}`

// 根据模式显示不同示例
const jsonImportExample = computed(() => {
  return jsonImportMode.value === 'codingPlan' ? codingPlanFormatExample : arrayFormatExample
})

// 上传配置
const uploadAction = '/api/v1/llm-resources/import'
const uploadHeaders = computed(() => {
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
})

// 对话框相关
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isAddMode = ref(false)
const resourceFormRef = ref(null)
const apiKeyVisible = ref(false)
const resourceForm = reactive({
  id: '',
  name: '',
  type: 'chat',
  driver: 'openai',
  model: '',
  base_url: '',
  api_key: '',
  status: 'active',
  // 类型特定配置
  embedding_config: {
    supported_dimensions: '', // 逗号分隔的字符串，如 "512,1024,1536"
    default_dimension: null,
    normalize: false,
    max_input_tokens: null
  },
  rerank_config: {
    default_top_n: null,
    max_documents: null,
    max_query_length: null,
    max_document_length: null
  },
  chat_config: {
    max_tokens: null,
    context_window: null,
    supports_streaming: true,
    supports_function_calling: false,
    supports_vision: false
  },
  image_config: {
    max_image_size: null,
    supported_formats: '',
    max_images_per_request: 1,
    quality: 'standard'
  },
  audio_config: {
    max_audio_duration: null,
    supported_formats: '',
    max_file_size: null,
    language: ''
  }
})

// 表单验证规则
const resourceRules = computed(() => ({
  name: [{ required: true, message: t('llmResources.nameRequired'), trigger: 'blur' }],
  type: [{ required: true, message: t('llmResources.typeRequired'), trigger: 'change' }],
  driver: [{ required: true, message: t('llmResources.driverRequired'), trigger: 'change' }],
  model: [{ required: true, message: t('llmResources.modelRequired'), trigger: 'blur' }],
  base_url: [{ required: true, message: t('llmResources.endpointRequired'), trigger: 'blur' }],
  api_key: [{ required: true, message: t('llmResources.apiKeyRequired'), trigger: 'blur' }],
  status: [{ required: true, message: t('llmResources.statusRequired'), trigger: 'change' }]
}))

// 获取LLM资源列表
const getResourceList = async () => {
  try {
    loading.value = true
    const response = await api.getLLMResources({
      page: currentPage.value,
      page_size: pageSize.value
    })
    if (response && response.data) {
      // 后端返回的data可能是数组，也可能是分页对象
      if (Array.isArray(response.data)) {
        // 直接是数组格式
        resources.value = response.data
        total.value = response.data.length
      } else if (response.data.items) {
        // 分页格式
        resources.value = response.data.items || []
        total.value = response.data.total || 0
      } else {
        // 其他格式，尝试直接使用data
        resources.value = []
        total.value = 0
      }
    }
  } catch (error) {
    console.error('获取LLM资源列表失败:', error)
    ElMessage.error(t('llmResources.getListFailed'))
  } finally {
    loading.value = false
  }
}

// 过滤资源
const filteredResources = computed(() => {
  let result = resources.value

  // 搜索过滤（支持资源名称、基础URL和模型标识的模糊搜索）
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase().trim()
    result = result.filter(resource => {
      // 搜索资源名称
      const nameMatch = resource.name && resource.name.toLowerCase().includes(query)
      // 搜索基础URL（模糊匹配）
      const urlMatch = resource.base_url && resource.base_url.toLowerCase().includes(query)
      // 搜索模型标识（可选，增强搜索体验）
      const modelMatch = resource.model && resource.model.toLowerCase().includes(query)
      // 只要有一个字段匹配就返回true
      return nameMatch || urlMatch || modelMatch
    })
  }

  // 类型过滤
  if (typeFilter.value) {
    result = result.filter(resource => resource.type === typeFilter.value)
  }

  // 状态过滤
  if (statusFilter.value) {
    result = result.filter(resource => resource.status === statusFilter.value)
  }

  // 分页
  const startIndex = (currentPage.value - 1) * pageSize.value
  const endIndex = startIndex + pageSize.value
  return result.slice(startIndex, endIndex)
})

// 获取模型类别标签
const getTypeLabel = type => {
  const labels = {
    chat: t('llmResources.typeChat'),
    image: t('llmResources.typeImage'),
    embedding: t('llmResources.typeEmbedding'),
    rerank: t('llmResources.typeRerank'),
    audio: t('llmResources.typeAudio'),
    video: t('llmResources.typeVideo')
  }
  return labels[type] || type
}

// 驱动标签映射
const driverLabels = {
  openai: 'OpenAI'
}

// 获取驱动标签
const getDriverLabel = driver => {
  return driverLabels[driver] || driver
}

// 驱动到BaseURL的映射
const driverToBaseURL = {
  openai: 'https://api.openai.com/v1'
}

// 处理驱动变化
const handleDriverChange = driver => {
  if (!driver) return

  // 根据驱动自动填充BaseURL
  const defaultURL = driverToBaseURL[driver]
  if (defaultURL) {
    // 只有在添加模式，或者当前base_url为空，或者是默认URL时才自动填充
    if (
      isAddMode.value ||
      !resourceForm.base_url ||
      Object.values(driverToBaseURL).includes(resourceForm.base_url)
    ) {
      resourceForm.base_url = defaultURL
    }
  }
}

// 处理类型变化
const handleTypeChange = type => {
  // 切换类型时重置类型特定配置
  resetTypeConfig()
}

// 重置类型特定配置
const resetTypeConfig = () => {
  resourceForm.embedding_config = {
    supported_dimensions: '',
    default_dimension: null,
    normalize: false,
    max_input_tokens: null
  }
  resourceForm.rerank_config = {
    default_top_n: null,
    max_documents: null,
    max_query_length: null,
    max_document_length: null
  }
  resourceForm.chat_config = {
    max_tokens: null,
    context_window: null,
    supports_streaming: true,
    supports_function_calling: false,
    supports_vision: false
  }
  resourceForm.image_config = {
    max_image_size: null,
    supported_formats: '',
    max_images_per_request: 1,
    quality: 'standard'
  }
  resourceForm.audio_config = {
    max_audio_duration: null,
    supported_formats: '',
    max_file_size: null,
    language: ''
  }
}

// 处理添加资源
const handleAddResource = () => {
  isAddMode.value = true
  dialogTitle.value = t('llmResources.addResource')
  // 重置表单
  Object.assign(resourceForm, {
    id: '',
    name: '',
    type: 'chat',
    driver: 'openai',
    model: '',
    base_url: driverToBaseURL['openai'], // 默认填充OpenAI的URL
    api_key: '',
    status: 'active'
  })
  resetTypeConfig()
  apiKeyVisible.value = false
  dialogVisible.value = true
}

// 处理编辑资源
const handleEditResource = resource => {
  isAddMode.value = false
  dialogTitle.value = t('llmResources.editResource')
  // 填充基础字段
  Object.assign(resourceForm, {
    id: resource.id,
    name: resource.name,
    type: resource.type,
    driver: resource.driver,
    model: resource.model,
    base_url: resource.base_url,
    api_key: resource.api_key,
    status: resource.status
  })

  // 填充类型特定配置
  resetTypeConfig()
  if (resource.type === 'embedding' && resource.embedding_config) {
    const config = resource.embedding_config
    resourceForm.embedding_config = {
      supported_dimensions: Array.isArray(config.supported_dimensions)
        ? config.supported_dimensions.join(',')
        : config.supported_dimensions || '',
      default_dimension: config.default_dimension || null,
      normalize: config.normalize || false,
      max_input_tokens: config.max_input_tokens || null
    }
  } else if (resource.type === 'rerank' && resource.rerank_config) {
    const config = resource.rerank_config
    resourceForm.rerank_config = {
      default_top_n: config.default_top_n || null,
      max_documents: config.max_documents || null,
      max_query_length: config.max_query_length || null,
      max_document_length: config.max_document_length || null
    }
  } else if (resource.type === 'chat' && resource.chat_config) {
    const config = resource.chat_config
    resourceForm.chat_config = {
      max_tokens: config.max_tokens || null,
      context_window: config.context_window || null,
      supports_streaming: config.supports_streaming !== false,
      supports_function_calling: config.supports_function_calling || false,
      supports_vision: config.supports_vision || false
    }
  }

  apiKeyVisible.value = false
  dialogVisible.value = true
}

// 解析维度字符串为数组
const parseDimensions = () => {
  if (!resourceForm.embedding_config.supported_dimensions) return
  const dims = resourceForm.embedding_config.supported_dimensions
    .split(',')
    .map(d => d.trim())
    .filter(d => d && !isNaN(d))
    .map(d => parseInt(d))
  resourceForm.embedding_config.supported_dimensions = dims.join(',')
}

// 准备提交数据
const prepareSubmitData = () => {
  const data = {
    id: resourceForm.id,
    name: resourceForm.name,
    type: resourceForm.type,
    driver: resourceForm.driver,
    model: resourceForm.model,
    base_url: resourceForm.base_url,
    api_key: resourceForm.api_key,
    status: resourceForm.status
  }

  // 根据类型添加特定配置
  if (resourceForm.type === 'embedding') {
    // 解析维度字符串为数组
    const dims = resourceForm.embedding_config.supported_dimensions
      ? resourceForm.embedding_config.supported_dimensions
          .split(',')
          .map(d => d.trim())
          .filter(d => d && !isNaN(d))
          .map(d => parseInt(d))
      : []

    data.embedding_config = {
      supported_dimensions: dims,
      default_dimension: resourceForm.embedding_config.default_dimension,
      normalize: resourceForm.embedding_config.normalize,
      max_input_tokens: resourceForm.embedding_config.max_input_tokens
    }
  } else if (resourceForm.type === 'rerank') {
    data.rerank_config = {
      default_top_n: resourceForm.rerank_config.default_top_n,
      max_documents: resourceForm.rerank_config.max_documents,
      max_query_length: resourceForm.rerank_config.max_query_length,
      max_document_length: resourceForm.rerank_config.max_document_length
    }
  } else if (resourceForm.type === 'chat') {
    data.chat_config = {
      max_tokens: resourceForm.chat_config.max_tokens,
      context_window: resourceForm.chat_config.context_window,
      supports_streaming: resourceForm.chat_config.supports_streaming,
      supports_function_calling: resourceForm.chat_config.supports_function_calling,
      supports_vision: resourceForm.chat_config.supports_vision
    }
  }

  return data
}

// 处理保存资源
const handleSaveResource = async () => {
  try {
    // 验证表单
    await resourceFormRef.value.validate()

    // 准备提交数据
    const submitData = prepareSubmitData()

    if (isAddMode.value) {
      // 创建资源
      await api.createLLMResource(submitData)
      ElMessage.success(t('llmResources.createSuccess'))
    } else {
      // 更新资源
      await api.updateLLMResource(resourceForm.id, submitData)
      ElMessage.success(t('llmResources.updateSuccess'))
    }

    // 关闭对话框
    dialogVisible.value = false
    // 重新获取资源列表
    getResourceList()
  } catch (error) {
    console.error('保存LLM资源失败:', error)
    ElMessage.error(t('llmResources.saveFailed'))
  }
}

// 处理删除资源
const handleDeleteResource = async id => {
  try {
    await ElMessageBox.confirm(
      t('llmResources.deleteConfirmMessage'),
      t('llmResources.deleteConfirm'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'danger'
      }
    )

    await api.deleteLLMResource(id)
    ElMessage.success(t('llmResources.deleteSuccess'))
    // 重新获取资源列表
    getResourceList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除LLM资源失败:', error)
      ElMessage.error(t('llmResources.deleteFailed'))
    }
  }
}

// 处理测试资源
const handleTestResource = async resource => {
  // 检查资源状态
  if (resource.status !== 'active') {
    ElMessage.warning(t('llmResources.testOnlyActive'))
    return
  }

  // 设置测试状态
  testingResources.value[resource.id] = true

  try {
    const result = await api.testLLMResource(resource.id)

    if (result && result.success) {
      // 测试成功，显示详细信息
      let message = t('llmResources.testSuccess') + '！\n'
      message += t('llmResources.duration') + `: ${result.duration_ms}ms\n`

      if (result.model) {
        message += t('llmResources.model') + `: ${result.model}\n`
      }

      if (result.response) {
        message += t('llmResources.response') + `: ${result.response}\n`
      }

      if (result.embedding_dimension) {
        message += t('llmResources.embeddingDimension') + `: ${result.embedding_dimension}\n`
      }

      if (result.usage) {
        message += t('llmResources.tokenUsage') + `: ${JSON.stringify(result.usage, null, 2)}`
      }

      ElMessageBox.alert(message, t('llmResources.testSuccess'), {
        confirmButtonText: t('common.confirm'),
        type: 'success',
        dangerouslyUseHTMLString: false
      })
      // 刷新资源列表以显示更新后的测试状态
      await getResourceList()
    } else {
      // 测试失败
      const errorMsg = result?.message || result?.error || '测试失败'
      ElMessage.error(errorMsg)
      // 即使失败也刷新资源列表以显示更新后的测试状态
      await getResourceList()
    }
  } catch (error) {
    console.error('测试LLM资源失败:', error)
    const errorMsg =
      error.response?.data?.error || error.response?.data?.message || error.message || '测试失败'
    ElMessage.error(errorMsg)
  } finally {
    // 清除测试状态
    testingResources.value[resource.id] = false
  }
}

// 查看资源下的模型

// 分页处理
const handleSizeChange = size => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = current => {
  currentPage.value = current
}

// 下载导入模板
const handleDownloadTemplate = async () => {
  try {
    const response = await api.downloadLLMResourcesTemplate()

    // 响应拦截器对于blob响应返回的是整个response对象
    // response.data应该是Blob类型
    if (!response || !response.data) {
      ElMessage.error(t('llmResources.downloadFailedFormat'))
      return
    }

    const blob = response.data instanceof Blob ? response.data : new Blob([response.data])

    // 验证blob是否有效
    if (!blob || blob.size === 0) {
      ElMessage.error(t('llmResources.downloadEmpty'))
      return
    }

    // 验证文件类型（Excel文件应该以ZIP格式开头）
    const firstBytes = await blob.slice(0, 4).arrayBuffer()
    const uint8Array = new Uint8Array(firstBytes)

    // Excel文件（.xlsx）是ZIP格式，ZIP文件头是 50 4B 03 04 (PK..)
    const isValidExcel =
      uint8Array[0] === 0x50 &&
      uint8Array[1] === 0x4b &&
      uint8Array[2] === 0x03 &&
      uint8Array[3] === 0x04

    if (!isValidExcel) {
      // 可能是错误响应，尝试读取错误信息
      const text = await blob.text()
      try {
        const errorData = JSON.parse(text)
        ElMessage.error(errorData.error || t('llmResources.downloadFailed'))
      } catch (e) {
        ElMessage.error(t('llmResources.downloadInvalidFormat'))
      }
      return
    }

    // 创建下载链接
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'llm_resources_import_template.xlsx'
    link.style.display = 'none'
    document.body.appendChild(link)
    link.click()

    // 清理
    setTimeout(() => {
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    }, 100)

    ElMessage.success('模板下载成功')
  } catch (error) {
    console.error('下载模板失败:', error)
    ElMessage.error(
      error.response?.data?.error || error.message || t('llmResources.downloadFailed')
    )
  }
}

// 上传前验证
const beforeUpload = file => {
  const isExcel =
    file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' ||
    file.type === 'application/vnd.ms-excel' ||
    file.name.endsWith('.xlsx') ||
    file.name.endsWith('.xls')
  if (!isExcel) {
    ElMessage.error(t('llmResources.uploadOnlyExcel'))
    return false
  }
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error(t('llmResources.fileSizeLimit'))
    return false
  }
  return true
}

// 导入成功
const handleImportSuccess = response => {
  if (response && response.success !== undefined) {
    const message = t('llmResources.importComplete', {
      success: response.success,
      fail: response.fail
    })
    if (response.fail > 0 && response.errors && response.errors.length > 0) {
      ElMessageBox.alert(
        `${message}\n\n${t('llmResources.errorDetails')}：\n${response.errors.slice(0, 10).join('\n')}${response.errors.length > 10 ? '\n...' : ''}`,
        t('llmResources.importResult'),
        {
          confirmButtonText: t('common.confirm'),
          type: response.success > 0 ? 'warning' : 'error'
        }
      )
    } else {
      ElMessage.success(message)
    }
    // 刷新列表
    getResourceList()
  } else {
    ElMessage.success(t('llmResources.importSuccess'))
    getResourceList()
  }
}

// 导入失败
const handleImportError = error => {
  console.error('导入失败:', error)
  const errorMsg = error.response?.data?.error || t('llmResources.importFailed')
  ElMessage.error(errorMsg)
}

// 一键填充 JSON 示例
const fillJsonExample = () => {
  jsonImportText.value = jsonImportExample.value
}

// 处理 JSON 批量导入
const handleJsonImport = async () => {
  if (!jsonImportText.value.trim()) {
    ElMessage.error(t('llmResources.jsonDataRequired'))
    return
  }

  let data
  try {
    data = JSON.parse(jsonImportText.value)
  } catch (e) {
    console.error('JSON 解析错误:', e)
    ElMessage.error(t('llmResources.jsonFormatError'))
    return
  }

  jsonImportLoading.value = true
  try {
    // 判断格式：Coding Plan 格式有 base_url 和 models 字段
    if (data.base_url && data.api_key && Array.isArray(data.models)) {
      // Coding Plan 格式
      const expandedModels = []
      data.models.forEach(m => {
        // 支持 types 数组或单个 type
        const types = m.types || (m.type ? [m.type] : ['chat'])
        types.forEach(type => {
          expandedModels.push({
            name: m.name?.trim(),
            type: type || 'chat',
            display_name: m.display_name?.trim() || m.name?.trim()
          })
        })
      })

      const codingPlanData = {
        base_url: data.base_url,
        api_key: data.api_key,
        models: expandedModels
      }

      const response = await api.importLLMResourcesFromBailian(codingPlanData)
      handleImportSuccess(response)
    } else if (Array.isArray(data) && data.length > 0) {
      // 数组格式
      const response = await api.importLLMResourcesByJSON(data)
      handleImportSuccess(response)
    } else {
      ElMessage.error(t('llmResources.jsonInvalidFormat'))
    }

    jsonImportDialogVisible.value = false
  } catch (error) {
    console.error('JSON 导入失败:', error)
  } finally {
    jsonImportLoading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  getResourceList()
})
</script>

<style scoped>
.llm-resources {
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

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
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

.json-import-example {
  margin-top: 16px;
  margin-bottom: 12px;
}

.json-import-example-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
}

.json-example-block {
  background-color: #f5f5f5;
  border-radius: 4px;
  padding: 10px 12px;
  font-family: Menlo, Monaco, Consolas, 'Courier New', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 260px;
  overflow: auto;
}

.model-item {
  margin-bottom: 12px;
}

.model-card {
  background-color: #fafafa;
}

.model-card :deep(.el-card__body) {
  padding: 16px;
}

.json-import-tabs {
  margin: 16px 0;
  text-align: center;
}
</style>
