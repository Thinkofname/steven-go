//go:generate protocol_builder $GOFILE Play clientbound

package protocol

import (
	"github.com/thinkofdeath/steven/chat"
)

// KeepAliveClientbound is sent by a server to check if the
// client is still responding and keep the connection open.
// The client should reply with the KeepAliveServerbound
// packet setting ID to the same as this one.
//
// Currently the packet id is: 0x00
type KeepAliveClientbound struct {
	ID VarInt
}

// JoinGame is sent after completing the login process. This
// sets the initial state for the client.
//
// Currently the packet id is: 0x01
type JoinGame struct {
	// The entity id the client will be referenced by
	EntityID int32
	// The starting gamemode of the client
	Gamemode byte
	// The dimension the client is starting in
	Dimension int8
	// The difficuilty setting for the server
	Difficulty byte
	// The max number of players on the server
	MaxPlayers byte
	// The level type of the server
	LevelType string
	// Whether the client should reduce the amount of debug
	// information it displays in F3 mode
	ReducedDebugInfo bool
}

// ServerMessage is a message sent by the server. It could be from a player
// or just a system message. The Type field controls the location the
// message is displayed at and when the message is displayed.
//
// Currently the packet id is: 0x02
type ServerMessage struct {
	Message chat.AnyComponent `as:"json"`
	// 0 - Chat message, 1 - System message, 2 - Action bar message
	Type byte
}

// TimeUpdate is sent to sync the world's time to the client, the client
// will manually tick the time itself so this doesn't need to sent repeatedly
// but if the server or client has issues keeping up this can fall out of sync
// so it is a good idea to sent this now and again
//
// Currently the packet id is: 0x03
type TimeUpdate struct {
	WorldAge  int64
	TimeOfDay int64
}
