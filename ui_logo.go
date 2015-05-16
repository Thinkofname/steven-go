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
	"math"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/thinkofdeath/steven/native"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/ui"
	"github.com/thinkofdeath/steven/ui/scene"
)

type uiLogo struct {
	scene *scene.Type

	text          *ui.Text
	origX         float64
	textBaseScale float64
}

var (
	logoTextures = []string{
		"blocks/cobblestone",
		"blocks/netherrack",
		"blocks/dirt",
		"blocks/planks_oak",
		"blocks/brick",
		"blocks/snow",
		"blocks/sand",
		"blocks/gravel",
		"blocks/hardened_clay",
		"blocks/clay",
		"blocks/bedrock",
		"blocks/obsidian",
		"blocks/end_stone",
		"blocks/stone_andesite",
		"blocks/dirt_podzol_top",
		"blocks/portal",
		"blocks/prismarine_rough",
		"blocks/soul_sand",
		"blocks/lava_still",
		"blocks/hay_block_top",
		"blocks/log_acacia",
		"blocks/red_sandstone_carved",
	}
	r             = rand.New(rand.NewSource(time.Now().UnixNano()))
	logoTexture   = logoTextures[r.Intn(len(logoTextures))]
	logoText      = stevenLogoLines[r.Intn(len(stevenLogoLines))]
	logoTextTimer float64
)

func (u *uiLogo) init(scene *scene.Type) {
	u.scene = scene
	row := 0
	tex := render.GetTexture(logoTexture)
	titleBox := (&ui.Container{
		Y: 8,
	}).Attach(ui.Top, ui.Center)
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
			shadow.A = 100
			shadow.Parent = titleBox
			u.scene.AddDrawable(shadow)
			img := ui.NewImage(
				tex,
				float64(x), float64(y), 4, 8,
				float64(x%16)/16.0, float64(y%16)/16.0, 4/16.0, 8/16.0,
				rr, gg, bb,
			)
			img.Parent = titleBox
			u.scene.AddDrawable(img)
			if titleBox.W < float64(x+4) {
				titleBox.W = float64(x + 4)
			}
		}
		row++
	}
	titleBox.H = float64(row) * 8.0

	txt := ui.NewText(logoText, 0, -8, 255, 255, 0).Attach(ui.Bottom, ui.Right)
	txt.Parent = titleBox
	txt.Rotation = -math.Pi / 8
	u.scene.AddDrawable(txt)
	u.text = txt
	width, _ := txt.Size()
	u.textBaseScale = 300 / width
	if u.textBaseScale > 1 {
		u.textBaseScale = 1
	}
	txt.X = (-txt.Width / 2) * u.textBaseScale
	u.origX = txt.X
}

func (u *uiLogo) tick(delta float64) {
	logoTextTimer += delta
	if logoTextTimer > 60 {
		logoTextTimer -= 60
	}
	off := (logoTextTimer / 30)
	if off > 1.0 {
		off = 2.0 - off
	}
	off = (math.Cos(off*math.Pi) + 1) / 2
	u.text.ScaleX = (0.7 + (off / 3)) * u.textBaseScale
	u.text.ScaleY = (0.7 + (off / 3)) * u.textBaseScale
	u.text.X = u.origX
	u.text.X *= u.text.ScaleX * u.textBaseScale
}

var stevenLogoLines = []string{
	"I blame Xor",
	"Its not a bug its a feature!",
	"Don't go to #think, tis a silly place",
	"Tested! (In production)",
	"Not in scala!",
	"Its steven not phteven!",
	"Now webscale!",
	"Meow",
	"I bet one of cindy's cats broke it!",
	"=^.^=",
	"ಠ_ಠ",
	"Commit reverted in 5..4..3...",
	"Latest is greatest!",
	"[This space is intentionally left blank]",
	"ThinkBot: .... *Thinkofdeath damn it",
	"Now with more bugs!",
	"I blame Mojang",
	"The logo is totally not ascii art rendered as textures",
	"Look, it works on my machine.",
	"Open Source! https://github.com/thinkofdeath/steven",
	"Built with Go!",
	"Your machine uses " + native.Order.String() + " byte order!",
	fmt.Sprintf("You have %d CPUs!", runtime.NumCPU()),
	fmt.Sprintf("Compiled for %s with a %s CPU!", runtime.GOOS, runtime.GOARCH),
	"Compiled with " + runtime.Version() + "!",
	"try { } catch (Exception e) { }",
	"panic(recover())",
	"// Abandon hope all ye who enter here",
	"Its like I'm racing vanilla to see who can have the most bugs",
	"Using ascii art for the logo seemed like a bad idea at first",
	"... and still does.",
	"Help! I'm trapped in the splash text!",
	"Linux support!",
	"Windows support!",
	"Mac support! (in theory)",
	"Could have used vanilla's splash text!",
	"Come chat on IRC!",
	"Knowing Murphy's Law doesn't help",
	"Minecraft Multi-processing: breaking three things at once",
	"Silly Mortal...",
	"Software isn't released. It's allowed to escape.",
	"General System Error: Please sacrifice a cow and two chickens to continue",
	"Do you want to build a client?",
	"sudo rm -rf --no-preserve-root /",
	fmt.Sprintf("Splash generated at %d", time.Now().Unix()),
	"Thinkofdeath.getClass().getField(\"sanity\").set(Thinkofdeath, null);",
	"There is no God, only Zuul",
	"Now with potatoes!",
	"ask :: String -> Int; ask x = 42",
	"And then you cleanse them in a ball of atomic fire!",
	"I'm a little matrix, square and stout,",
	"this is my transpose, this is my count.",
	"<?php if(\"6 geese\" + \"4 chickens\" == \"10 birds\") echo(\"lolphp\");",
	"All hail the Cat Godess!",
	"var f = function() { return f; };",
	"unsafe.allocateObject(Unsafe.class);",
	":(){ :|:& };:",
	"You must pay to view this content.",
	"No touchy the topic",
	"Ceci n'est pas un splash.",
	"Xor is not actually a cat.",
	"The MD5 of md_5 is e14cfacdd442a953343ebd8529138680",
	"Gas powered stick! It never runs out of gas!",
}

const stevenLogo = `
   SSSSSSSSSSSSSSS          tttt                                                                                             
 SS:::::::::::::::S      ttt:::t                                                                                             
S:::::SSSSSS::::::S      t:::::t                                                                                             
S:::::S     SSSSSSS      t:::::t                                                                                             
S:::::S            ttttttt:::::ttttttt        eeeeeeeeeeee    vvvvvvv           vvvvvvv    eeeeeeeeeeee    nnnn  nnnnnnnn    
S:::::S            t:::::::::::::::::t      ee::::::::::::ee   v:::::v         v:::::v   ee::::::::::::ee  n:::nn::::::::nn  
 S::::SSSS         t:::::::::::::::::t     e::::::eeeee:::::ee  v:::::v       v:::::v   e::::::eeeee:::::een::::::::::::::nn 
  SS::::::SSSSS    tttttt:::::::tttttt    e::::::e     e:::::e   v:::::v     v:::::v   e::::::e     e:::::enn:::::::::::::::n
    SSS::::::::SS        t:::::t          e:::::::eeeee::::::e    v:::::v   v:::::v    e:::::::eeeee::::::e  n:::::nnnn:::::n
       SSSSSS::::S       t:::::t          e:::::::::::::::::e      v:::::v v:::::v     e:::::::::::::::::e   n::::n    n::::n
            S:::::S      t:::::t          e::::::eeeeeeeeeee        v:::::v:::::v      e::::::eeeeeeeeeee    n::::n    n::::n
            S:::::S      t:::::t    tttttte:::::::e                  v:::::::::v       e:::::::e             n::::n    n::::n
SSSSSSS     S:::::S      t::::::tttt:::::te::::::::e                  v:::::::v        e::::::::e            n::::n    n::::n
S::::::SSSSSS:::::S      tt::::::::::::::t e::::::::eeeeeeee           v:::::v          e::::::::eeeeeeee    n::::n    n::::n
S:::::::::::::::SS         tt:::::::::::tt  ee:::::::::::::e            v:::v            ee:::::::::::::e    n::::n    n::::n
 SSSSSSSSSSSSSSS             ttttttttttt      eeeeeeeeeeeeee             vvv               eeeeeeeeeeeeee    nnnnnn    nnnnnn                                         
`
