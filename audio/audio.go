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
	"unsafe"
)

// #cgo LDFLAGS: -lcsfml-audio
// #include <SFML/Audio/SoundBuffer.h>
// #include <SFML/Audio/Sound.h>
// #include <stdlib.h>
import "C"

type SoundBuffer struct {
	internal *C.sfSoundBuffer
}

func NewSoundBuffer(file string) SoundBuffer {
	return SoundBuffer{C.sfSoundBuffer_createFromFile(C.CString(file))}
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

func (s Sound) Status() Status {
	return Status(C.sfSound_getStatus(s.internal))
}

func (s Sound) Free() {
	C.sfSound_destroy(s.internal)
}
