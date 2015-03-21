//go:generate protocol_builder $GOFILE Login clientbound

package protocol

import (
	"github.com/thinkofdeath/steven/chat"
)

// LoginDisconnect is sent by the server if there was any issues
// authenticating the player during login or the general server
// issues (e.g. too many players).
//
// Currently the packet id is: 0x00
type LoginDisconnect struct {
	Reason chat.AnyComponent `as:"json"`
}

// EncryptionRequest is sent by the server if the server is in
// online mode. If it is not sent then its assumed the server is
// in offline mode.
//
// Currently the packet id is: 0x01
type EncryptionRequest struct {
	// Generally empty, left in from legacy auth
	// but is still used by the client if provided
	ServerID string
	// A RSA Public key serialized in x.509 PRIX format
	PublicKey []byte `length:"VarInt"`
	// Token used by the server to verify encryption is working
	// correctly
	VerifyToken []byte `length:"VarInt"`
}

// LoginSuccess is sent by the server if the player successfully
// authenicates with the session servers (online mode) or straight
// after LoginStart (offline mode).
//
// Currently the packet id is: 0x02
type LoginSuccess struct {
	// String encoding of a uuid (with hyphens)
	UUID     string
	Username string
}

// SetInitialCompression sets the compression threshold during the
// login state.
//
// Currently the packet id is: 0x03
type SetInitialCompression struct {
	// Threshold where a packet should be sent compressed
	Threshold VarInt
}
