package session

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ofavor/socket-gw/internal/log"
	"github.com/ofavor/socket-gw/transport"
)

var (
	// SendBufferSize size of send buffer
	SendBufferSize = 100

	// RecvBufferSize size of recv buffer
	RecvBufferSize = 100

	// ErrPacketTypeInvalid packet type invalid
	ErrPacketTypeInvalid = errors.New("packet type was invalid")

	// ErrSessionAuthFailed session auth failed
	ErrSessionAuthFailed = errors.New("session auth failed")

	// ErrSessionAuthTimeout session auth timeout
	ErrSessionAuthTimeout = errors.New("session auth timeout")

	// ErrSessionNotFound session not found error
	ErrSessionNotFound = errors.New("session not found")
)

// Status session status
type Status int

const (
	// StatusInitial initial
	StatusInitial Status = iota

	// StatusReady ready
	StatusReady

	// StatusClosed closed
	StatusClosed
)

// Session represent client connection
type Session struct {
	sync.RWMutex
	id        string
	conn      transport.Conn
	handler   Handler
	transport transport.Transport
	needAuth  bool

	sendCh chan *transport.Packet
	recvCh chan *transport.Packet
	stopCh chan interface{}
	status Status
	meta   map[string]string
}

// NewSession create new session
func NewSession(conn transport.Conn, t transport.Transport, handler Handler, auth bool) *Session {
	return &Session{
		id:        uuid.New().String(),
		conn:      conn,
		transport: t,
		handler:   handler,
		needAuth:  auth,
		sendCh:    make(chan *transport.Packet, SendBufferSize),
		recvCh:    make(chan *transport.Packet, RecvBufferSize),
		stopCh:    make(chan interface{}),
		status:    StatusInitial,
		meta:      make(map[string]string),
	}
}

// ID get session id
func (s *Session) ID() string {
	return s.id
}

func (s *Session) auth() error {
	ch := make(chan error)
	go func() {
		if p, err := s.transport.Read(s.conn); err != nil {
			ch <- err
		} else if p.Type != transport.PacketTypeAuth {
			log.Warn("Session %s first packet is not an auth packet", s.id)
			ch <- ErrSessionAuthFailed
		} else {
			err := s.handler.OnSessionAuth(s, p)
			if err != nil {
				ch <- err
				return
			}
			p = transport.NewPacket(transport.PacketTypeAuthAck, nil)
			if err := s.transport.Write(s.conn, p); err != nil {
				log.Errorf("Session %s send auth ACK error: %s", s.id, err)
				ch <- ErrSessionAuthFailed
				return
			}
			ch <- nil
		}
	}()
	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()
	select {
	case err := <-ch:
		return err
	case <-timeout.C:
		return ErrSessionAuthTimeout
	}
}

// Meta get session meta
func (s *Session) Meta() map[string]string {
	return s.meta
}

// Run start session
func (s *Session) Run() error {
	if s.needAuth {
		if err := s.auth(); err != nil {
			log.Errorf("Session %s auth failed: %s", s.id, err)
			s.conn.Close()
			s.status = StatusClosed
			return err
		}
	}
	s.Lock()
	s.status = StatusReady
	s.Unlock()
	s.handler.OnSessionEstablished(s)
	go func() { // send
		for p := range s.sendCh {
			if err := s.transport.Write(s.conn, p); err != nil {
				log.Errorf("Session %s send error:%s", s.id, err)
				return
			}
		}
	}()
	go func() { // recv
		for {
			p, err := s.transport.Read(s.conn)
			if err != nil {
				log.Errorf("Session %s recv error:%s", s.id, err)
				s.Close()
				return
			}
			log.Debugf("Session %s recv packet:%v", s.id, p)
			s.recvCh <- p
		}
	}()
LOOP:
	for {
		select {
		case p := <-s.recvCh:
			if p.Type > transport.PacketTypeCustom {
				if err := s.handler.OnSessionReceived(s, p); err != nil {
					log.Errorf("Session %s handle received packet error:%s", s.id, err)
				}
				continue
			}
			switch p.Type {
			case transport.PacketTypePing:
				pp := transport.NewPacket(transport.PacketTypePong, nil)
				s.transport.Write(s.conn, pp)
			default:
				log.Errorf("Session %s got packet with invalid type:%d", s.id, p.Type)
			}
		case <-s.stopCh:
			break LOOP
		}
	}

	if err := s.conn.Close(); err != nil {
		log.Errorf("Session %s connection close error:%s", s.id, err)
	}
	close(s.sendCh)
	close(s.recvCh)
	s.handler.OnSessionClosed(s)
	return nil
}

// Send data to client
func (s *Session) Send(p *transport.Packet) error {
	if p.Type <= transport.PacketTypeCustom {
		return ErrPacketTypeInvalid
	}
	s.RLock()
	if s.status == StatusClosed {
		s.RUnlock()
		return nil
	}
	s.RUnlock()
	s.sendCh <- p
	return nil
}

// Close session
func (s *Session) Close() error {
	s.RLock()
	if s.status == StatusClosed {
		s.RUnlock()
		return nil
	}
	s.RUnlock()
	s.stopCh <- 1
	return nil
}
