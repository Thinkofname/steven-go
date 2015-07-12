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

package audio

import (
	"io"
	"io/ioutil"
	"os"
	"unsafe"
)

// #cgo LDFLAGS: -lcsfml-audio
// #include <SFML/Audio/SoundBuffer.h>
// #include <SFML/Audio/Sound.h>
// #include <SFML/Audio/Listener.h>
// #include <SFML/Audio/Music.h>
// #include <SFML/System/Vector3.h>
// #include <SFML/System/InputStream.h>
// #include <stdlib.h>
import "C"

type SoundBuffer struct {
	internal *C.sfSoundBuffer
}

func NewSoundBuffer(file string) SoundBuffer {
	str := C.CString(file)
	defer C.free(unsafe.Pointer(str))
	return SoundBuffer{C.sfSoundBuffer_createFromFile(str)}
}

func NewSoundBufferData(data []byte) SoundBuffer {
	return SoundBuffer{C.sfSoundBuffer_createFromMemory(unsafe.Pointer(&data[0]), (C.size_t)(len(data)))}
}

func (sb SoundBuffer) Free() {
	C.sfSoundBuffer_destroy(sb.internal)
}

type Sound struct {
	internal *C.sfSound
}

const (
	StatStopped Status = C.sfStopped
	StatPaused  Status = C.sfPaused
	StatPlaying Status = C.sfPlaying
)

type Status C.sfSoundStatus

func NewSound() Sound {
	return Sound{C.sfSound_create()}
}

func (s Sound) Play() {
	C.sfSound_play(s.internal)
}

func (s Sound) SetBuffer(sb SoundBuffer) {
	C.sfSound_setBuffer(s.internal, sb.internal)
}

func (s Sound) SetVolume(v float64) {
	C.sfSound_setVolume(s.internal, (C.float)(v))
}

func (s Sound) SetPitch(v float64) {
	C.sfSound_setPitch(s.internal, C.float(v))
}

func (s Sound) Status() Status {
	return Status(C.sfSound_getStatus(s.internal))
}

func (s Sound) SetRelative(rel bool) {
	if rel {
		C.sfSound_setRelativeToListener(s.internal, C.sfBool(1))
	} else {
		C.sfSound_setRelativeToListener(s.internal, C.sfBool(0))
	}
}

func (s Sound) SetPosition(x, y, z float32) {
	C.sfSound_setPosition(s.internal, C.sfVector3f{
		C.float(x),
		C.float(y),
		C.float(z),
	})
}

func (s Sound) SetMinDistance(dist float64) {
	C.sfSound_setMinDistance(s.internal, C.float(dist))
}

func (s Sound) SetAttenuation(att float64) {
	C.sfSound_setAttenuation(s.internal, C.float(att))
}

func (s Sound) Free() {
	C.sfSound_destroy(s.internal)
}

func SetListenerPosition(x, y, z float32) {
	C.sfListener_setPosition(C.sfVector3f{
		C.float(x),
		C.float(y),
		C.float(z),
	})
}

func SetListenerDirection(x, y, z float32) {
	C.sfListener_setDirection(C.sfVector3f{
		C.float(x),
		C.float(y),
		C.float(z),
	})
}

func SetGlobalVolume(vol float64) {
	C.sfListener_setGlobalVolume(C.float(vol))
}

type Music struct {
	internal *C.sfMusic
	f        string
}

// NewMusic creates a new music object from the reader
// Note: this will close the Reader if closable
func NewMusic(r io.Reader) Music {
	if c, ok := r.(io.Closer); ok {
		defer c.Close()
	}
	m := Music{}
	f, ok := r.(*os.File)
	if !ok {
		var err error
		f, err = ioutil.TempFile("", "stream_")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, err = io.Copy(f, r)
		if err != nil {
			panic(err)
		}
		m.f = f.Name()
	}
	str := C.CString(f.Name())
	defer C.free(unsafe.Pointer(str))
	m.internal = C.sfMusic_createFromFile(str)
	return m
}

func (m Music) Play() {
	C.sfMusic_play(m.internal)
}

func (m Music) Stop() {
	C.sfMusic_stop(m.internal)
}

func (m Music) SetVolume(v float64) {
	C.sfMusic_setVolume(m.internal, (C.float)(v))
}

func (m Music) SetPitch(v float64) {
	C.sfMusic_setPitch(m.internal, C.float(v))
}

func (m Music) Status() Status {
	return Status(C.sfMusic_getStatus(m.internal))
}

func (m Music) Free() {
	C.sfMusic_destroy(m.internal)
	if m.f != "" {
		os.Remove(m.f)
	}
}
