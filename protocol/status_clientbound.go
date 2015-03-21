//go:generate protocol_builder $GOFILE Status clientbound

package protocol

// StatusResponse is sent as a reply to a StatusRequest.
// The Status should contain a json encoded structure with
// version information, a player sample, a description/MOTD
// and optionally a favicon.
//
// The structure is as follows
//     {
//         "version": {
//             "name": "1.8.3",
//             "protocol": 47,
//         },
//         "players": {
//             "max": 20,
//             "online": 1,
//             "sample": [
//                 {"name": "Thinkofdeath", "id": "4566e69f-c907-48ee-8d71-d7ba5aa00d20"}
//             ]
//         },
//         "description": "Hello world",
//         "favicon": "data:image/png;base64,<data>"
//     }
//
// Currently the packet id is: 0x00
type StatusResponse struct {
	Status StatusReply `as:"json"`
}

// StatusPong is sent as a reply to a StatusPing.
// The Time field should be exactly the same as the
// one sent by the client.
//
// Currently the packet id is: 0x01
type StatusPong struct {
	Time int64
}
