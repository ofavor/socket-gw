package client

import (
	"github.com/ofavor/socket-gw/transport"
)

// Client interface
type Client interface {
	Connect() error

	Send(*transport.Packet) error

	Recv() (*transport.Packet, error)

	Close() error
}

// NewClient create new client
func NewClient(opts ...Option) Client {
	return newTCPClient(opts...)
}
