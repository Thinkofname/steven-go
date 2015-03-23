package main

import (
	"encoding/json"
	"fmt"
	"github.com/thinkofdeath/steven/platform"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/protocol/mojang"
	"github.com/thinkofdeath/steven/render"
	"os"
)

func main() {
	var username, uuid, accessToken string
	for i, arg := range os.Args {
		switch arg {
		case "--username":
			username = os.Args[i+1]
		case "--uuid":
			uuid = os.Args[i+1]
		case "--accessToken":
			accessToken = os.Args[i+1]
		}
	}

	conn, err := protocol.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	err = conn.LoginToServer(mojang.Profile{
		Username:    username,
		ID:          uuid,
		AccessToken: accessToken,
	})
	if err != nil {
		panic(err)
	}

preLogin:
	for {
		packet, err := conn.ReadPacket()
		if err != nil {
			panic(err)
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

	for {
		packet, err := conn.ReadPacket()
		if err != nil {
			panic(err)
		}

		switch packet := packet.(type) {
		case *protocol.KeepAliveClientbound:
			conn.WritePacket(&protocol.KeepAliveServerbound{ID: packet.ID})
		case *protocol.ChunkData, *protocol.ChunkDataBulk:
			continue
		case *protocol.TeleportPlayer:
			conn.WritePacket(&protocol.PlayerPositionLook{
				X:        packet.X,
				Y:        packet.Y,
				Z:        packet.Z,
				Yaw:      packet.Yaw,
				Pitch:    packet.Pitch,
				OnGround: true,
			})
		}

		b, _ := json.Marshal(packet)
		fmt.Printf("Got packet: %T%s\n", packet, b)
	}

	platform.Init(platform.Handler{
		Start: start,
		Draw:  draw,
	})
}

func start() {
	render.Start()
}

func draw() {
	render.Draw()
}
