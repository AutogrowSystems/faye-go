package memory

import (
	"container/list"
	"github.com/AutogrowSystems/faye-go/protocol"
	"sync"
)

type MemoryMsgStore struct {
	msgs *list.List
	lock *sync.RWMutex
}

func NewMemoryMsgStore() *MemoryMsgStore {
	return &MemoryMsgStore{list.New(), &sync.RWMutex{}}
}

func (m *MemoryMsgStore) EnqueueMessages(msgs []protocol.Message) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, msg := range msgs {
		m.msgs.PushBack(msg)
	}
}

func (m *MemoryMsgStore) GetAndClearMessages() []protocol.Message {
	m.lock.Lock()
	defer m.lock.Unlock()
	var msgArray = make([]protocol.Message, m.msgs.Len())
	i := 0
	for e := m.msgs.Front(); e != nil; e = e.Next() {
		msgArray[i] = e.Value.(protocol.Message)
		i = i + 1
	}
	m.msgs = &list.List{}
	return msgArray
}
