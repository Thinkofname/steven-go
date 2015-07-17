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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/thinkofdeath/steven/protocol/mojang"
)

// BUG(Think) LoginToServer doesn't support offline mode. Call it a feature?

// LoginToServer sends the necessary packets to join a server. This
// also authenticates the request with mojang for online mode connections.
// This stops before LoginSuccess (or any other preceding packets).
func (c *Conn) LoginToServer(profile mojang.Profile) (err error) {
	err = c.WritePacket(&Handshake{
		ProtocolVersion: SupportedProtocolVersion,
		Host:            c.host,
		Port:            c.port,
		Next:            VarInt(Login - 1),
	})
	if err != nil {
		return
	}
	c.State = Login
	if err = c.WritePacket(&LoginStart{
		Username: profile.Username,
	}); err != nil {
		return
	}

	var packet Packet
	if packet, err = c.ReadPacket(); err != nil {
		return
	}

	req, err := checkLoginPacket(c, packet)
	if err != nil {
		return err
	}
	var p interface{}
	if p, err = x509.ParsePKIXPublicKey(req.PublicKey); err != nil {
		return
	}
	pub := p.(*rsa.PublicKey)

	key := make([]byte, 16)
	n, err := rand.Read(key)
	if n != 16 || err != nil {
		return errors.New("crypto error")
	}

	sharedKey, err := rsa.EncryptPKCS1v15(rand.Reader, pub, key)
	if err != nil {
		return
	}
	verifyToken, err := rsa.EncryptPKCS1v15(rand.Reader, pub, req.VerifyToken)
	if err != nil {
		return
	}

	err = mojang.JoinServer(profile, []byte(req.ServerID), key, req.PublicKey)
	if err != nil {
		return
	}

	err = c.WritePacket(&EncryptionResponse{
		SharedSecret: sharedKey,
		VerifyToken:  verifyToken,
	})
	if err != nil {
		return
	}

	err = c.EnableEncryption(key)
	return
}

func checkLoginPacket(c *Conn, p Packet) (*EncryptionRequest, error) {
	switch p := p.(type) {
	case *EncryptionRequest:
		return p, nil
	case *LoginDisconnect:
		return nil, errors.New(p.Reason.String())
	case *LoginSuccess:
		return nil, errors.New("server is in offline mode which is currently unsupported")
	case *SetInitialCompression:
		c.SetCompression(int(p.Threshold))
		p2, err := c.ReadPacket()
		if err != nil {
			return nil, err
		}
		return checkLoginPacket(c, p2)
	default:
		return nil, fmt.Errorf("unexpected packet %#v", p)
	}
}
