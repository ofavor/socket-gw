package transport

import (
	"encoding/binary"
	"io"

	"github.com/ofavor/socket-gw/internal/log"
)

type defaultTransport struct {
}

type packet struct {
	length int
	typ    PacketType
	body   []byte
}

func (p *packet) Length() int {
	return p.length
}

func (p *packet) Type() PacketType {
	return p.typ
}

func (p *packet) Body() []byte {
	return p.body
}

func (t *defaultTransport) Read(conn Conn) (*Packet, error) {
	l := uint32(0)
	if err := binary.Read(conn, binary.LittleEndian, &l); err != nil {
		return nil, err
	}
	ty := PacketType(0)
	if err := binary.Read(conn, binary.LittleEndian, &ty); err != nil {
		return nil, err
	}
	body := make([]byte, l)
	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, err
	}
	return &Packet{
		Length: l,
		Type:   ty,
		Body:   body,
	}, nil
}

func (t *defaultTransport) Write(conn Conn, p *Packet) error {
	log.Debug("Transport write packet:", p)
	if err := binary.Write(conn, binary.LittleEndian, p.Length); err != nil {
		return err
	}
	if err := binary.Write(conn, binary.LittleEndian, p.Type); err != nil {
		return err
	}
	body := p.Body
	if body != nil {
		total := 0
		for total < int(p.Length) {
			n, err := conn.Write(body[total:])
			if err != nil {
				return err
			}
			total += n
		}
	}
	return nil
}

func newDefaultTransport() Transport {
	return &defaultTransport{}
}
