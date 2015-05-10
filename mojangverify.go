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
	"archive/zip"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

const (
	authlibURL     = "https://libraries.minecraft.net/com/mojang/authlib/%[1]s/authlib-%[1]s-sources.jar"
	authlibVersion = "1.5.21"
	authlibKeyPath = "yggdrasil_session_pubkey.der"
)

var (
	hasAuthlib  bool
	authlibLock sync.Mutex
	authlibKey  *rsa.PublicKey
)

func init() {
	if _, err := os.Stat(authlibKeyPath); os.IsNotExist(err) {
		authlibLock.Lock()
		go func() {
			getAuthlib()
			parseAuthlibKey()
			authlibLock.Unlock()
		}()
	} else {
		parseAuthlibKey()
		hasAuthlib = true
	}
}

func verifySkinSignature(data, sig []byte) error {
	s := sha1.New()
	s.Write(data)
	return rsa.VerifyPKCS1v15(getAuthlibKey(), crypto.SHA1, s.Sum(nil), sig)
}

func getAuthlibKey() *rsa.PublicKey {
	if hasAuthlib {
		return authlibKey
	}
	authlibLock.Lock()
	defer authlibLock.Unlock()
	if !hasAuthlib {
		panic("Invalid state, missing authlib key")
	}
	return authlibKey
}

func parseAuthlibKey() {
	f, err := os.Open(authlibKeyPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	p, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		panic(err)
	}
	authlibKey = p.(*rsa.PublicKey)
}

func getAuthlib() {
	target := fmt.Sprintf("authlib-%[1]s-sources.jar", authlibVersion)
	fmt.Printf("Obtaining authlib %s\n", authlibVersion)
	resp, err := http.Get(fmt.Sprintf(authlibURL, authlibVersion))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	os.MkdirAll("./", 0777)
	f, err := os.Create(target + ".tmp")
	if err != nil {
		panic(err)
	}
	defer os.Remove(target + ".tmp")
	defer f.Close()
	size, err := io.Copy(f, resp.Body)
	if err != nil {
		panic(err)
	}

	f.Seek(0, 0) // Go back to the start
	fr, err := zip.NewReader(f, size)

	for _, f := range fr.File {
		if f.Name != authlibKeyPath {
			continue
		}
		func() {
			path := authlibKeyPath
			w, err := os.Create(path)
			if err != nil {
				panic(err)
			}
			defer w.Close()
			r, err := f.Open()
			if err != nil {
				panic(err)
			}
			defer r.Close()
			_, err = io.Copy(w, r)
			if err != nil {
				panic(err)
			}
		}()
	}
	hasAuthlib = true
}
