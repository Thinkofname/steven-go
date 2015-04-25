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
	"runtime"
	"time"

	"github.com/thinkofdeath/steven/chat"
	"github.com/thinkofdeath/steven/render/ui"
)

var memoryStats runtime.MemStats

func (c *ClientState) initDebug() {
	c.debug.position = ui.NewText("X:0 Y:0 Z:0", 5, 5, 255, 255, 255).
		Attach(ui.Top, ui.Left)
	ui.AddDrawable(c.debug.position)
	c.debug.facing = ui.NewText("Facing: invalid", 5, 23, 255, 255, 255).
		Attach(ui.Top, ui.Left)
	ui.AddDrawable(c.debug.facing)

	c.debug.fps = ui.NewText("FPS: 0", 5, 5, 255, 255, 255).
		Attach(ui.Top, ui.Right)
	ui.AddDrawable(c.debug.fps)
	c.debug.memory = ui.NewText("0/0", 5, 23, 255, 255, 255).
		Attach(ui.Top, ui.Right)
	ui.AddDrawable(c.debug.memory)

	c.debug.target = ui.NewText("", 5, 41, 255, 255, 255).
		Attach(ui.Top, ui.Right)
	ui.AddDrawable(c.debug.target)
	c.debug.targetName = ui.NewText("", 5, 59, 255, 255, 255).
		Attach(ui.Top, ui.Right)
	ui.AddDrawable(c.debug.targetName)
	c.debug.enabled = true
	c.toggleDebug()
}

func (c *ClientState) toggleDebug() {
	c.debug.enabled = !c.debug.enabled
	e := c.debug.enabled
	c.debug.position.Visible = e
	c.debug.facing.Visible = e
	c.debug.fps.Visible = e
	c.debug.memory.Visible = e
	c.debug.target.Visible = e
	c.debug.targetName.Visible = e
	for _, t := range c.debug.targetInfo {
		t[0].Visible = e
		t[1].Visible = e
	}
}

func (c *ClientState) renderDebug() {
	if !c.debug.enabled {
		return
	}
	c.debug.position.Update(fmt.Sprintf("X: %.2f, Y: %.2f, Z: %.2f", c.X, c.Y, c.Z))
	c.debug.facing.Update(fmt.Sprintf("Facing: %s", c.facingDirection()))

	c.displayTargetInfo()

	runtime.ReadMemStats(&memoryStats)
	c.debug.memory.Update(fmt.Sprintf("%s/%s", formatMemory(memoryStats.Alloc), formatMemory(memoryStats.Sys)))

	now := time.Now()
	if now.Sub(c.lastCount) > time.Second {
		c.lastCount = now
		c.fps = c.frames
		c.frames = 0
	}
	c.debug.fps.Update(fmt.Sprintf("FPS: %d", c.fps))
}

func formatMemory(alloc uint64) string {
	const letters = "BKMG"
	i := 0
	for {
		check := alloc
		check >>= 10
		if check == 0 {
			break
		}
		alloc = check
		i++
	}
	l := string(letters[i])
	if l != "B" {
		l += "B"
	}
	return fmt.Sprintf("%d%s", alloc, l)
}

var debugStateColors = [...]chat.Color{
	cWhite:     chat.White,
	cOrange:    chat.Gold,
	cMagenta:   chat.LightPurple,
	cLightBlue: chat.Aqua,
	cYellow:    chat.Yellow,
	cLime:      chat.Green,
	cPink:      chat.Red,
	cGray:      chat.Gray,
	cSilver:    chat.DarkGray,
	cCyan:      chat.DarkAqua,
	cPurple:    chat.DarkPurple,
	cBlue:      chat.Blue,
	cBrown:     chat.Gold,
	cGreen:     chat.DarkGreen,
	cRed:       chat.DarkRed,
	cBlack:     chat.Black,
}

func (c *ClientState) displayTargetInfo() {
	tx, ty, tz, b := c.targetBlock()
	c.debug.target.Update(fmt.Sprintf("Target(%d,%d,%d)", tx, ty, tz))
	c.debug.targetName.Update(fmt.Sprintf("%s:%s", b.Plugin(), b.Name()))

	for i, s := range b.states() {
		var r, g, b int = 255, 255, 255
		text := fmt.Sprint(s.Value)
		switch val := s.Value.(type) {
		case bool:
			b = 0
			if val {
				g = 255
				r = 0
			} else {
				r = 255
				g = 0
			}
		case color:
			r, g, b = chatColorRGB(debugStateColors[val])
		}
		if i >= len(c.debug.targetInfo) {
			c.debug.targetInfo = append(c.debug.targetInfo, [2]*ui.Text{})
			c.debug.targetInfo[i] = [2]*ui.Text{
				ui.NewText("", 5, 59+18*(1+float64(i)), 255, 255, 255).Attach(ui.Top, ui.Right),
				ui.NewText("", 5, 59+18*(1+float64(i)), 255, 255, 255).Attach(ui.Top, ui.Right),
			}
			for _, t := range c.debug.targetInfo[i] {
				ui.AddDrawable(t)
			}
		}
		v := c.debug.targetInfo[i][0]
		v.Visible = true
		v.R, v.G, v.B = r, g, b
		v.Update(text)
		k := c.debug.targetInfo[i][1]
		k.Visible = true
		k.X = 7 + v.Width
		k.Update(fmt.Sprintf("%s=", s.Key))
	}
	for i := len(b.states()); i < len(c.debug.targetInfo); i++ {
		info := &c.debug.targetInfo[i]
		info[0].Visible = false
		info[1].Visible = false
	}
}
