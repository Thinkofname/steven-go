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
	"github.com/thinkofdeath/steven/ui"
)

var memoryStats runtime.MemStats

func (c *ClientState) initDebug() {
	c.debug.position = ui.NewText("X:0 Y:0 Z:0", 5, 5, 255, 255, 255).
		Attach(ui.Top, ui.Left)
	c.scene.AddDrawable(c.debug.position)
	c.debug.facing = ui.NewText("Facing: invalid", 5, 5+18, 255, 255, 255).
		Attach(ui.Top, ui.Left)
	c.scene.AddDrawable(c.debug.facing)
	c.debug.rotation = ui.NewText("Yaw: 0 Pitch: 0", 5, 5+18*2, 255, 255, 255).
		Attach(ui.Top, ui.Left)
	c.scene.AddDrawable(c.debug.rotation)

	c.debug.fps = ui.NewText("FPS: 0", 5, 5, 255, 255, 255).
		Attach(ui.Top, ui.Right)
	c.scene.AddDrawable(c.debug.fps)
	c.debug.memory = ui.NewText("0/0", 5, 5+18, 255, 255, 255).
		Attach(ui.Top, ui.Right)
	c.scene.AddDrawable(c.debug.memory)

	c.debug.target = ui.NewText("", 5, 5+18*2, 255, 255, 255).
		Attach(ui.Top, ui.Right)
	c.scene.AddDrawable(c.debug.target)
	c.debug.targetName = ui.NewText("", 5, 5+18*3, 255, 255, 255).
		Attach(ui.Top, ui.Right)
	c.scene.AddDrawable(c.debug.targetName)
	c.debug.enabled = true
	c.toggleDebug()
}

func (c *ClientState) toggleDebug() {
	c.debug.enabled = !c.debug.enabled
	e := c.debug.enabled
	c.debug.position.SetDraw(e)
	c.debug.facing.SetDraw(e)
	c.debug.rotation.SetDraw(e)
	c.debug.fps.SetDraw(e)
	c.debug.memory.SetDraw(e)
	c.debug.target.SetDraw(e)
	c.debug.targetName.SetDraw(e)
	for _, t := range c.debug.targetInfo {
		t[0].SetDraw(e)
		t[1].SetDraw(e)
	}
}

func (c *ClientState) renderDebug() {
	if !c.debug.enabled {
		return
	}
	c.debug.position.Update(fmt.Sprintf("X: %.2f, Y: %.2f, Z: %.2f", c.X, c.Y, c.Z))
	c.debug.facing.Update(fmt.Sprintf("Facing: %s", c.facingDirection()))
	c.debug.rotation.Update(fmt.Sprintf("Yaw: %.2f Pitch: %.2f", c.Yaw, c.Pitch))
	
	c.displayTargetInfo()

	runtime.ReadMemStats(&memoryStats)
	c.debug.memory.Update(fmt.Sprintf("%s/%s", formatMemory(memoryStats.Alloc), formatMemory(memoryStats.Sys)))

	c.debug.frames++
	now := time.Now()
	if now.Sub(c.debug.lastCount) >= time.Second {
		c.debug.lastCount = now
		c.debug.fpsValue = c.debug.frames
		c.debug.frames = 0
	}
	c.debug.fps.Update(fmt.Sprintf("FPS: %d", c.debug.fpsValue))
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
	t, b, face, _ := c.targetBlock()
	c.debug.target.Update(fmt.Sprintf("Target(%d,%d,%d)-%s", t.X, t.Y, t.Z, face))
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
				c.scene.AddDrawable(t)
			}
		}
		v := c.debug.targetInfo[i][0]
		v.SetDraw(true)
		v.SetR(r)
		v.SetG(g)
		v.SetB(b)
		v.Update(text)
		k := c.debug.targetInfo[i][1]
		k.SetDraw(true)
		k.SetX(7 + v.Width)
		k.Update(fmt.Sprintf("%s=", s.Key))
	}
	for i := len(b.states()); i < len(c.debug.targetInfo); i++ {
		info := &c.debug.targetInfo[i]
		info[0].SetDraw(false)
		info[1].SetDraw(false)
	}
}
