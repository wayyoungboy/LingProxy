<template>
  <div class="users">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" @click="handleAddUser">
            <el-icon><Plus /></el-icon>
            添加用户
          </el-button>
        </div>
      </template>
      
      <!-- 搜索和筛选 -->
      <div class="search-filter">
        <el-input
          v-model="searchQuery"
          placeholder="搜索用户名"
          prefix-icon="Search"
          style="width: 240px; margin-right: 10px"
        ></el-input>
        <el-select
          v-model="statusFilter"
          placeholder="筛选状态"
          style="width: 120px"
        >
          <el-option label="全部" value=""></el-option>
          <el-option label="活跃" value="active"></el-option>
          <el-option label="禁用" value="inactive"></el-option>
        </el-select>
      </div>
      
      <!-- 用户列表 -->
      <el-table
        v-loading="loading"
        :data="filteredUsers"
        style="width: 100%; margin-top: 20px"
        border
      >
        <el-table-column prop="id" label="ID" width="180" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="role" label="角色" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.role === 'admin' ? 'danger' : 'primary'">
              {{ scope.row.role === 'admin' ? '管理员' : '用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="api_key" label="API Key" width="300">
          <template #default="scope">
            <el-tooltip content="点击复制" placement="top">
              <span @click="copyApiKey(scope.row.api_key)" style="cursor: pointer;">
                {{ scope.row.api_key || '未生成' }}
              </span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag
              :type="scope.row.status === 'active' ? 'success' : scope.row.status === 'suspended' ? 'warning' : 'danger'"
            >
              {{ scope.row.status === 'active' ? '活跃' : scope.row.status === 'suspended' ? '暂停' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              @click="handleEditUser(scope.row)"
              style="margin-right: 5px"
            >
              编辑
            </el-button>
            <el-button
              type="warning"
              size="small"
              @click="handleResetAPIKey(scope.row.id)"
              style="margin-right: 5px"
            >
              重置Key
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDeleteUser(scope.row.id)"
            >
              删除
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
    
    <!-- 添加/编辑用户对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
    >
      <el-form
        ref="userFormRef"
        :model="userForm"
        :rules="userRules"
        label-width="80px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="isAddMode">
          <el-input
            v-model="userForm.password"
            type="password"
            placeholder="请输入密码（至少6位）"
            show-password
          ></el-input>
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="userForm.role" placeholder="请选择角色">
            <el-option label="管理员" value="admin"></el-option>
            <el-option label="普通用户" value="user"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="API Key" prop="api_key">
          <el-input
            v-model="userForm.api_key"
            placeholder="API Key（创建时自动生成）"
            :disabled="true"
          >
            <template #append>
              <el-button @click="handleResetAPIKey(userForm.id)" type="text" :disabled="isAddMode">
                重置
              </el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="userForm.status" placeholder="请选择状态">
            <el-option label="活跃" value="active"></el-option>
            <el-option label="禁用" value="inactive"></el-option>
            <el-option label="暂停" value="suspended"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSaveUser">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, CopyDocument } from '@element-plus/icons-vue'
import api from '../api'

const loading = ref(false)
const users = ref([])
const searchQuery = ref('')
const statusFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 对话框相关
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

// 表单验证规则
const userRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度在3到50个字符', trigger: 'blur' }
  ],
  password: [
    { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

// 获取用户列表
const getUserList = async () => {
  try {
    loading.value = true
    const response = await api.getUsers()
    if (response && response.data) {
      // 后端返回的是数组，不是分页对象
      users.value = Array.isArray(response.data) ? response.data : []
      total.value = users.value.length
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 过滤用户
const filteredUsers = computed(() => {
  let result = users.value
  
  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(user => 
      user.username.toLowerCase().includes(query)
    )
  }
  
  // 状态过滤
  if (statusFilter.value) {
    result = result.filter(user => user.status === statusFilter.value)
  }
  
  // 分页
  const startIndex = (currentPage.value - 1) * pageSize.value
  const endIndex = startIndex + pageSize.value
  return result.slice(startIndex, endIndex)
})

// 处理添加用户
const handleAddUser = () => {
  isAddMode.value = true
  dialogTitle.value = '添加用户'
  // 重置表单
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

// 处理编辑用户
const handleEditUser = (user) => {
  isAddMode.value = false
  dialogTitle.value = '编辑用户'
  // 填充表单
  Object.assign(userForm, user)
  dialogVisible.value = true
}

// 处理保存用户
const handleSaveUser = async () => {
  try {
    // 验证表单
    await userFormRef.value.validate()
    
    if (isAddMode.value) {
      // 创建用户
      await api.createUser(userForm)
      ElMessage.success('用户创建成功')
    } else {
      // 更新用户
      await api.updateUser(userForm.id, userForm)
      ElMessage.success('用户更新成功')
    }
    
    // 关闭对话框
    dialogVisible.value = false
    // 重新获取用户列表
    getUserList()
  } catch (error) {
    console.error('保存用户失败:', error)
    ElMessage.error('保存用户失败')
  }
}

// 处理删除用户
const handleDeleteUser = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这个用户吗？',
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'danger'
      }
    )
    
    await api.deleteUser(id)
    ElMessage.success('用户删除成功')
    // 重新获取用户列表
    getUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除用户失败:', error)
      ElMessage.error('删除用户失败')
    }
  }
}

// 重置API Key
const handleResetAPIKey = async (userId) => {
  if (!userId) {
    ElMessage.warning('请先保存用户')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      '确定要重置这个用户的API Key吗？重置后旧的API Key将失效。',
      '重置确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const response = await api.resetAPIKey(userId)
    if (response && response.data) {
      userForm.api_key = response.data.api_key
      ElMessage.success('API Key重置成功')
      // 重新获取用户列表
      getUserList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重置API Key失败:', error)
      ElMessage.error('重置API Key失败')
    }
  }
}

// 复制API Key
const copyApiKey = (apiKey) => {
  navigator.clipboard.writeText(apiKey).then(() => {
    ElMessage.success('API Key 已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
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
  getUserList()
})
</script>

<style scoped>
.users {
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
</style>