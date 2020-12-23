package session

import (
	"sync"

	"github.com/ofavor/socket-gw/internal/log"
	"github.com/ofavor/socket-gw/transport"
)

// Handler session handler
type Handler interface {
	OnSessionAuth(*Session, *transport.Packet) error

	OnSessionEstablished(*Session) error

	OnSessionReceived(*Session, *transport.Packet) error

	OnSessionClosed(*Session) error
}

// BaseHandler basic session handler
type BaseHandler struct {
	sync.RWMutex
	sessions map[string]*Session
}

// NewHandler create new base session handler
func NewHandler() Handler {
	return &BaseHandler{
		sessions: make(map[string]*Session),
	}
}

// NewBaseHandler create new base session handler
func NewBaseHandler() BaseHandler {
	return BaseHandler{
		sessions: make(map[string]*Session),
	}
}

// GetSession get session by id
func (h *BaseHandler) GetSession(id string) (*Session, error) {
	h.RLock()
	defer h.RUnlock()
	s, ok := h.sessions[id]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return s, nil
}

// GetTotalSessions get total session count
func (h *BaseHandler) GetTotalSessions() int {
	h.RLock()
	defer h.RUnlock()
	return len(h.sessions)
}

// OnSessionAuth called when session request to auth
func (h *BaseHandler) OnSessionAuth(s *Session, p *transport.Packet) error {
	return nil
}

// OnSessionEstablished called when session is established
func (h *BaseHandler) OnSessionEstablished(s *Session) error {
	log.Debugf("Session %s was established", s.ID())
	h.Lock()
	h.sessions[s.ID()] = s
	h.Unlock()
	return nil
}

// OnSessionReceived called when session received data
func (h *BaseHandler) OnSessionReceived(s *Session, p *transport.Packet) error {
	log.Debugf("Session %s got data packet:%v", s.ID(), p)
	return nil
}

// OnSessionClosed called when session is closed
func (h *BaseHandler) OnSessionClosed(s *Session) error {
	log.Debugf("Session %s was closed", s.ID())
	h.Lock()
	delete(h.sessions, s.ID())
	h.Unlock()
	return nil
}
