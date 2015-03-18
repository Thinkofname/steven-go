//go:generate protocol_builder $GOFILE Status serverbound

package protocol

// StatusRequest is sent by the client instantly after
// switching to the Status protocol state and is used
// to signal the server to send a StatusResponse to the
// client
//
// Currently the packet id is: 0x00
type StatusRequest struct {
}

// StatusPing is sent by the client after recieving a
// StatusResponse. The client uses the time from sending
// the ping until the time of recieving a pong to measure
// the latency between the client and the server.
//
// Currently the packet id is: 0x01
type StatusPing struct {
	// The time when the ping was sent
	Time int64
}
