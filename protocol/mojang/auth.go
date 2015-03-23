package mojang

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const joinURL = "https://sessionserver.mojang.com/session/minecraft/join"

// Profile contains information about the player required
// to connect to a server
type Profile struct {
	Username    string
	ID          string
	AccessToken string
}

type joinData struct {
	AccessToken     string `json:"accessToken"`
	SelectedProfile string `json:"selectedProfile"`
	ServerID        string `json:"serverId"`
}

type mojError struct {
	Message string `json:"errorMessage"`
	Type    string `json:"error"`
}

func (m mojError) Error() string {
	return fmt.Sprintf("%s: %s", m.Type, m.Message)
}

// JoinServer tries to mark the server has joined on mojang's session servers
// using the passed profile and bytes (as the server hash). The hash is normally
// the serverID + secret key + public key.
func JoinServer(profile Profile, serverHash ...[]byte) error {
	h := sha1.New()
	for _, sh := range serverHash {
		h.Write(sh)
	}
	hash := h.Sum(nil)

	// Mojang uses a hex method which allows for
	// negatives so we have to account for that.
	negative := (hash[0] & 0x80) == 0x80
	if negative {
		twosCompliment(hash)
	}
	serverID := hex.EncodeToString(hash)
	serverID = strings.TrimLeft(serverID, "0")
	if negative {
		serverID = "-" + serverID
	}

	b, err := json.Marshal(joinData{
		AccessToken:     profile.AccessToken,
		SelectedProfile: profile.ID,
		ServerID:        serverID,
	})
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	resp, err := http.Post(joinURL, "application/json", r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if len(reply) != 0 {
		var e mojError
		json.Unmarshal(reply, &e)
		return e
	}

	return nil
}

func twosCompliment(p []byte) {
	carry := true
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = ^p[i]
		if carry {
			carry = p[i] == 0xFF
			p[i]++
		}
	}
}
