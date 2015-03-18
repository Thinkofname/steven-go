package protocol

import (
	"bytes"
	"reflect"
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
