package client

import (
	"sync"
)

// ClientManager 管理不同类型的客户端实例
type ClientManager struct {
	clients map[string]Client
	mu      sync.RWMutex
}

// NewClientManager 创建新的客户端管理器
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]Client),
	}
}

// Register 注册客户端实例
func (cm *ClientManager) Register(name string, client Client) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.clients[name] = client
}

// Get 获取客户端实例
func (cm *ClientManager) Get(name string) Client {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.clients[name]
}

// GetOpenAIClient 获取 OpenAI 客户端实例
func (cm *ClientManager) GetOpenAIClient(name string) OpenAIClient {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if client, ok := cm.clients[name]; ok {
		if openAIClient, ok := client.(OpenAIClient); ok {
			return openAIClient
		}
	}
	return nil
}

// Remove 移除客户端实例
func (cm *ClientManager) Remove(name string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if client, ok := cm.clients[name]; ok {
		client.Close()
		delete(cm.clients, name)
	}
}

// List 列出所有客户端实例
func (cm *ClientManager) List() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	var names []string
	for name := range cm.clients {
		names = append(names, name)
	}
	return names
}

// Close 关闭所有客户端实例
func (cm *ClientManager) Close() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for name, client := range cm.clients {
		client.Close()
		delete(cm.clients, name)
	}
}
