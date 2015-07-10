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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/thinkofdeath/steven/audio"
	"github.com/thinkofdeath/steven/console"
	"github.com/thinkofdeath/steven/resource"
)

const (
	assetsVersion    = "1.8"
	assetIndexURL    = "https://s3.amazonaws.com/Minecraft.Download/indexes/%s.json"
	assetResourceURL = "http://resources.download.minecraft.net/%s"
)

var (
	assets       assetIndex
	loadedSounds = map[pluginKey]audio.SoundBuffer{}
	soundList    []audio.Sound
)

type assetIndex struct {
	Objects map[string]struct {
		Hash string `json:"hash"`
		Size int    `json:"size"`
	} `json:"objects"`
}

func PlaySound(plugin, name string) {
	key := pluginKey{plugin, name}
	sb, ok := loadedSounds[key]
	if !ok {
		f, err := resource.Open(plugin, "sounds/"+name+".ogg")
		if err != nil {
			v, ok := assets.Objects[fmt.Sprintf("%s/sounds/%s.ogg", plugin, name)]
			if !ok {
				console.Text("Missing sound %s", key)
				return
			}
			loc := fmt.Sprintf("./resources/%s", hashPath(v.Hash))
			f, err = os.Open(loc)
			if err != nil {
				console.Text("Missing sound %s", key)
				return
			}
		}
		defer f.Close()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		sb = audio.NewSoundBufferData(data)
		loadedSounds[key] = sb
	}
	for _, s := range soundList {
		if s.Status() == audio.StatStopped {
			s.SetBuffer(sb)
			s.SetVolume(100.0)
			s.Play()
			return
		}
	}
	s := audio.NewSound()
	s.SetBuffer(sb)
	s.Play()
	s.SetVolume(100.0)
	soundList = append(soundList, s)
}

func init() {
	defLocation := "./resources"
	loc := fmt.Sprintf("%s/%s.index", defLocation, assetsVersion)
	_, err := os.Stat(loc)
	if os.IsNotExist(err) {
		getAssetIndex()
	} else {
		f, err := os.Open(fmt.Sprintf("./resources/%s.index", assetsVersion))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		err = json.NewDecoder(f).Decode(&assets)
		if err != nil {
			panic(err)
		}
	}
	go downloadAssets()
}

func getAssetIndex() {
	console.Text("Getting asset index")
	defLocation := "./resources"
	resp, err := http.Get(fmt.Sprintf(assetIndexURL, assetsVersion))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&assets); err != nil {
		panic(err)
	}
	os.MkdirAll("./resources", 0777)
	f, err := os.Create(fmt.Sprintf("%s/%s.index", defLocation, assetsVersion))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(assets)
	console.Text("Got asset index for %s", assetsVersion)
}

func downloadAssets() {
	defLocation := "./resources"
	for _, v := range assets.Objects {
		path := hashPath(v.Hash)
		loc := fmt.Sprintf("%s/%s", defLocation, path)
		_, err := os.Stat(loc)
		if !os.IsNotExist(err) {
			continue
		}
		func() {
			resp, err := http.Get(fmt.Sprintf(assetResourceURL, path))
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			os.MkdirAll(filepath.Dir(loc), 0777)
			f, err := os.Create(loc + ".tmp")
			if err != nil {
				panic(err)
			}
			defer f.Close()
			n, err := io.Copy(f, resp.Body)
			if err != nil {
				panic(err)
			}
			if n != int64(v.Size) {
				panic(fmt.Sprintf("Got: %d, Wanted: %d for %s", n, v.Size, fmt.Sprintf(assetResourceURL, path)))
			}
			console.Text("Downloaded: %s", loc)
		}()
		os.Rename(loc+".tmp", loc)
	}
}
func hashPath(hash string) string {
	return hash[:2] + "/" + hash
}
