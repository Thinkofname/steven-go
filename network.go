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

package steven

import (
	"fmt"

	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/protocol/mojang"
)

var (
	// TODO(Think) Tweek the values
	writeChan = make(chan protocol.Packet, 200)
	readChan  = make(chan protocol.Packet, 200)
	errorChan = make(chan error, 1)
	killChan  = make(chan struct{})
	conn      *protocol.Conn
)

func startConnection(profile mojang.Profile, server string) {
	var err error
	conn, err = protocol.Dial(server)
	if err != nil {
		closeWithError(err)
		return
	}

	err = conn.LoginToServer(profile)
	if err != nil {
		closeWithError(err)
		return
	}

	defer fmt.Println("Read handler closed")
preLogin:
	for {
		packet, err := conn.ReadPacket()
		if err != nil {
			closeWithError(err)
			return
		}
		switch packet := packet.(type) {
		case *protocol.SetInitialCompression:
			conn.SetCompression(int(packet.Threshold))
		case *protocol.LoginSuccess:
			conn.State = protocol.Play
			break preLogin
		default:
			panic(fmt.Errorf("unhandled packet %T", packet))
		}
	}

	go writeHandler(conn)

	for {
		packet, err := conn.ReadPacket()
		if err != nil {
			closeWithError(err)
			return
		}

		// Handle keep alives async as there is no need to process them
		switch packet := packet.(type) {
		case *protocol.KeepAliveClientbound:
			writeChan <- &protocol.KeepAliveServerbound{ID: packet.ID}
		case *protocol.SetCompression:
			conn.SetCompression(int(packet.Threshold))
		default:
			readChan <- packet
		}
	}
}

func closeWithError(err error) {
	// Try to save the error if one isn't already there
	select {
	case errorChan <- err:
	default:
	}
}

func writeHandler(conn *protocol.Conn) {
	defer fmt.Println("Write handler closed")
	for {
		select {
		case packet := <-writeChan:
			err := conn.WritePacket(packet)
			if err != nil {
				closeWithError(err)
				<-killChan
				return
			}
		case <-killChan:
			return
		}
	}
}
