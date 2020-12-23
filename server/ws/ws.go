package ws

import (
	"github.com/ofavor/socket-gw/server"
)

type wsServer struct {
	opts server.Options
}

func newWsServer(opts ...server.Option) server.Server {
	return &wsServer{}
}

// NewServer create new websocket server
func NewServer(opts ...server.Option) server.Server {
	return newWsServer(opts...)
}
