// 全局配置
const API_BASE_URL = '/api/v1';

// 工具函数
function showMessage(message, type = 'info') {
    const messageDiv = document.createElement('div');
    messageDiv.className = `message message-${type}`;
    messageDiv.textContent = message;
    document.body.appendChild(messageDiv);
    
    setTimeout(() => {
        messageDiv.remove();
    }, 3000);
}

function showLoading(element) {
    element.innerHTML = '<div class="loading">加载中...</div>';
}

function hideLoading(element) {
    element.querySelector('.loading')?.remove();
}

// 页面切换
function showSection(sectionId) {
    // 隐藏所有section
    const sections = document.querySelectorAll('.section');
    sections.forEach(section => section.classList.remove('active'));
    
    // 显示目标section
    document.getElementById(sectionId).classList.add('active');
    
    // 更新导航按钮状态
    const navBtns = document.querySelectorAll('.nav-btn');
    navBtns.forEach(btn => btn.classList.remove('active'));
    event.target.classList.add('active');
    
    // 加载对应数据
    loadSectionData(sectionId);
}

// 加载各section数据
function loadSectionData(sectionId) {
    switch(sectionId) {
        case 'dashboard':
            loadDashboardData();
            break;
        case 'users':
            loadUsers();
            break;
        case 'providers':
            loadLLMResources();
            break;
        case 'requests':
            loadRequests();
            break;
    }
}

// API调用
async function apiCall(method, url, data = null) {
    try {
        const options = {
            method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
            }
        };
        
        if (data) {
            options.body = JSON.stringify(data);
        }
        
        const response = await fetch(API_BASE_URL + url, options);
        const result = await response.json();
        
        if (!response.ok) {
            throw new Error(result.message || '请求失败');
        }
        
        return result;
    } catch (error) {
        showMessage(error.message, 'error');
        throw error;
    }
}

// 仪表盘数据
async function loadDashboardData() {
    try {
        // 模拟数据（实际项目中从API获取）
        document.getElementById('totalUsers').textContent = '12';
        document.getElementById('totalLLMResources').textContent = '3';
        document.getElementById('totalRequests').textContent = '1,234';
        document.getElementById('successRate').textContent = '98.5%';
        
        // 加载最近请求
        loadRecentRequests();
    } catch (error) {
        console.error('加载仪表盘数据失败:', error);
    }
}

async function loadRecentRequests() {
    const container = document.getElementById('recentRequests');
    showLoading(container);
    
    try {
        // 模拟数据
        const requests = [
            { id: 1, user: 'user1', provider: 'OpenAI', status: 'success', time: '2024-01-29 10:30:00' },
            { id: 2, user: 'user2', provider: 'Claude', status: 'success', time: '2024-01-29 10:25:00' },
            { id: 3, user: 'user3', provider: 'Gemini', status: 'failed', time: '2024-01-29 10:20:00' }
        ];
        
        container.innerHTML = `
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>用户</th>
                        <th>LLM资源</th>
                        <th>状态</th>
                        <th>时间</th>
                    </tr>
                </thead>
                <tbody>
                    ${requests.map(req => `
                        <tr>
                            <td>${req.id}</td>
                            <td>${req.user}</td>
                            <td>${req.provider}</td>
                            <td><span class="status ${req.status}">${req.status}</span></td>
                            <td>${req.time}</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
    } catch (error) {
        container.innerHTML = '<div class="error">加载失败</div>';
    }
}

// 用户管理
async function loadUsers() {
    const tbody = document.getElementById('usersTableBody');
    showLoading(tbody);
    
    try {
        // 模拟数据
        const users = [
            { id: 1, username: 'admin', status: 'active', created_at: '2024-01-01' },
            { id: 2, username: 'user1', status: 'active', created_at: '2024-01-15' },
            { id: 3, username: 'user2', status: 'inactive', created_at: '2024-01-20' }
        ];
        
        tbody.innerHTML = users.map(user => `
            <tr>
                <td>${user.id}</td>
                <td>${user.username}</td>
                <td><span class="status ${user.status}">${user.status}</span></td>
                <td>${user.created_at}</td>
                <td>
                    <button class="btn small" onclick="editUser(${user.id})")">编辑</button>
                    <button class="btn small danger" onclick="deleteUser(${user.id})")">删除</button>
                </td>
            </tr>
        `).join('');
    } catch (error) {
        tbody.innerHTML = '<tr><td colspan="6" class="error">加载失败</td></tr>';
    }
}

function addUser() {
    showModal('添加用户', `
        <form onsubmit="createUser(event)">
            <div class="form-group">
                <label>用户名:</label>
                <input type="text" id="newUsername" required>
            </div>
            <div class="form-group">
                <label>密码:</label>
                <input type="password" id="newPassword" required>
            </div>
            <button type="submit" class="btn primary">创建用户</button>
        </form>
    `);
}

function createUser(event) {
    event.preventDefault();
    const username = document.getElementById('newUsername').value;
    const password = document.getElementById('newPassword').value;
    
    // 模拟创建用户
    showMessage('用户创建成功');
    closeModal();
    loadUsers();
}

function editUser(id) {
    showModal('编辑用户', `
        <form onsubmit="updateUser(${id}, event)">
            <div class="form-group">
                <label>用户名:</label>
                <input type="text" id="editUsername" value="user${id}" required>
            </div>
            <button type="submit" class="btn primary">更新用户</button>
        </form>
    `);
}

function updateUser(id, event) {
    event.preventDefault();
    showMessage('用户更新成功');
    closeModal();
    loadUsers();
}

function deleteUser(id) {
    if (confirm('确定要删除这个用户吗？')) {
        showMessage('用户删除成功');
        loadUsers();
    }
}

// LLM资源管理
async function loadLLMResources() {
    const tbody = document.getElementById('llmResourcesTableBody');
    showLoading(tbody);
    
    try {
        // 模拟数据
        const resources = [
            { id: 1, name: 'OpenAI', type: 'openai', base_url: 'https://api.openai.com/v1', api_key: 'sk-****************************************************************************************************', status: 'active', created_at: '2024-01-01' },
            { id: 2, name: 'Claude', type: 'anthropic', base_url: 'https://api.anthropic.com/v1', api_key: 'sk-****************************************************************************************************', status: 'active', created_at: '2024-01-15' },
            { id: 3, name: 'Gemini', type: 'google', base_url: 'https://generativelanguage.googleapis.com/v1', api_key: '****************************************************************************************************', status: 'inactive', created_at: '2024-01-20' }
        ];
        
        tbody.innerHTML = resources.map(resource => `
            <tr>
                <td>${resource.id}</td>
                <td>${resource.name}</td>
                <td>${resource.type}</td>
                <td>${resource.base_url}</td>
                <td>${maskAPIKey(resource.api_key)}</td>
                <td><span class="status ${resource.status}">${resource.status}</span></td>
                <td>${resource.created_at}</td>
                <td>
                    <button class="btn small" onclick="editLLMResource(${resource.id})")">编辑</button>
                    <button class="btn small danger" onclick="deleteLLMResource(${resource.id})")">删除</button>
                </td>
            </tr>
        `).join('');
    } catch (error) {
        tbody.innerHTML = '<tr><td colspan="8" class="error">加载失败</td></tr>';
    }
}

// 加密API Key显示
function maskAPIKey(apiKey) {
    if (apiKey.length <= 8) {
        return apiKey;
    }
    const prefix = apiKey.substring(0, 4);
    const suffix = apiKey.substring(apiKey.length - 4);
    const masked = '*'.repeat(apiKey.length - 8);
    return prefix + masked + suffix;
}

function addLLMResource() {
    showModal('添加LLM资源', `
        <form onsubmit="createLLMResource(event)">
            <div class="form-group">
                <label>名称:</label>
                <input type="text" id="newResourceName" required>
            </div>
            <div class="form-group">
                <label>类型:</label>
                <select id="newResourceType" required>
                    <option value="openai">OpenAI</option>
                    <option value="anthropic">Anthropic</option>
                    <option value="google">Google</option>
                </select>
            </div>
            <div class="form-group">
                <label>模型:</label>
                <input type="text" id="newResourceModel" required placeholder="例如: gpt-3.5-turbo, glm-4.5-flash">
            </div>
            <div class="form-group">
                <label>Base URL:</label>
                <input type="text" id="newResourceBaseURL" required>
            </div>
            <div class="form-group">
                <label>API Key:</label>
                <input type="password" id="newResourceAPIKey" required>
            </div>
            <div class="form-group">
                <label>状态:</label>
                <select id="newResourceStatus" required>
                    <option value="active">活跃</option>
                    <option value="inactive">非活跃</option>
                </select>
            </div>
            <button type="submit" class="btn primary">创建LLM资源</button>
        </form>
    `);
}

function createLLMResource(event) {
    event.preventDefault();
    const name = document.getElementById('newResourceName').value;
    const type = document.getElementById('newResourceType').value;
    const model = document.getElementById('newResourceModel').value;
    const baseURL = document.getElementById('newResourceBaseURL').value;
    const apiKey = document.getElementById('newResourceAPIKey').value;
    const status = document.getElementById('newResourceStatus').value;
    
    // 模拟创建LLM资源
    showMessage('LLM资源创建成功');
    closeModal();
    loadLLMResources();
}

async function editLLMResource(id) {
    try {
        // 显示加载中状态
        showModal('编辑LLM资源', '<div class="loading">加载资源信息中...</div>');
        
        // 获取LLM资源详情
        const response = await fetch(`/api/v1/llm-resources/${id}`);
        if (!response.ok) {
            if (response.status === 404) {
                throw new Error('LLM资源不存在');
            }
            throw new Error('获取资源信息失败');
        }
        const data = await response.json();
        const resource = data.data;
        
        // 显示编辑表单并填充现有数据
        showModal('编辑LLM资源', `
            <form onsubmit="updateLLMResource(${id}, event)">
                <div class="form-group">
                    <label>名称:</label>
                    <input type="text" id="editResourceName" value="${resource.name || ''}" required>
                </div>
                <div class="form-group">
                    <label>类型:</label>
                    <input type="text" id="editResourceType" value="${resource.type || ''}" required>
                </div>
                <div class="form-group">
                    <label>模型:</label>
                    <input type="text" id="editResourceModel" value="${resource.model || ''}" required placeholder="例如: gpt-3.5-turbo, glm-4.5-flash">
                </div>
                <div class="form-group">
                    <label>Base URL:</label>
                    <input type="text" id="editResourceBaseURL" value="${resource.base_url || ''}" required>
                </div>
                <div class="form-group">
                    <label>API Key:</label>
                    <input type="password" id="editResourceAPIKey" value="${resource.api_key || ''}" required>
                </div>
                <div class="form-group">
                    <label>状态:</label>
                    <select id="editResourceStatus" required>
                        <option value="active" ${resource.status === 'active' ? 'selected' : ''}>活跃</option>
                        <option value="inactive" ${resource.status === 'inactive' ? 'selected' : ''}>非活跃</option>
                    </select>
                </div>
                <button type="submit" class="btn primary">更新LLM资源</button>
            </form>
        `);
    } catch (error) {
        console.error('加载资源信息失败:', error);
        showMessage(error.message, 'error');
        closeModal();
    }
}

function updateLLMResource(id, event) {
    event.preventDefault();
    const name = document.getElementById('editResourceName').value;
    const type = document.getElementById('editResourceType').value;
    const model = document.getElementById('editResourceModel').value;
    const baseURL = document.getElementById('editResourceBaseURL').value;
    const apiKey = document.getElementById('editResourceAPIKey').value;
    const status = document.getElementById('editResourceStatus').value;
    
    // 模拟更新LLM资源
    showMessage('LLM资源更新成功');
    closeModal();
    loadLLMResources();
}

function deleteLLMResource(id) {
    if (confirm('确定要删除这个LLM资源吗？')) {
        showMessage('LLM资源删除成功');
        loadLLMResources();
    }
}

// 请求记录
async function loadRequests() {
    const tbody = document.getElementById('requestsTableBody');
    showLoading(tbody);
    
    try {
        // 模拟数据
        const requests = [
            { id: 1, user: 'user1', provider: 'OpenAI', status: 'success', duration: '120ms', time: '2024-01-29 10:30:00' },
            { id: 2, user: 'user2', provider: 'Claude', status: 'success', duration: '95ms', time: '2024-01-29 10:25:00' },
            { id: 3, user: 'user3', provider: 'Gemini', status: 'failed', duration: '500ms', time: '2024-01-29 10:20:00' }
        ];
        
        tbody.innerHTML = requests.map(req => `
            <tr>
                <td>${req.id}</td>
                <td>${req.user}</td>
                <td>${req.provider}</td>
                <td><span class="status ${req.status}">${req.status}</span></td>
                <td>${req.duration}</td>
                <td>${req.time}</td>
            </tr>
        `).join('');
    } catch (error) {
        tbody.innerHTML = '<tr><td colspan="6" class="error">加载失败</td></tr>';
    }
}

// 系统设置
function saveSettings() {
    showMessage('设置保存成功');
}

// 模态框
function showModal(title, content) {
    const modal = document.getElementById('modal');
    const modalBody = document.getElementById('modalBody');
    
    modalBody.innerHTML = `<h3>${title}</h3>${content}`;
    modal.style.display = 'block';
}

function closeModal() {
    document.getElementById('modal').style.display = 'none';
}

// 刷新功能
function refreshUsers() {
    loadUsers();
    showMessage('用户列表已刷新');
}

function refreshLLMResources() {
    loadLLMResources();
    showMessage('LLM资源列表已刷新');
}

function refreshRequests() {
    loadRequests();
    showMessage('请求记录已刷新');
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    // 加载初始数据
    loadDashboardData();
    
    // 关闭模态框点击事件
    window.onclick = function(event) {
        const modal = document.getElementById('modal');
        if (event.target === modal) {
            closeModal();
        }
    };
});