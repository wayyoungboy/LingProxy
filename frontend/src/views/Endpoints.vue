<template>
  <div class="endpoints-container">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span>端点管理</span>
          <el-button type="primary" @click="openAddDialog">添加端点</el-button>
        </div>
      </template>
      
      <el-form :inline="true" :model="searchForm" class="mb-4">
        <el-form-item label="名称">
          <el-input v-model="searchForm.name" placeholder="请输入端点名称" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态">
            <el-option label="启用" value="enabled" />
            <el-option label="禁用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getEndpoints">查询</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
      
      <el-table :data="endpointsList" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="端点名称" />
        <el-table-column prop="path" label="路径" />
        <el-table-column prop="model" label="模型" />
        <el-table-column prop="method" label="方法" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'enabled' ? 'success' : 'danger'">
              {{ scope.row.status === 'enabled' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="openEditDialog(scope.row)">
              编辑
            </el-button>
            <el-button type="danger" size="small" @click="deleteEndpoint(scope.row.id)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination-container mt-4">
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
    </el-card>
    
    <!-- 添加/编辑端点对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '添加端点' : '编辑端点'"
      width="600px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="端点名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入端点名称" />
        </el-form-item>
        <el-form-item label="路径" prop="path">
          <el-input v-model="form.path" placeholder="请输入路径，如 /v1/chat/completions" />
        </el-form-item>
        <el-form-item label="模型" prop="model">
          <el-select v-model="form.model" placeholder="请选择模型">
            <el-option
              v-for="model in models"
              :key="model.id"
              :label="model.name"
              :value="model.name"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="方法" prop="method">
          <el-select v-model="form.method" placeholder="请选择方法">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" placeholder="请选择状态">
            <el-option label="启用" value="enabled" />
            <el-option label="禁用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="form.description"
            type="textarea"
            rows="3"
            placeholder="请输入端点描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

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

const rules = {
  name: [{ required: true, message: '请输入端点名称', trigger: 'blur' }],
  path: [{ required: true, message: '请输入路径', trigger: 'blur' }],
  model: [{ required: true, message: '请选择模型', trigger: 'change' }],
  method: [{ required: true, message: '请选择方法', trigger: 'change' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

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
    ElMessage.error('获取端点列表失败')
  } finally {
    loading.value = false
  }
}

const getModels = async () => {
  try {
    const response = await api.getModels({ page: 1, page_size: 100 })
    models.value = response.data.items
  } catch (error) {
    ElMessage.error('获取模型列表失败')
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

const openEditDialog = (row) => {
  dialogType.value = 'edit'
  Object.assign(form, row)
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (dialogType.value === 'add') {
          await api.createEndpoint(form)
          ElMessage.success('添加端点成功')
        } else {
          await api.updateEndpoint(form.id, form)
          ElMessage.success('更新端点成功')
        }
        dialogVisible.value = false
        getEndpoints()
      } catch (error) {
        ElMessage.error(dialogType.value === 'add' ? '添加端点失败' : '更新端点失败')
      }
    }
  })
}

const deleteEndpoint = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除该端点吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await api.deleteEndpoint(id)
    ElMessage.success('删除端点成功')
    getEndpoints()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除端点失败')
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

const handleSizeChange = (size) => {
  pagination.size = size
  getEndpoints()
}

const handleCurrentChange = (current) => {
  pagination.current = current
  getEndpoints()
}

onMounted(() => {
  getEndpoints()
  getModels()
})
</script>

<style scoped>
.endpoints-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
}
</style>