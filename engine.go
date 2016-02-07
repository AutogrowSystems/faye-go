package faye

import (
	"fmt"
	"github.com/AutogrowSystems/faye-go/memory"
	"github.com/AutogrowSystems/faye-go/protocol"
	"github.com/AutogrowSystems/faye-go/transport"
	"github.com/AutogrowSystems/faye-go/utils"
	"strconv"
)

type Engine struct {
	ns      memory.MemoryNamespace
	clients *memory.ClientRegister
	logger  utils.Logger
}

func NewEngine(logger utils.Logger) Engine {
	return Engine{
		ns:      memory.NewMemoryNamespace(),
		clients: memory.NewClientRegister(logger),
		logger:  logger,
	}
}

func (m Engine) responseFromRequest(request protocol.Message) protocol.Message {
	response := protocol.Message{}
	response["channel"] = request.Channel().Name()
	if reqId, ok := request["id"]; ok {
		response["id"] = reqId.(string)
	}

	return response
}

func (m Engine) GetClient(clientId string) *protocol.Client {
	return m.clients.GetClient(clientId)
}

func (m Engine) NewClient(conn protocol.Connection) *protocol.Client {
	newClientId := m.ns.Generate()
	msgStore := memory.NewMemoryMsgStore()
	newClient := protocol.NewClient(newClientId, msgStore, m.logger)
	m.clients.AddClient(&newClient)
	return &newClient
}

func (m Engine) AddSubscription(clientId string, subscriptions []string) {
	m.logger.Infof("SUBSCRIBE %s subscription: %v", clientId, subscriptions)
	m.clients.AddSubscription(clientId, subscriptions)
}

func (m Engine) Handshake(request protocol.Message, conn protocol.Connection) string {
	newClientId := ""
	version := request["version"].(string)
	m.logger.Debugf("New handshake request received for Bayeux Protocol %s", version)

	response := m.responseFromRequest(request)
	response["successful"] = false

	if version == protocol.BAYEUX_VERSION {
		newClientId = m.NewClient(conn).Id()

		response.Update(map[string]interface{}{
			"clientId":                 newClientId,
			"channel":                  protocol.META_PREFIX + protocol.META_HANDSHAKE_CHANNEL,
			"version":                  protocol.BAYEUX_VERSION,
			"advice":                   protocol.DEFAULT_ADVICE,
			"supportedConnectionTypes": []string{"websocket", "long-polling"},
			"successful":               true,
		})

		m.logger.Debugf("New client given ID of %s, advised: %+v", newClientId, response["advice"])

	} else {
		response["error"] = fmt.Sprintf("Only supported version is '%s'", protocol.BAYEUX_VERSION)
	}

	// Answer directly
	conn.Send([]protocol.Message{response})
	return newClientId
}

func (m Engine) Connect(request protocol.Message, client *protocol.Client, conn protocol.Connection) {
	m.logger.Debugf("Connect request from %s", client.Id())
	response := m.responseFromRequest(request)
	response["successful"] = true

	timeout, _ := strconv.Atoi(protocol.DEFAULT_ADVICE["timeout"])

	response.Update(protocol.Message{
		"advice": protocol.DEFAULT_ADVICE,
	})
	client.Connect(timeout, 0, response, conn)
}

func (m Engine) SubscribeService(chanOut chan<- protocol.Message, subscription []string) {
	conn := transport.InternalConnection{chanOut}
	newClient := m.NewClient(conn)
	newClient.Connect(-1, 0, nil, conn)
	newClient.SetConnection(conn)
	m.AddSubscription(newClient.Id(), subscription)
}

func (m Engine) SubscribeClient(request protocol.Message, client *protocol.Client) {
	m.logger.Debugf("Subscribe request from %s", client.Id())
	response := m.responseFromRequest(request)
	response["successful"] = true

	subscription := request["subscription"]
	response["subscription"] = subscription

	var subs []string
	switch subscription.(type) {
	case []string:
		subs = subscription.([]string)
	case string:
		subs = []string{subscription.(string)}
	}

	for _, s := range subs {
		// Do not register clients subscribing to a service channel
		// They will be answered directly instead of through the normal subscription system
		if !protocol.NewChannel(s).IsService() {
			m.AddSubscription(client.Id(), []string{s})
		}
	}

	client.Queue(response)
}

func (m Engine) Disconnect(request protocol.Message, client *protocol.Client, conn protocol.Connection) {
	m.logger.Debugf("Disconnect request from %s", client.Id())
	response := m.responseFromRequest(request)
	response["successful"] = true
	clientId := request.ClientId()
	m.logger.Debugf("Client %s disconnected", clientId)
}

func (m Engine) Publish(request protocol.Message, conn protocol.Connection) {

	response := m.responseFromRequest(request)
	response["successful"] = true
	data := request["data"]
	channel := request.Channel()
	response.SetClientId(request.ClientId())

	if m.clients.GetClient(request.ClientId()) == nil {
		m.logger.Warnf("PUBLISH from unknown client %s", request)
		m.failAndAdvise("handshake", response, conn)
		return
	}

	m.logger.Debugf("Publish request from %s to %s", request.ClientId(), channel.Name())

	conn.Send([]protocol.Message{response})

	go func() {
		// Prepare msg to send to subscribers
		msg := protocol.Message{}
		msg["channel"] = channel.Name()
		msg["data"] = data

		// Get clients with subscriptions
		recipients := m.clients.GetClients(channel.Expand())
		m.logger.Debugf("PUBLISH from %s on %s to %d recipients", request.ClientId(), channel.Name(), len(recipients))
		// Queue messages
		for _, c := range recipients {
			m.clients.GetClient(c).Queue(msg)
		}
	}()

}

func (m Engine) failAndAdvise(reconnect string, response protocol.Message, conn protocol.Connection) {
	advice := protocol.DEFAULT_ADVICE
	advice["reconnect"] = reconnect
	response["advice"] = advice
	response["successful"] = false

	m.logger.Warnf("Advising %s to reconnect via %s", response.ClientId(), reconnect)
	conn.Send([]protocol.Message{response})
}

// Publish message directly to client
// msg should have "channel" which the client is expecting, e.g. "/service/echo"
func (m Engine) PublishFromService(recipientId string, msg protocol.Message) {
	// response["successful"] = true
	m.clients.GetClient(recipientId).Queue(msg)
}
