//go:generate protocol_builder $GOFILE Play clientbound

package protocol

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
