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
	"math/rand"
	"os"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/audio"
	"github.com/thinkofdeath/steven/console"
	"github.com/thinkofdeath/steven/resource"
)

var (
	loadedSounds = map[pluginKey]audio.SoundBuffer{}
	soundList    []audio.Sound
	soundInfo    = map[string]soundData{}
	soundRandom  = rand.New(rand.NewSource(time.Now().Unix()))
)

type soundCategory string

const (
	sndAmbient soundCategory = "ambient"
	sndWeather soundCategory = "weather"
	sndPlayer  soundCategory = "player"
	sndNeutral soundCategory = "neutral"
	sndHostile soundCategory = "hostile"
	sndBlock   soundCategory = "block"
	sndMusic   soundCategory = "music"
)

func PlaySoundAt(name string, vol, pitch float64, v mgl32.Vec3) {
	snd, ok := soundInfo[name]
	if !ok {
		return
	}
	playSoundInternal(snd.Sounds[soundRandom.Intn(len(snd.Sounds))], vol, pitch, true, v)
}

func PlaySound(name string) {
	snd, ok := soundInfo[name]
	if !ok {
		return
	}
	playSoundInternal(snd.Sounds[soundRandom.Intn(len(snd.Sounds))], 1, 1, false, mgl32.Vec3{})
}

func playSoundInternal(snd sound, vol, pitch float64, rel bool, pos mgl32.Vec3) {
	name := snd.Name
	key := pluginKey{"minecraft", name}
	sb, ok := loadedSounds[key]
	if !ok {
		f, err := resource.Open("minecraft", "sounds/"+name+".ogg")
		if err != nil {
			v, ok := assets.Objects[fmt.Sprintf("minecraft/sounds/%s.ogg", name)]
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
	var s audio.Sound
	n := true
	for _, sn := range soundList {
		if sn.Status() == audio.StatStopped {
			s = sn
			n = false
			break
		}
	}
	if n {
		s = audio.NewSound()
		soundList = append(soundList, s)
	}
	s.SetBuffer(sb)
	s.Play()
	s.SetVolume(snd.Volume * vol * 100.0)
	s.SetMinDistance(5)
	s.SetAttenuation(0.008)
	s.SetPitch(pitch)
	if rel {
		s.SetRelative(true)
		s.SetPosition(pos.X(), pos.Y(), pos.Z())
	} else {
		s.SetRelative(false)
	}
}

type soundData struct {
	Category soundCategory
	Sounds   []sound
}

type sound struct {
	Name   string
	Stream bool
	Volume float64
}

func (s *sound) UnmarshalJSON(b []byte) error {
	type js struct {
		Name   string
		Stream bool
		Volume *float64
	}
	s.Volume = 1
	if b[0] == '"' {
		return json.Unmarshal(b, &s.Name)
	}
	var val js
	err := json.Unmarshal(b, &val)
	s.Name = val.Name
	s.Stream = val.Stream
	if val.Volume != nil {
		s.Volume = *val.Volume
	}
	return err
}

func loadSoundData() {
	v := assets.Objects["minecraft/sounds.json"]
	loc := fmt.Sprintf("./resources/%s", hashPath(v.Hash))
	f, err := os.Open(loc)
	if err != nil {
		panic(err)
	}
	snds, _ := resource.OpenAll("minecraft", "sounds.json")
	snds = append([]io.ReadCloser{f}, snds...)
	for _, s := range snds {
		func() {
			defer s.Close()
			data := map[string]soundData{}
			err := json.NewDecoder(s).Decode(&data)
			if err != nil {
				panic(err)
			}
			for k, v := range data {
				soundInfo[k] = v
			}
		}()
	}

}
