package pnet

import "sync"

// ClientPool 客户端管理池
type ClientPool struct {
	mutex *sync.Mutex
	pool  map[string]*ClientInfo
}

// NewClientPool 创建客户端管理池
func NewClientPool() *ClientPool {
	clientPool := &ClientPool{
		mutex: &sync.Mutex{},
		pool:  make(map[string]*ClientInfo),
	}
	return clientPool
}

// lock 加锁
func (cp *ClientPool) lock() {
	cp.mutex.Lock()
}

// unlock 解锁
func (cp *ClientPool) unlock() {
	cp.mutex.Unlock()
}

// Set 设置数据
func (cp *ClientPool) Set(name string, clientInfo *ClientInfo) {
	cp.lock()
	cp.pool[name] = clientInfo
	cp.unlock()
}

// Del 删除数据
func (cp *ClientPool) Del(name string) {
	cp.lock()
	if _, ok := cp.pool[name]; ok {
		delete(cp.pool, name)
	}
	cp.unlock()
}

// Get 获取数据
func (cp *ClientPool) Get(name string) (*ClientInfo, bool) {
	cp.lock()
	clientInfo, ok := cp.pool[name]
	cp.unlock()
	return clientInfo, ok
}
