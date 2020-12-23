package gw

// Gateway interface
type Gateway interface {
	// Run start gateway
	Run() error
	Stop() error
}

type gateway struct {
	opts Options
}

func (g *gateway) Run() error {
	return g.opts.Server.Run()
}

func (g *gateway) Stop() error {
	return g.opts.Server.Stop()
}

// NewGateway create new gateway
func NewGateway(opts ...Option) Gateway {
	options := defaultOptions()
	for _, o := range opts {
		o(&options)
	}

	return &gateway{
		opts: options,
	}
}
