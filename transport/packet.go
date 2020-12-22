package transport

// PacketType packet type
type PacketType uint8

const (
	// PacketTypeAuth auth
	PacketTypeAuth PacketType = 0

	// PacketTypeAuthAck auth ACK
	PacketTypeAuthAck PacketType = 1

	// PacketTypePing ping
	PacketTypePing PacketType = 2

	// PacketTypePong pong
	PacketTypePong PacketType = 3

	// PacketTypeCustom custom packet type
	PacketTypeCustom PacketType = 10
)

// Packet data packet
type Packet struct {
	Length uint32

	Type PacketType

	Body []byte
}

// NewPacket create new packet
func NewPacket(t PacketType, body []byte) *Packet {
	return &Packet{
		Length: uint32(len(body)),
		Type:   t,
		Body:   body,
	}
}
