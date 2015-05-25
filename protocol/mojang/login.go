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
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	loginURL    = "https://authserver.mojang.com/authenticate"
	refreshURL  = "https://authserver.mojang.com/refresh"
	validateURL = "https://authserver.mojang.com/validate"
)

type loginRequest struct {
	Agent struct {
		Name    string `json:"name"`
		Version int    `json:"version"`
	} `json:"agent"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ClientToken string `json:"clientToken"`
}

type loginReply struct {
	AccessToken     string `json:"accessToken"`
	SelectedProfile struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
}

// Login tries to login using the passed username (or email) and password
// and returns the complete profile. error is non-nil if the login
// fails.
func Login(username, password, token string) (Profile, error) {
	req := loginRequest{
		Username:    username,
		Password:    password,
		ClientToken: token,
	}
	req.Agent.Name = "Minecraft"
	req.Agent.Version = 1
	b, err := json.Marshal(req)
	if err != nil {
		return Profile{}, err
	}
	r := bytes.NewReader(b)
	resp, err := http.Post(loginURL, "application/json", r)
	if err != nil {
		return Profile{}, err
	}
	defer resp.Body.Close()

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Profile{}, err
	}
	var me Error
	err = json.Unmarshal(reply, &me)
	if err == nil && me.Type != "" {
		return Profile{}, me
	}
	var lr loginReply
	err = json.Unmarshal(reply, &lr)
	return Profile{
		AccessToken: lr.AccessToken,
		Username:    lr.SelectedProfile.Name,
		ID:          lr.SelectedProfile.ID,
	}, err
}

type refreshRequest struct {
	AccessToken string `json:"accessToken"`
	ClientToken string `json:"clientToken"`
}

// Refresh attempts to refresh the passed profile's accessToken
// for futher use. The passed token should be the same as the
// one passed to Login.
func Refresh(profile Profile, token string) (Profile, error) {
	req := refreshRequest{
		AccessToken: profile.AccessToken,
		ClientToken: token,
	}
	b, err := json.Marshal(req)
	if err != nil {
		return Profile{}, err
	}
	// Try to reuse old token
	r := bytes.NewReader(b)
	resp, err := http.Post(validateURL, "application/json", r)
	if err == nil {
		defer resp.Body.Close()
		reply, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			var me Error
			err = json.Unmarshal(reply, &me)
			if err != nil || me.Type == "" {
				return profile, nil
			}
		}
	}
	r = bytes.NewReader(b)

	// Try and get a updated one
	resp, err = http.Post(refreshURL, "application/json", r)
	if err != nil {
		return Profile{}, err
	}
	defer resp.Body.Close()

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Profile{}, err
	}
	var me Error
	err = json.Unmarshal(reply, &me)
	if err == nil && me.Type != "" {
		return Profile{}, me
	}
	var lr loginReply
	err = json.Unmarshal(reply, &lr)
	return Profile{
		AccessToken: lr.AccessToken,
		Username:    lr.SelectedProfile.Name,
		ID:          lr.SelectedProfile.ID,
	}, err
}
