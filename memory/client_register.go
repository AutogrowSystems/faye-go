package memory

import (
	"github.com/autogrowsystems/faye-go/protocol"
	"github.com/autogrowsystems/faye-go/utils"
	creg "github.com/roncohen/cleaningRegister"
	"sync"
	"time"
)

type ClientRegister struct {
	clients              *creg.CleaningRegister
	subscriptionRegister *SubscriptionRegister
	logger               utils.Logger
	lock                 *sync.RWMutex
}

func NewClientRegister(logger utils.Logger) *ClientRegister {
	subReg := NewSubscriptionRegister()

	shouldRemove := func(key interface{}, item interface{}) bool {
		client := item.(*protocol.Client)
		return client.IsExpired()
	}

	removed := func(key interface{}, item interface{}) {
		client := item.(*protocol.Client)
		logger.Infof("Removing client %s due to inactivity", client.Id())
		subReg.RemoveClient(client.Id())
	}

	clientreg := ClientRegister{
		clients:              creg.New(1*time.Minute, shouldRemove, removed),
		subscriptionRegister: subReg,
		logger:               logger,
		lock:                 &sync.RWMutex{},
	}

	return &clientreg
}

func (cr ClientRegister) AddClient(client *protocol.Client) {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	cr.clients.Put(client.Id(), client)
}

func (cr ClientRegister) removeClient(clientId string) {
	// TODO: More cleanups
	cr.lock.Lock()
	defer cr.lock.Unlock()
	cr.subscriptionRegister.RemoveClient(clientId)
}

func (cr ClientRegister) GetClient(clientId string) *protocol.Client {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	client, ok := cr.clients.Get(clientId)
	if ok {
		return client.(*protocol.Client)
	} else {
		return nil
	}
}

/* Front for SubscriptionRegister */

func (cr ClientRegister) AddSubscription(clientId string, patterns []string) {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	cr.subscriptionRegister.AddSubscription(clientId, patterns)
}

func (cr ClientRegister) RemoveSubscription(clientId string, patterns []string) {
	cr.lock.Lock()
	defer cr.lock.Unlock()
	cr.subscriptionRegister.RemoveSubscription(clientId, patterns)
}

func (cr ClientRegister) GetClients(patterns []string) []string {
	return cr.subscriptionRegister.GetClients(patterns)
}
