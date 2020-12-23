package client

import (
	"net"

	"github.com/ofavor/socket-gw/transport"
)

type tcpClient struct {
	conn transport.Conn
	opts Options
}

func (c *tcpClient) Connect() error {
	conn, err := net.Dial("tcp", c.opts.Addr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *tcpClient) Send(p *transport.Packet) error {
	return c.opts.Transport.Write(c.conn, p)
}

func (c *tcpClient) Recv() (*transport.Packet, error) {
	return c.opts.Transport.Read(c.conn)
}

func (c *tcpClient) Close() error {
	return c.conn.Close()
}

func newTCPClient(opts ...Option) Client {
	options := defaultOptions()
	for _, o := range opts {
		o(&options)
	}
	return &tcpClient{
		opts: options,
	}
}
