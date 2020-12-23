package server

import (
	"github.com/ofavor/socket-gw/session"
	"github.com/ofavor/socket-gw/transport"
)

// Options for server
type Options struct {
	Addr string

	SessionAuth    bool
	SessionHandler session.Handler
	Transport      transport.Transport
}

// Option function to set server options
type Option func(opts *Options)

func defaultOptions() Options {
	return Options{
		Addr:           ":9999",
		SessionAuth:    false,
		SessionHandler: nil,
		Transport:      transport.NewTransport(),
	}
}

// Address set server address, "host:port"
func Address(addr string) Option {
	return func(opts *Options) {
		opts.Addr = addr
	}
}

// SessionAuth set session auth flag
func SessionAuth(auth bool) Option {
	return func(opts *Options) {
		opts.SessionAuth = auth
	}
}

// SessionHandler set session handler
func SessionHandler(sh session.Handler) Option {
	return func(opts *Options) {
		opts.SessionHandler = sh
	}
}

// Transport set transport
func Transport(tp transport.Transport) Option {
	return func(opts *Options) {
		opts.Transport = tp
	}
}
