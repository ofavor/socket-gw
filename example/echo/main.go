package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ofavor/socket-gw/session"
	"github.com/ofavor/socket-gw/transport"

	gw "github.com/ofavor/socket-gw"
	"github.com/ofavor/socket-gw/internal/log"
)

type myHandler struct {
	session.BaseHandler
}

func (h *myHandler) OnSessionAuth(s *session.Session, p *transport.Packet) error {
	token := string(p.Body)
	if token == "abcd" {
		return nil
	}
	return errors.New("Token invalid")
}

func (h *myHandler) OnSessionReceived(s *session.Session, p *transport.Packet) error {
	switch p.Type {
	case 11:
		s.Send(p)
	}
	return nil
}

func main() {
	container := &myHandler{
		BaseHandler: session.NewBaseHandler(),
	}
	gw := gw.NewGateway(
		gw.LogLevel("debug"),
		gw.Address(":9999"),
		gw.SessionAuth(true),
		gw.SessionHandler(container),
	)
	go func() {
		for range time.Tick(5 * time.Second) {
			log.Infof("Current session count:%d", container.GetTotalSessions())
		}
	}()
	if err := gw.Run(); err != nil {
		log.Fatal("Gateway run error:", err)
	}

	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	select {
	case <-sc:
	}
	gw.Stop()
}
