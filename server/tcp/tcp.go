package tcp

import (
	"github.com/ofavor/socket-gw/server"
)

// NewServer create new TCP server
func NewServer(opts ...server.Option) server.Server {
	return server.NewServer(opts...)
}
