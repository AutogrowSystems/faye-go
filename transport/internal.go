package transport

import (
	"github.com/autogrowsystems/faye-go/protocol"
	// "github.com/autogrowsystems/faye-go/utils"
)

type InternalConnection struct {
	Channel chan<- protocol.Message
}

func (i InternalConnection) Send(msgs []protocol.Message) error {
	for _, m := range msgs {
		i.Channel <- m
	}
	return nil
}

func (i InternalConnection) IsConnected() bool {
	return true
}

func (i InternalConnection) IsSingleShot() bool {
	return false
}

func (i InternalConnection) Close() {
	close(i.Channel)
}

func (i InternalConnection) Priority() int {
	return 1
}
