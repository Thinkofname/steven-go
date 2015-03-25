package protocol

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"

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
	req, ok := packet.(*EncryptionRequest)
	if !ok {
		return ErrUnexpectedPacket
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
