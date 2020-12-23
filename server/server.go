package server

// Server interface
type Server interface {
	// Init server with option
	Init(Option)

	// Options get options
	Options() Options

	// Run start server
	Run() error

	// Stop server
	Stop() error
}

// NewServer create new server(TCP)
func NewServer(opts ...Option) Server {
	return newTCPServer(opts...)
}
