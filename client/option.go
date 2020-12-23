package client

import (
	"github.com/ofavor/socket-gw/internal/log"
	"github.com/ofavor/socket-gw/transport"
)

// Options for client
type Options struct {
	Addr      string
	Transport transport.Transport
}

// Option function to set client options
type Option func(opts *Options)

func defaultOptions() Options {
	return Options{
		Addr:      "127.0.0.1:9999",
		Transport: transport.NewTransport(),
	}
}

// LogLevel set log level
func LogLevel(lv string) Option {
	return func(opts *Options) {
		log.SetLevel(lv)
	}
}

// Address set address
func Address(addr string) Option {
	return func(opts *Options) {
		opts.Addr = addr
	}
}

// Transport set transport
func Transport(t transport.Transport) Option {
	return func(opts *Options) {
		opts.Transport = t
	}
}
