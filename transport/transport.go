package transport

// Conn connection interface
type Conn interface {
	// Read reads data from the connection.
	// Read can be made to time out and return an error after a fixed
	// time limit; see SetDeadline and SetReadDeadline.
	Read(b []byte) (n int, err error)

	// Write writes data to the connection.
	// Write can be made to time out and return an error after a fixed
	// time limit; see SetDeadline and SetWriteDeadline.
	Write(b []byte) (n int, err error)

	// Close closes the connection.
	// Any blocked Read or Write operations will be unblocked and return errors.
	Close() error
}

// Transport interface for reading/writing data packet
type Transport interface {
	// Read data packet
	Read(Conn) (*Packet, error)

	// Write data packet
	Write(Conn, *Packet) error
}

// NewTransport create new transport
func NewTransport() Transport {
	return newDefaultTransport()
}
