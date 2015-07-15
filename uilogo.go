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
	r                         = rand.New(rand.NewSource(time.Now().UnixNano()))
	logoTexture               = logoTextures[r.Intn(len(logoTextures))]
	logoTargetTexture         = logoTextures[r.Intn(len(logoTextures))]
	logoText                  string
	logoTextTimer             float64
	logoLayers                [2][]*ui.Image
	logoTimer, logoTransTimer float64
)

func (u *uiLogo) init(scene *scene.Type) {
	if logoText == "" {
		nextLogoText()
	}
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
		logoTexture = logoTargetTexture
		logoTargetTexture = logoTextures[r.Intn(len(logoTextures))]
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
<<<<<<< HEAD:uilogo.go
=======
	"Splashes by: `git log --format='%aN' ui_logo.go | sort -u`",
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
>>>>>>> origin/master:ui_logo.go
	"Your machine uses " + native.Order.String() + " byte order!",
	fmt.Sprintf("You have %d CPUs!", runtime.NumCPU()),
	fmt.Sprintf("Compiled for %s with a %s CPU!", runtime.GOOS, runtime.GOARCH),
	"Compiled with " + runtime.Version() + "!",
<<<<<<< HEAD:uilogo.go
	fmt.Sprintf("Splash generated at %d", time.Now().Unix()),
=======
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
	fmt.Sprintf("Splash generated at %s", time.Now().Unix()),
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
	"We are we are we are we are we are the Engineers!",
	"There is no spoon!",
	"re-cur-sion (n.) - see: recursion",
	"throw null;",
	"CafeBabe, kill Akkarin",
	"ls steven | grep 'Herobrine' | rm",
	"Buy one cat, get one free!",
	"3 = π for large values of 3!",
	"Buzzword driven development for synergetic paradigm shifts!",
	"Real programmers use cat and echo!",
	"10 is sufficiently prime for our purposes!",
	"3 divides every number 100% of the time 33% of the time!",
	"void* boom(void) { return &(*NULL); }",
	"Netherrack is also in Go!",
	"More \"indended behaviors\" than PHP!",
	"Loved by about 50 people, give or take!",
	"Supports Xor's broken JSON files!",
	"Passes the Turing test more often than Xor!",
	fmt.Sprintf("Number of the day: %d", r.Int()),
	"I am a taco.",
	"Vegan options available!",
	"Are you not entertained?!",
	"Does this look like the face of mercy?",
	"(╯°□°）╯︵ ┻━┻",
	"chmod -x chmod",
	"Release the kraken!",
	"We're all mad here!",
	"double sin(double theta) { return theta; }",
	"That's not how any of this works!",
	"Hail Hydrate!",
	"The elevator is worthy!",
	"They're illusions, Michael!",
	"Akkarin set mode on #think: [+g *scala*]",
	"[] + [] == ''",
	"[] + {} == [Object object]",
	"{} + [] == 0",
	"{} + {} == NaN",
	"Array(16).join('wat' - 1) + ' Batman!'",
	"Gson.setLenient(true);",
	"V1c5MUlHeHZiMnRsWkNFSwo=",
	"Heil Akkarin!",
	"Praise Helix!",
	"Praise Dome!",
	"GunfighterJ doesn't actually own guns!",
	"Enterprise™ Quality!",
	"My splashes are better than yours!",
	"Hypixel you!",
	"It's the hash-slinging slasher!",
	"Try thinkmap, it's unfinished!",
	"git reset HEAD --hard",
	"This algorithm assumes today is Monday.",
	"Don't touch my Easy-Bake oven!",
	"It's a fanfic about yawkat reading fanfic!",
	"If you're seeing this, you haven't crashed (yet)!",
	"Woo spigotmc.org!",
	"Woo github!",
	"f = λx.f",
	"lazy val x: Nothing = x",
	"Akkarin is actually a sentient potato!",
	"I love you too!",
	"i^2 = j^2 = k^2 = ijk = -1",
	"Wow, such client, much protocol, very multiplayer.",
	"Don't feed the cats!",
	"Retina support!",
	"Projects need names. This one's named Steven!",
	"Not the duct tape!",
	"sudo love me!",
	"Oh how we laughed and laughed! (Only I wasn't laughing.)",
	"Better than Half-Life 3!",
	"Free AND not on Steam!",
	"Presented by Thinkofdeath!",
	"What do you mean you're out of doughnuts?!",
	"How Can Our Eyes Be Real If Our Eyes Arent Real.",
	"¿Porque no los dos?",
	"Elizabeth II plays it!",
	"steven.stevenLogoLines[139]",
	"Phrasing!",
	"PotatOS support!",
	"WAT",
>>>>>>> origin/master:ui_logo.go
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
