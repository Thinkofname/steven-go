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

//go:generate protocol_builder $GOFILE Handshaking serverbound

package protocol

// Handshake is the first packet sent in the protocol.
// Its used for deciding if the request is a client
// is requesting status information about the server
// (MOTD, players etc) or trying to login to the server.
//
// The host and port fields are not used by the vanilla
// server but are there for virtual server hosting to
// be able to redirect a client to a target server with
// a single address + port.
//
// Some modified servers/proxies use the handshake field
// differently, packing information into the field other
// than the hostname due to the protocol not providing
// any system for custom information to be transfered
// by the client to the server until after login.
//
// This is a Minecraft packet
type Handshake struct {
	// The protocol version of the connecting client
	ProtocolVersion VarInt
	// The hostname the client connected to
	Host string
	// The port the client connected to
	Port uint16
	// The next protocol state the client wants
	Next VarInt
}
