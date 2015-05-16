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
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	h := &Handshake{
		ProtocolVersion: 4,
		Host:            "",
		Port:            25565,
		Next:            1,
	}
	buf := &bytes.Buffer{}
	h.write(buf)

	h2 := &Handshake{}
	h2.read(bytes.NewReader(buf.Bytes()))

	if !reflect.DeepEqual(h, h2) {
		t.Fail()
	}
}

// Created by fuzzing
var testData = []string{
	"\xff\xfd\x80\xc0(",
	"00000000000\xc000000000" +
		"00000000000000000000" +
		"00000000000",
	"0d000000000000000000" +
		"00000000000000000000" +
		"00000000000",
	"0 000000000000000000" +
		"00000000000000000000" +
		"00000000000000000000" +
		"0000",
	"000\x800000000000000000" +
		"00000000000000000000" +
		"000000000000",
	"0 000000000000000000" +
		"0000000000\x1e\v00000000" +
		"00000000000000000000" +
		"000",
	"0 0000000000000ee000" +
		"00000000000000000000" +
		"00000000000000000000" +
		"0000",
	"00000000000\xc000000000" +
		"00000000000000000000" +
		"000000000",
	"\x13\x1c0\x80\xc0\xfd\x80\xc0(00000000000",
}

func TestCrashers(t *testing.T) {
	for _, str := range testData {
		c := &Conn{
			r:                    strings.NewReader(str),
			w:                    nil,
			net:                  nil,
			direction:            serverbound,
			State:                Play,
			compressionThreshold: -1,
		}
		c.readPacket()
	}
}
