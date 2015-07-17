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
	"crypto/rand"
	"encoding/json"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/format"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
	"github.com/thinkofdeath/steven/resource"
)

type resourceList struct {
	baseUI
	scene *scene.Type
	logo  uiLogo

	packs []*resourceListItem
	ret   func() screen
}

type resourceListItem struct {
	*scene.Type

	X, Y      float64
	container *ui.Container
	offset    float64
	id        string
}

func newResourceList(ret func() screen) screen {
	rl := &resourceList{
		scene: scene.New(true),
		ret:   ret,
	}
	rl.logo.init(rl.scene)

	uiFooter(rl.scene)

	rl.redraw()

	refresh, txt := newButtonText("Refresh", 300, -50-15, 100, 30)
	rl.scene.AddDrawable(refresh.Attach(ui.Center, ui.Middle))
	rl.scene.AddDrawable(txt)
	refresh.AddClick(rl.redraw)

	done, txt := newButtonText("Done", 200, -50-15, 100, 30)
	rl.scene.AddDrawable(done.Attach(ui.Center, ui.Middle))
	rl.scene.AddDrawable(txt)
	done.AddClick(func() {
		setScreen(newOptionMenu(rl.ret))
	})

	return rl
}

func (rl *resourceList) init() {
	window.SetScrollCallback(rl.onScroll)
	window.SetKeyCallback(rl.handleKey)
}

func (rl *resourceList) onScroll(w *glfw.Window, xoff float64, yoff float64) {
	if len(rl.packs) == 0 {
		return
	}
	diff := yoff / 10
	if s := rl.packs[len(rl.packs)-1]; s.offset+diff <= 2 {
		diff = 2 - s.offset
	}
	if s := rl.packs[0]; s.offset+diff >= 0 {
		diff = -s.offset
	}
	for _, s := range rl.packs {
		s.offset += diff
		s.updatePosition()
	}
}

func (ri *resourceListItem) updatePosition() {
	if ri.offset < 0 {
		ri.Y = ri.offset * 200
	} else if ri.offset >= 2 {
		ri.Y = ri.offset * 100
	} else {
		ri.Y = ri.offset * 100
	}
}

func (rl *resourceList) redraw() {
	for _, s := range rl.packs {
		s.Hide()
		render.FreeIcon(s.id)
	}
	rl.packs = rl.packs[:0]

	os.MkdirAll("./resource-packs", 0777)
	files, _ := ioutil.ReadDir("./resource-packs")

	for i, f := range files {
		f := f
		if !strings.HasSuffix(f.Name(), ".zip") {
			continue
		}

		fullName := filepath.Join("./resource-packs", f.Name())
		desc, iimg, ok := getPackInfo(fullName)
		if !ok {
			continue
		}

		sc := scene.New(true)
		container := ui.NewContainer(0, float64(i)*100, 700, 100).
			Attach(ui.Center, ui.Middle)
		r := make([]byte, 20)
		rand.Read(r)
		si := &resourceListItem{
			Type:      sc,
			container: container,
			offset:    float64(i),
			id:        "servericon:" + string(r),
		}
		si.updatePosition()
		rl.packs = append(rl.packs, si)

		var rr, gg, bb int
		if resource.IsActive(fullName) {
			rr = 200
			gg = 200
		}

		bck := ui.NewImage(render.GetTexture("solid"), 0, 0, 700, 100, 0, 0, 1, 1, rr, gg, bb).Attach(ui.Top, ui.Left)
		bck.SetA(100)
		bck.AttachTo(container)
		sc.AddDrawable(bck)
		txt := ui.NewText(f.Name(), 90+10, 5, 255, 255, 255).Attach(ui.Top, ui.Left)
		txt.AttachTo(container)
		sc.AddDrawable(txt)

		var tex render.TextureInfo
		if iimg == nil {
			tex = render.GetTexture("misc/unknown_pack")
		} else {
			render.AddIcon(si.id, iimg)
			tex = render.Icon(si.id)
		}
		icon := ui.NewImage(tex, 5, 5, 90, 90, 0, 0, 1, 1, 255, 255, 255).
			Attach(ui.Top, ui.Left)
		icon.AttachTo(container)
		sc.AddDrawable(icon)

		msg := format.Wrap(&format.TextComponent{Text: desc})
		format.ConvertLegacy(msg)
		motd := ui.NewFormattedWidth(msg, 90+10, 5+18, 700-(90+10+5)).Attach(ui.Top, ui.Left)
		motd.AttachTo(container)
		sc.AddDrawable(motd)
		container.ClickFunc = func() {
			if resource.IsActive(fullName) {
				RemovePack(fullName)
			} else {
				AddPack(fullName)
			}
			setScreen(newResourceList(rl.ret))
		}
		container.HoverFunc = func(over bool) {
			if over {
				bck.SetA(200)
			} else {
				bck.SetA(100)
			}
		}

		sc.AddDrawable(container)
	}
}

func getPackInfo(name string) (desc string, icon image.Image, ok bool) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return
	}
	z, err := zip.NewReader(f, stat.Size())
	if err != nil {
		return
	}
	for _, e := range z.File {
		if e.Name == "pack.mcmeta" {
			func() {
				f, err := e.Open()
				if err != nil {
					return
				}
				defer f.Close()
				type meta struct {
					Pack struct {
						Description string
					}
				}
				m := &meta{}
				err = json.NewDecoder(f).Decode(m)
				if err != nil {
					return
				}
				desc = m.Pack.Description
			}()
		}
		if e.Name == "pack.png" {
			func() {
				f, err := e.Open()
				if err != nil {
					return
				}
				defer f.Close()
				icon, _ = png.Decode(f)
			}()
		}
	}
	ok = true
	return
}

func (rl *resourceList) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Release {
		setScreen(newOptionMenu(rl.ret))
	}
}

func (rl *resourceList) tick(delta float64) {
	rl.logo.tick(delta)
	for _, s := range rl.packs {
		dx := s.X - s.container.X()
		dy := s.Y - s.container.Y()
		if dx*dx > 1 {
			s.container.SetX(s.container.X() + delta*dx*0.1)
		} else {
			s.container.SetX(s.X)
		}
		if dy*dy > 1 {
			s.container.SetY(s.container.Y() + delta*dy*0.1)
		} else {
			s.container.SetY(s.Y)
		}
	}
}

func (rl *resourceList) remove() {
	window.SetScrollCallback(onScroll)
	window.SetKeyCallback(onKey)
	rl.scene.Hide()
	for _, s := range rl.packs {
		s.Hide()
		render.FreeIcon(s.id)
	}
}
