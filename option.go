package gw

import (
	"github.com/ofavor/socket-gw/internal/log"
	"github.com/ofavor/socket-gw/server"
	"github.com/ofavor/socket-gw/session"
	"github.com/ofavor/socket-gw/transport"
)

// Options for gateway
type Options struct {
	Server server.Server
}

// Option function to set gateway options
type Option func(opts *Options)

func defaultOptions() Options {
	return Options{
		Server: server.NewServer(),
	}
}

// LogLevel set log level
func LogLevel(lv string) Option {
	return func(opts *Options) {
		log.SetLevel(lv)
	}
}

// Address set server address
func Address(addr string) Option {
	return func(opts *Options) {
		opts.Server.Init(server.Address(addr))
	}
}

// SessionAuth set session auth flag
func SessionAuth(auth bool) Option {
	return func(opts *Options) {
		opts.Server.Init(server.SessionAuth(auth))
	}
}

// SessionHandler set session handler
func SessionHandler(h session.Handler) Option {
	return func(opts *Options) {
		opts.Server.Init(server.SessionHandler(h))
	}
}

// Transport set server transport
func Transport(tp transport.Transport) Option {
	return func(opts *Options) {
		opts.Server.Init(server.Transport(tp))
	}
}
