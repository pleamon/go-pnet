package pnet

import "sync"

type ClientPool struct {
	mutex *sync.Mutex
	pool  map[string]*ClientInfo
}

func NewClientPool() *ClientPool {
	clientPool := &ClientPool{
		mutex: &sync.Mutex{},
		pool:  make(map[string]*ClientInfo),
	}
	return clientPool
}

func (cp *ClientPool) lock() {
	cp.mutex.Lock()
}

func (cp *ClientPool) unlock() {
	cp.mutex.Unlock()
}

func (cp *ClientPool) Set(name string, clientInfo *ClientInfo) {
	cp.lock()
	cp.pool[name] = clientInfo
	cp.unlock()
}

func (cp *ClientPool) Del(name string) {
	cp.lock()
	if _, ok := cp.pool[name]; ok {
		delete(cp.pool, name)
	}
	cp.unlock()
}

func (cp *ClientPool) Get(name string) *ClientInfo {
	return cp.pool[name]
}
