package protocol

import (
	"bytes"
	"errors"
)

var (
	// ErrUnexpectedPacket is returned if the method got a different
	// packet to the one it expected.
	ErrUnexpectedPacket = errors.New("unexpected packet")
)

// Serializable is a type which can be serialized into a packet.
// This is used by protocol_builder when the struct tag 'as' is set
// to "raw".
type Serializable interface {
	Serialize(w *bytes.Buffer) error
	Deserialize(r *bytes.Reader) error
}

// State defined which state the protocol is in.
type State int

// States of the protocol.
// Handshaking is default.
const (
	Handshaking State = 0
	Play        State = 1
	Status      State = 2
	Login       State = 3
)

const (
	// SupportedProtocolVersion is current protocol version this package defines
	SupportedProtocolVersion = 47
)

const (
	clientbound = iota
	serverbound
)
const maxPacketCount = 100

var packetCreator [4][2][maxPacketCount]func() Packet

// VarInt is a variable length integer with a cap of
// 32 bits
type VarInt int32

// VarLong is a variable length integer with a cap of
// 64 bits
type VarLong int64

// Packet is a structure that can be serialized or deserialized from
// Minecraft connection
type Packet interface {
	write(*bytes.Buffer) error
	read(*bytes.Reader) error
	id() int
}
