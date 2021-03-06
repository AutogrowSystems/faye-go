package memory

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"sync"
)

type MemoryNamespace struct {
	idMap  map[string]bool
	lastId int
	mutex  sync.RWMutex
}

func NewMemoryNamespace() MemoryNamespace {
	return MemoryNamespace{
		idMap: make(map[string]bool),
		mutex: sync.RWMutex{},
	}
}

func (m MemoryNamespace) IsUsed(id string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, ok := m.idMap[id]
	return ok
}

func (m MemoryNamespace) generate() string {
	uuid := make([]byte, 16)
	io.ReadFull(rand.Reader, uuid)
	return hex.EncodeToString(uuid)
}

func (m MemoryNamespace) Generate() string {
	for {
		newId := m.generate()
		if !m.IsUsed(newId) {
			m.mutex.Lock()
			m.idMap[newId] = true
			m.mutex.Unlock()
			return newId
		}
	}

}

func (m MemoryNamespace) Expire(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.idMap, id)
}
