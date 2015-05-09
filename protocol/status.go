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

package protocol

import (
	"fmt"
	"time"

	"github.com/thinkofdeath/phteven/chat"
)

// StatusReply is the reply retrieved from a server when pinging
// it.
type StatusReply struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int            `json:"max"`
		Online int            `json:"online"`
		Sample []StatusPlayer `json:"sample,omitempty"`
	} `json:"players"`
	Description chat.AnyComponent `json:"description"`
	Favicon     string            `json:"favicon"`
}

// StatusPlayer is one of the sample players in a StatusReply
type StatusPlayer struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// RequestStatus starts a status request to the server and
// returns the results of the request. The connection will
// be closed after this request.
func (c *Conn) RequestStatus() (response StatusReply, ping time.Duration, err error) {
	defer c.Close()

	err = c.WritePacket(&Handshake{
		ProtocolVersion: SupportedProtocolVersion,
		Host:            c.host,
		Port:            c.port,
		Next:            VarInt(Status - 1),
	})
	if err != nil {
		return
	}
	c.State = Status
	if err = c.WritePacket(&StatusRequest{}); err != nil {
		return
	}

	// Get the reply
	var packet Packet
	if packet, err = c.ReadPacket(); err != nil {
		return
	}

	resp, ok := packet.(*StatusResponse)
	if !ok {
		err = fmt.Errorf("unexpected packet %#v", packet)
		return
	}
	response = resp.Status

	t := time.Now()
	if err = c.WritePacket(&StatusPing{
		Time: t.UnixNano(),
	}); err != nil {
		return
	}

	// Get the pong reply
	packet, err = c.ReadPacket()
	if err != nil {
		return
	}

	_, ok = packet.(*StatusPong)
	if !ok {
		err = fmt.Errorf("unexpected packet %#v", packet)
	}
	ping = time.Now().Sub(t)
	return
}
