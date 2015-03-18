//go:generate protocol_builder $GOFILE Login clientbound

package protocol

// LoginDisconnect is sent by the server if there was any issues
// authenticating the player during login or the general server
// issues (e.g. too many players)
//
// Currently the packet id is: 0x00
type LoginDisconnect struct {
	// JSON string
	Reason string
}

// Currently the packet id is: 0x01
type EncryptionRequest struct {
	// Generally empty, left in from legacy auth
	// but is still used by the client if provided
	ServerID    string
	PublicKey   []byte `VarInt`
	VerifyToken []byte `VarInt`

	Test []int32 `byte`
}
