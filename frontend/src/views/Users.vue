<template>
  <div class="users-page">
    <div class="page-header">
      <h1 class="page-title">{{ $t('users.title') }}</h1>
      <el-button type="primary" @click="handleAddUser">
        <el-icon><Plus /></el-icon>
        {{ $t('users.addUser') }}
      </el-button>
    </div>

    <!-- Search -->
    <el-card class="filter-card">
      <div class="search-filter">
        <el-input
          v-model="searchQuery"
          :placeholder="$t('users.searchUsername')"
          prefix-icon="Search"
          style="width: 240px"
        ></el-input>
        <el-select v-model="statusFilter" :placeholder="$t('users.filterStatus')" style="width: 140px">
          <el-option :label="$t('users.all')" value=""></el-option>
          <el-option :label="$t('users.active')" value="active"></el-option>
          <el-option :label="$t('users.inactive')" value="inactive"></el-option>
        </el-select>
      </div>
    </el-card>

    <!-- Users Table -->
    <el-card class="table-card">
      <el-table v-loading="loading" :data="filteredUsers" style="width: 100%">
        <el-table-column prop="id" :label="$t('users.id')" width="180" />
        <el-table-column prop="username" :label="$t('users.username')" />
        <el-table-column prop="role" :label="$t('users.role')" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.role === 'admin' ? 'danger' : 'primary'" size="small">
              {{ scope.row.role === 'admin' ? $t('users.admin') : $t('users.user') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="api_key" :label="$t('users.apiKey')" width="300">
          <template #default="scope">
            <el-tooltip :content="$t('users.clickToCopy')" placement="top">
              <span @click="copyApiKey(scope.row.api_key)" class="api-key-text">
                {{ scope.row.api_key || $t('users.notGenerated') }}
              </span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('users.status')" width="100">
          <template #default="scope">
            <el-tag
              :type="scope.row.status === 'active' ? 'success' : scope.row.status === 'suspended' ? 'warning' : 'danger'"
              size="small"
            >
              {{
                scope.row.status === 'active'
                  ? $t('users.active')
                  : scope.row.status === 'suspended'
                    ? $t('users.suspended')
                    : $t('users.inactive')
              }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('users.createdAt')" width="180">
          <template #default="scope">
            <span class="text-caption">{{ formatDate(scope.row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('users.actions')" width="280">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEditUser(scope.row)">
              {{ $t('users.edit') }}
            </el-button>
            <el-button type="default" size="small" @click="handleResetAPIKey(scope.row.id)">
              {{ $t('users.resetKey') }}
            </el-button>
            <el-button type="danger" size="small" @click="handleDeleteUser(scope.row.id)">
              {{ $t('users.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

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

    <!-- Add/Edit Dialog -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form ref="userFormRef" :model="userForm" :rules="userRules" label-width="80px">
        <el-form-item :label="$t('users.username')" prop="username">
          <el-input v-model="userForm.username" :placeholder="$t('users.usernamePlaceholder')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('users.password')" prop="password" v-if="isAddMode">
          <el-input v-model="userForm.password" type="password" :placeholder="$t('users.passwordPlaceholder')" show-password></el-input>
        </el-form-item>
        <el-form-item :label="$t('users.role')" prop="role">
          <el-select v-model="userForm.role" :placeholder="$t('users.selectRole')">
            <el-option :label="$t('users.admin')" value="admin"></el-option>
            <el-option :label="$t('users.user')" value="user"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('users.apiKey')" prop="api_key">
          <el-input v-model="userForm.api_key" :placeholder="$t('users.apiKeyPlaceholder')" :disabled="true">
            <template #append>
              <el-button @click="handleResetAPIKey(userForm.id)" type="text" :disabled="isAddMode">
                {{ $t('users.reset') }}
              </el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('users.status')" prop="status">
          <el-select v-model="userForm.status" :placeholder="$t('users.selectStatus')">
            <el-option :label="$t('users.active')" value="active"></el-option>
            <el-option :label="$t('users.inactive')" value="inactive"></el-option>
            <el-option :label="$t('users.suspended')" value="suspended"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">{{ $t('users.cancel') }}</el-button>
          <el-button type="primary" @click="handleSaveUser">{{ $t('users.confirm') }}</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, CopyDocument } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import api from '../api'

const { t } = useI18n()

const loading = ref(false)
const users = ref([])
const searchQuery = ref('')
const statusFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const dialogVisible = ref(false)
const dialogTitle = ref('')
const isAddMode = ref(false)
const userFormRef = ref(null)
const userForm = reactive({
  id: '',
  username: '',
  password: '',
  role: 'user',
  api_key: '',
  status: 'active'
})

const userRules = computed(() => ({
  username: [
    { required: true, message: t('users.usernameRequired'), trigger: 'blur' },
    { min: 3, max: 50, message: t('users.usernameLength'), trigger: 'blur' }
  ],
  password: [{ min: 6, message: t('users.passwordMinLength'), trigger: 'blur' }],
  role: [{ required: true, message: t('users.roleRequired'), trigger: 'change' }],
  status: [{ required: true, message: t('users.statusRequired'), trigger: 'change' }]
}))

const getUserList = async () => {
  try {
    loading.value = true
    const response = await api.getUsers()
    if (response && response.data) {
      users.value = Array.isArray(response.data) ? response.data : []
      total.value = users.value.length
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error(t('users.getListFailed'))
  } finally {
    loading.value = false
  }
}

const filteredUsers = computed(() => {
  let result = users.value
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(user => user.username.toLowerCase().includes(query))
  }
  if (statusFilter.value) {
    result = result.filter(user => user.status === statusFilter.value)
  }
  const startIndex = (currentPage.value - 1) * pageSize.value
  const endIndex = startIndex + pageSize.value
  return result.slice(startIndex, endIndex)
})

const handleAddUser = () => {
  isAddMode.value = true
  dialogTitle.value = t('users.addUser')
  Object.assign(userForm, {
    id: '',
    username: '',
    password: '',
    role: 'user',
    api_key: '',
    status: 'active'
  })
  dialogVisible.value = true
}

const handleEditUser = user => {
  isAddMode.value = false
  dialogTitle.value = t('users.editUser')
  Object.assign(userForm, user)
  dialogVisible.value = true
}

const handleSaveUser = async () => {
  try {
    await userFormRef.value.validate()

    if (isAddMode.value) {
      await api.createUser(userForm)
      ElMessage.success(t('users.createSuccess'))
    } else {
      await api.updateUser(userForm.id, userForm)
      ElMessage.success(t('users.updateSuccess'))
    }

    dialogVisible.value = false
    getUserList()
  } catch (error) {
    console.error('保存用户失败:', error)
    ElMessage.error(t('users.saveFailed'))
  }
}

const handleDeleteUser = async id => {
  try {
    await ElMessageBox.confirm(t('users.deleteConfirmMessage'), t('users.deleteConfirm'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'danger'
    })

    await api.deleteUser(id)
    ElMessage.success(t('users.deleteSuccess'))
    getUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除用户失败:', error)
      ElMessage.error(t('users.deleteFailed'))
    }
  }
}

const handleResetAPIKey = async userId => {
  if (!userId) {
    ElMessage.warning(t('users.saveUserFirst'))
    return
  }

  try {
    await ElMessageBox.confirm(t('users.resetKeyConfirmMessage'), t('users.resetKeyConfirm'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })

    const response = await api.resetAPIKey(userId)
    if (response && response.data) {
      userForm.api_key = response.data.api_key
      ElMessage.success(t('users.resetKeySuccess'))
      getUserList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重置API Key失败:', error)
      ElMessage.error(t('users.resetKeyFailed'))
    }
  }
}

const copyApiKey = apiKey => {
  navigator.clipboard
    .writeText(apiKey)
    .then(() => {
      ElMessage.success(t('users.apiKeyCopied'))
    })
    .catch(() => {
      ElMessage.error(t('users.copyFailed'))
    })
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
  getUserList()
})
</script>

<style scoped>
.users-page {
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

.filter-card {
  margin-bottom: 16px;
}

.search-filter {
  display: flex;
  gap: 12px;
  align-items: center;
}

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

.api-key-text {
  cursor: pointer;
  color: var(--claude-text-secondary);
  transition: color 0.2s;
}

.api-key-text:hover {
  color: var(--claude-terracotta);
}

.text-caption {
  font-size: 14px;
  color: var(--claude-text-secondary);
}
</style>