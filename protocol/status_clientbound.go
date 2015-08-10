// Copyright 2015 Matthew Collins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
// This is a Minecraft packet
type StatusResponse struct {
	Status StatusReply `as:"json"`
}

// StatusPong is sent as a reply to a StatusPing.
// The Time field should be exactly the same as the
// one sent by the client.
//
// This is a Minecraft packet
type StatusPong struct {
	Time int64
}
