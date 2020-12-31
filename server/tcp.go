package server

import (
	"net"
	"sync"

	"github.com/ofavor/socket-gw/internal/log"
	"github.com/ofavor/socket-gw/session"
)

type tcpServer struct {
	opts     Options
	listener net.Listener
	stopCh   chan interface{}
	wg       sync.WaitGroup
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
	s.listener = l
	log.Info("Server is listened on:", s.opts.Addr)
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.stopCh: // server stopped
					return
				default:
				}
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
	close(s.stopCh) // close the stop channel
	s.listener.Close()
	s.wg.Wait()
	return nil
}

func newTCPServer(opts ...Option) Server {
	options := defaultOptions()
	for _, o := range opts {
		o(&options)
	}
	return &tcpServer{
		opts:   options,
		stopCh: make(chan interface{}),
	}
}
