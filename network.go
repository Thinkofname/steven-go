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

type networkManager struct {
	conn      *protocol.Conn
	writeChan chan protocol.Packet
	readChan  chan protocol.Packet
	errorChan chan error
	closeChan chan struct{}
}

func (n *networkManager) init() {
	n.writeChan = make(chan protocol.Packet, 200)
	n.readChan = make(chan protocol.Packet, 200)
	n.errorChan = make(chan error, 1)
	n.closeChan = make(chan struct{}, 1)
}

func (n *networkManager) Connect(profile mojang.Profile, server string) {
	go func() {
		var err error
		n.conn, err = protocol.Dial(server)
		if err != nil {
			n.SignalClose(err)
			return
		}

		err = n.conn.LoginToServer(profile)
		if err != nil {
			n.SignalClose(err)
			return
		}

	preLogin:
		for {
			packet, err := n.conn.ReadPacket()
			if err != nil {
				n.SignalClose(err)
				return
			}
			switch packet := packet.(type) {
			case *protocol.SetInitialCompression:
				n.conn.SetCompression(int(packet.Threshold))
			case *protocol.LoginSuccess:
				n.conn.State = protocol.Play
				break preLogin
			default:
				n.SignalClose(fmt.Errorf("unhandled packet %T", packet))
			}
		}

		first := true
		for {
			packet, err := n.conn.ReadPacket()
			if err != nil {
				n.SignalClose(err)
				return
			}
			if first {
				go n.writeHandler()
				first = false
			}

			// Handle keep alives async as there is no need to process them
			switch packet := packet.(type) {
			case *protocol.KeepAliveClientbound:
				n.Write(&protocol.KeepAliveServerbound{ID: packet.ID})
			case *protocol.SetCompression:
				n.conn.SetCompression(int(packet.Threshold))
			default:
				n.readChan <- packet
			}
		}
	}()
}

func (n *networkManager) writeHandler() {
	for packet := range n.writeChan {
		err := n.conn.WritePacket(packet)
		if err != nil {
			n.SignalClose(err)
			return
		}
	}
}

func (n *networkManager) SignalClose(err error) {
	// Try to save the error if one isn't already there
	select {
	case n.errorChan <- err:
	default:
	}
}

func (n *networkManager) Error() <-chan error {
	return n.errorChan
}

func (n *networkManager) Read() <-chan protocol.Packet {
	return n.readChan
}

func (n *networkManager) Write(packet protocol.Packet) {
	select {
	case n.writeChan <- packet:
	case <-n.closeChan:
		n.closeChan <- struct{}{} // Keep the closed state
		return
	}
}

func (n *networkManager) Close() {
	n.closeChan <- struct{}{}
	n.conn.Close()
}
