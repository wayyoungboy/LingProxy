<template>
  <div class="endpoints-page">
    <div class="page-header">
      <h1 class="page-title">{{ $t('endpoints.title') }}</h1>
      <el-button type="primary" @click="openAddDialog">
        {{ $t('endpoints.addEndpoint') }}
      </el-button>
    </div>

    <!-- Search -->
    <div class="filter-section">
      <div class="search-filter">
        <el-form :inline="true" :model="searchForm">
          <el-form-item :label="$t('endpoints.name')">
            <el-input v-model="searchForm.name" :placeholder="$t('endpoints.namePlaceholder')" />
          </el-form-item>
          <el-form-item :label="$t('endpoints.status')">
            <el-select v-model="searchForm.status" :placeholder="$t('endpoints.selectStatus')">
              <el-option :label="$t('endpoints.enabled')" value="enabled" />
              <el-option :label="$t('endpoints.disabled')" value="disabled" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="getEndpoints">{{ $t('endpoints.query') }}</el-button>
            <el-button @click="resetForm">{{ $t('endpoints.reset') }}</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <!-- Table -->
    <div class="table-section">
      <el-table :data="endpointsList" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" :label="$t('endpoints.id')" width="80" />
        <el-table-column prop="name" :label="$t('endpoints.endpointName')" />
        <el-table-column prop="path" :label="$t('endpoints.path')" />
        <el-table-column prop="model" :label="$t('endpoints.model')" />
        <el-table-column prop="method" :label="$t('endpoints.method')" width="100" />
        <el-table-column prop="status" :label="$t('endpoints.status')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'enabled' ? 'success' : 'danger'" size="small">
              {{
                scope.row.status === 'enabled' ? $t('endpoints.enabled') : $t('endpoints.disabled')
              }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('endpoints.createdAt')" width="180">
          <template #default="scope">
            <span class="text-caption">{{ scope.row.created_at }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('endpoints.actions')" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="openEditDialog(scope.row)">
              {{ $t('endpoints.edit') }}
            </el-button>
            <el-button type="danger" size="small" @click="deleteEndpoint(scope.row.id)">
              {{ $t('endpoints.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.size"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="pagination.total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? $t('endpoints.addEndpoint') : $t('endpoints.editEndpoint')"
      width="600px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('endpoints.endpointName')" prop="name">
          <el-input v-model="form.name" :placeholder="$t('endpoints.namePlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('endpoints.path')" prop="path">
          <el-input v-model="form.path" :placeholder="$t('endpoints.pathPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('endpoints.model')" prop="model">
          <el-select v-model="form.model" :placeholder="$t('endpoints.selectModel')">
            <el-option
              v-for="model in models"
              :key="model.id"
              :label="model.name"
              :value="model.name"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('endpoints.method')" prop="method">
          <el-select v-model="form.method" :placeholder="$t('endpoints.selectMethod')">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('endpoints.status')" prop="status">
          <el-select v-model="form.status" :placeholder="$t('endpoints.selectStatus')">
            <el-option :label="$t('endpoints.enabled')" value="enabled" />
            <el-option :label="$t('endpoints.disabled')" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('endpoints.description')">
          <el-input
            v-model="form.description"
            type="textarea"
            rows="3"
            :placeholder="$t('endpoints.descriptionPlaceholder')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('endpoints.cancel') }}</el-button>
          <el-button type="primary" @click="submitForm">{{ $t('endpoints.confirm') }}</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import api from '../api'

const { t } = useI18n()

const endpointsList = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogType = ref('add')
const formRef = ref(null)
const models = ref([])

const searchForm = reactive({
  name: '',
  status: ''
})

const pagination = reactive({
  current: 1,
  size: 10,
  total: 0
})

const form = reactive({
  id: '',
  name: '',
  path: '',
  model: '',
  method: 'POST',
  status: 'enabled',
  description: ''
})

const rules = computed(() => ({
  name: [{ required: true, message: t('endpoints.nameRequired'), trigger: 'blur' }],
  path: [{ required: true, message: t('endpoints.pathRequired'), trigger: 'blur' }],
  model: [{ required: true, message: t('endpoints.modelRequired'), trigger: 'change' }],
  method: [{ required: true, message: t('endpoints.methodRequired'), trigger: 'change' }],
  status: [{ required: true, message: t('endpoints.statusRequired'), trigger: 'change' }]
}))

const getEndpoints = async () => {
  loading.value = true
  try {
    const response = await api.getEndpoints({
      page: pagination.current,
      page_size: pagination.size,
      name: searchForm.name,
      status: searchForm.status
    })
    endpointsList.value = response.data.items
    pagination.total = response.data.total
  } catch (error) {
    ElMessage.error(t('endpoints.getListFailed'))
  } finally {
    loading.value = false
  }
}

const getModels = async () => {
  try {
    const response = await api.getModels({ page: 1, page_size: 100 })
    models.value = response.data.items
  } catch (error) {
    ElMessage.error(t('endpoints.getModelsFailed'))
  }
}

const openAddDialog = () => {
  dialogType.value = 'add'
  Object.assign(form, {
    id: '',
    name: '',
    path: '',
    model: '',
    method: 'POST',
    status: 'enabled',
    description: ''
  })
  dialogVisible.value = true
}

const openEditDialog = row => {
  dialogType.value = 'edit'
  Object.assign(form, row)
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async valid => {
    if (valid) {
      try {
        if (dialogType.value === 'add') {
          await api.createEndpoint(form)
          ElMessage.success(t('endpoints.addSuccess'))
        } else {
          await api.updateEndpoint(form.id, form)
          ElMessage.success(t('endpoints.updateSuccess'))
        }
        dialogVisible.value = false
        getEndpoints()
      } catch (error) {
        ElMessage.error(
          dialogType.value === 'add' ? t('endpoints.addFailed') : t('endpoints.updateFailed')
        )
      }
    }
  })
}

const deleteEndpoint = async id => {
  try {
    await ElMessageBox.confirm(t('endpoints.deleteConfirmMessage'), t('endpoints.deleteConfirm'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    await api.deleteEndpoint(id)
    ElMessage.success(t('endpoints.deleteSuccess'))
    getEndpoints()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(t('endpoints.deleteFailed'))
    }
  }
}

const resetForm = () => {
  Object.assign(searchForm, {
    name: '',
    status: ''
  })
  getEndpoints()
}

const handleSizeChange = size => {
  pagination.size = size
  getEndpoints()
}

const handleCurrentChange = current => {
  pagination.current = current
  getEndpoints()
}

onMounted(() => {
  getEndpoints()
  getModels()
})
</script>

<style scoped>
.endpoints-page {
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

.filter-section {
  margin-bottom: 20px;
  background: var(--claude-ivory);
  border: 1px solid var(--claude-border-cream);
  border-radius: var(--radius-comfortable);
  padding: 20px;
}

.search-filter {
  display: flex;
  gap: 12px;
  align-items: center;
}

.table-section {
  margin-top: 20px;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.text-caption {
  font-size: 14px;
  color: var(--claude-text-secondary);
}
</style>
