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
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/thinkofdeath/steven/native"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/render/ui"
	"github.com/thinkofdeath/steven/render/ui/scene"
	"github.com/thinkofdeath/steven/resource"
)

type uiLogo struct {
	scene *scene.Type

	text          *ui.Text
	origX         float64
	textBaseScale float64
}

var (
	r                         = rand.New(rand.NewSource(time.Now().UnixNano()))
	logoTexture               = "blocks/cobblestone"
	logoTargetTexture         = "blocks/cobblestone"
	logoText                  string
	logoTextTimer             float64
	logoLayers                [2][]*ui.Image
	logoTimer, logoTransTimer float64
	stevenLogo                string
)

func (u *uiLogo) init(scene *scene.Type) {
	if logoText == "" {
		nextLogoText()
	}
	if logoTexture == "" {
		nextLogoTexture()
		logoTexture = logoTargetTexture
	}
	readStevenLogo()
	u.scene = scene
	row := 0
	tex, tex2 := render.GetTexture(logoTexture), render.GetTexture(logoTargetTexture)
	titleBox := ui.NewContainer(0, 8, 0, 0).Attach(ui.Top, ui.Center)
	logoTimer = r.Float64() * 60 * 30
	logoTransTimer = 120
	for _, line := range strings.Split(stevenLogo, "\n") {
		if line == "" {
			continue
		}
		for i, r := range line {
			if r == ' ' {
				continue
			}
			x, y := i*4, row*8
			rr, gg, bb := 255, 255, 255
			if r != ':' {
				rr, gg, bb = 170, 170, 170
			}
			shadow := ui.NewImage(
				render.GetTexture("solid"),
				float64(x+2), float64(y+4), 4, 8,
				float64(x%16)/16.0, float64(y%16)/16.0, 4/16.0, 8/16.0,
				0, 0, 0,
			)
			shadow.SetA(100)
			shadow.AttachTo(titleBox)
			u.scene.AddDrawable(shadow)

			img := ui.NewImage(
				tex,
				float64(x), float64(y), 4, 8,
				float64(x%16)/16.0, float64(y%16)/16.0, 4/16.0, 8/16.0,
				rr, gg, bb,
			)
			img.AttachTo(titleBox)
			u.scene.AddDrawable(img)
			logoLayers[0] = append(logoLayers[0], img)

			img = ui.NewImage(
				tex2,
				float64(x), float64(y), 4, 8,
				float64(x%16)/16.0, float64(y%16)/16.0, 4/16.0, 8/16.0,
				rr, gg, bb,
			)
			img.AttachTo(titleBox)
			img.SetA(0)
			u.scene.AddDrawable(img)
			logoLayers[1] = append(logoLayers[1], img)
			if titleBox.Width() < float64(x+4) {
				titleBox.SetWidth(float64(x + 4))
			}
		}
		row++
	}
	titleBox.SetHeight(float64(row) * 8.0)

	txt := ui.NewText(logoText, 0, -8, 255, 255, 0).Attach(ui.Bottom, ui.Right)
	txt.AttachTo(titleBox)
	txt.SetRotation(-math.Pi / 8)
	u.scene.AddDrawable(txt)
	u.text = txt
	width, _ := txt.Size()
	u.textBaseScale = 300 / width
	if u.textBaseScale > 1 {
		u.textBaseScale = 1
	}
	txt.SetX((-txt.Width / 2) * u.textBaseScale)
	u.origX = txt.X()
}

func (u *uiLogo) tick(delta float64) {
	if logoTimer > 0 {
		logoTimer -= delta
	} else if logoTransTimer < 0 {
		logoTransTimer = 120
		logoTimer = r.Float64() * 60 * 30
		readStevenLogo()
		logoTexture = logoTargetTexture
		nextLogoTexture()
		nextLogoText()
		u.text.Update(logoText)
		width, _ := u.text.Size()
		u.textBaseScale = 300 / width
		if u.textBaseScale > 1 {
			u.textBaseScale = 1
		}
		u.text.SetX((-u.text.Width / 2) * u.textBaseScale)
		u.origX = u.text.X()
	} else {
		logoTransTimer -= delta
	}

	tex, tex2 := render.GetTexture(logoTexture), render.GetTexture(logoTargetTexture)
	for i := range logoLayers[0] {
		logoLayers[0][i].SetTexture(tex)
		logoLayers[1][i].SetTexture(tex2)

		logoLayers[0][i].SetA(int(255 * (logoTransTimer / 120)))
		logoLayers[1][i].SetA(int(255 * (1 - (logoTransTimer / 120))))
	}

	logoTextTimer += delta
	if logoTextTimer > 60 {
		logoTextTimer -= 60
	}
	off := (logoTextTimer / 30)
	if off > 1.0 {
		off = 2.0 - off
	}
	off = (math.Cos(off*math.Pi) + 1) / 2
	u.text.SetScaleX((0.7 + (off / 3)) * u.textBaseScale)
	u.text.SetScaleY((0.7 + (off / 3)) * u.textBaseScale)
	u.text.SetX(u.origX * u.text.ScaleX() * u.textBaseScale)
}

func nextLogoText() {
	lines := make([]string, len(stevenLogoLines))
	copy(lines, stevenLogoLines)

	rs, _ := resource.OpenAll("minecraft", "texts/splashes.txt")
	for _, r := range rs {
		func() {
			defer r.Close()
			s := bufio.NewScanner(r)
			for s.Scan() {
				line := s.Text()
				if line != "" && !strings.ContainsRune(line, '§') {
					switch line {
					case "Now Java 6!":
						line = "Now Go!"
					case "OpenGL 2.1 (if supported)!":
						line = "OpenGL 3.2!"
					}
					lines = append(lines, line)
				}
			}
		}()
	}

	logoText = lines[r.Intn(len(lines))]
}

var stevenLogoLines = []string{
	"Your machine uses " + native.Order.String() + " byte order!",
	fmt.Sprintf("You have %d CPUs!", runtime.NumCPU()),
	fmt.Sprintf("Compiled for %s with a %s CPU!", runtime.GOOS, runtime.GOARCH),
	"Compiled with " + runtime.Version() + "!",
	fmt.Sprintf("Splash generated at %d", time.Now().Unix()),
}

func readStevenLogo() {
	r, _ := resource.Open("steven", "logo/logo.txt")
	var str string
	func() {
		defer r.Close()
		data, _ := ioutil.ReadAll(r)
		str = string(data)
	}()
	stevenLogo = str
}

func nextLogoTexture() {
	textures := make([]string, 0)

	rs, _ := resource.OpenAll("steven", "logo/textures.txt")
	for _, r := range rs {
		func() {
			defer r.Close()
			s := bufio.NewScanner(r)
			for s.Scan() {
				texture := s.Text()
				textures = append(textures, texture)
			}
		}()
	}

	if len(textures) > 0 {
		logoTargetTexture = textures[r.Intn(len(textures))]
	}
}
