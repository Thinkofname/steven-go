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

type joinData struct {
	AccessToken     string `json:"accessToken"`
	SelectedProfile string `json:"selectedProfile"`
	ServerID        string `json:"serverId"`
}

type Error struct {
	Message string `json:"errorMessage"`
	Type    string `json:"error"`
}

func (m Error) Error() string {
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
		var e Error
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
