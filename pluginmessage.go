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
	"bytes"
	"io"
	"log"
	"reflect"

	"github.com/thinkofdeath/steven/protocol"
)

var (
	pluginMessagesServerbound = map[string]reflect.Type{}
	pluginMessagesClientbound = map[string]reflect.Type{}
)

func registerPluginMessage(pm pluginMessage, serverbound bool) {
	t := reflect.TypeOf(pm).Elem()
	if serverbound {
		pluginMessagesServerbound[pm.channel()] = t
	} else {
		pluginMessagesClientbound[pm.channel()] = t
	}
}

func (h handler) handlePluginMessage(channel string, r io.Reader, serverbound bool) {
	var pm reflect.Type
	var ok bool
	if serverbound {
		pm, ok = pluginMessagesServerbound[channel]
	} else {
		pm, ok = pluginMessagesClientbound[channel]
	}
	if !ok {
		log.Printf("Unhandled plugin message %s\n", channel)
		return
	}
	p := reflect.New(pm).Interface().(pluginMessage)
	err := p.read(r)
	if err != nil {
		log.Printf("Failed to handle plugin message %s: %s", channel, err)
		return
	}
	h.Handle(p)
}

func sendPluginMessage(pm pluginMessage) {
	var buf bytes.Buffer
	pm.write(&buf)
	Client.network.Write(&protocol.PluginMessageServerbound{
		Channel: pm.channel(),
		Data:    buf.Bytes(),
	})
}

//go:generate protocol_builder $GOFILE

type pluginMessage interface {
	write(io.Writer) error
	read(io.Reader) error
	channel() string
}

// This is a packet
type pmMinecraftBrand struct {
	Brand string
}

func (*pmMinecraftBrand) channel() string {
	return "MC|Brand"
}

func init() {
	registerPluginMessage((*pmMinecraftBrand)(nil), true)
	registerPluginMessage((*pmMinecraftBrand)(nil), false)
}
