package server

import (
	"net"

	"github.com/ofavor/socket-gw/internal/log"
	"github.com/ofavor/socket-gw/session"
)

type tcpServer struct {
	opts Options
}

func (s *tcpServer) Init(o Option) {
	o(&s.opts)
}

func (s *tcpServer) Options() Options {
	return s.opts
}

func (s *tcpServer) Run() error {
	l, err := net.Listen("tcp", s.opts.Addr)
	if err != nil {
		return err
	}
	log.Info("Server is listened on:", s.opts.Addr)
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Error("Accept connection error:", err)
				continue
			}
			log.Info("Got connection:", conn.RemoteAddr())
			session := session.NewSession(conn, s.opts.Transport, s.opts.SessionHandler, s.opts.SessionAuth)
			go session.Run()
		}
	}()
	return nil
}

func (s *tcpServer) Stop() error {
	return nil
}

func newTCPServer(opts ...Option) Server {
	options := defaultOptions()
	for _, o := range opts {
		o(&options)
	}
	return &tcpServer{
		opts: options,
	}
}
