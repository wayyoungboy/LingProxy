package balancer

import (
	"fmt"
	"sync"
)

// LoadBalancer 负载均衡器接口
type LoadBalancer interface {
	Select(servers []string) (string, error)
}

// RoundRobinBalancer 轮询负载均衡器
type RoundRobinBalancer struct {
	currentIndex int
	mutex        sync.RWMutex
}

// NewRoundRobinBalancer 创建新的轮询负载均衡器
func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{
		currentIndex: 0,
	}
}

// Select 选择下一个服务器
func (b *RoundRobinBalancer) Select(servers []string) (string, error) {
	if len(servers) == 0 {
		return "", ErrNoServersAvailable
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	server := servers[b.currentIndex]
	b.currentIndex = (b.currentIndex + 1) % len(servers)
	return server, nil
}

// ErrNoServersAvailable 没有可用服务器的错误
var ErrNoServersAvailable = fmt.Errorf("no servers available")