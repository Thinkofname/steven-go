package protocol

import (
	"errors"
	"fmt"
	"io"
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
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
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

// Position is a location in the world packed into a 64 bit integer
type Position uint64

func NewPosition(x, y, z int) Position {
	return ((Position(x) & 0x3FFFFFF) << 38) |
		((Position(y) & 0xFFF) << 26) |
		(Position(z) & 0x3FFFFFF)
}

// X returns the X component of the position
func (p Position) X() int {
	return int(int64(p) >> 38)
}

// Y returns the Y component of the position
func (p Position) Y() int {
	return int((int64(p) >> 26) & 0xFFF)
}

// Z returns the Z component of the position
func (p Position) Z() int {
	return int(int64(p) << 38 >> 38)
}

// String returns a string representation of the position
func (p Position) String() string {
	return fmt.Sprintf("%d,%d,%d", p.X(), p.Y(), p.Z())
}

// UUID is an unique identifier
type UUID [16]byte

func (u *UUID) Serialize(w io.Writer) error {
	_, err := w.Write(u[:])
	return err
}
func (u *UUID) Deserialize(r io.Reader) error {
	_, err := io.ReadFull(r, u[:])
	return err
}

// Packet is a structure that can be serialized or deserialized from
// Minecraft connection
type Packet interface {
	write(io.Writer) error
	read(io.Reader) error
	id() int
}
