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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/thinkofdeath/steven/audio"
	"github.com/thinkofdeath/steven/console"
	"github.com/thinkofdeath/steven/resource"
)

var (
	loadedSounds = map[pluginKey]audio.SoundBuffer{}
	soundList    []audio.Sound
)

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
