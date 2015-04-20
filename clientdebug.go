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

package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/thinkofdeath/steven/chat"
	"github.com/thinkofdeath/steven/render"
)

var memoryStats runtime.MemStats

func (c *ClientState) renderDebug() {
	render.DrawUIText(
		fmt.Sprintf("X: %.2f, Y: %.2f, Z: %.2f", c.X, c.Y, c.Z),
		5, 5, 255, 255, 255,
	)
	render.DrawUIText(
		fmt.Sprintf("Facing: %s", c.facingDirection()),
		5, 23, 255, 255, 255,
	)
	c.displayTargetInfo()

	runtime.ReadMemStats(&memoryStats)
	text := fmt.Sprintf("%s/%s", formatMemory(memoryStats.Alloc), formatMemory(memoryStats.Sys))
	render.DrawUIText(text, 800-5-float64(render.SizeOfString(text)), 23, 255, 255, 255)

	now := time.Now()
	if now.Sub(c.lastCount) > time.Second {
		c.lastCount = now
		c.fps = c.frames
		c.frames = 0
	}
	text = fmt.Sprintf("FPS: %d", c.fps)
	render.DrawUIText(text, 800-5-float64(render.SizeOfString(text)), 5, 255, 255, 255)
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
	text := fmt.Sprintf("Target(%d,%d,%d)", tx, ty, tz)
	render.DrawUIText(
		text,
		800-5-render.SizeOfString(text), 41, 255, 255, 255,
	)
	text = fmt.Sprintf("%s:%s", b.Plugin(), b.Name())
	render.DrawUIText(
		text,
		800-5-render.SizeOfString(text), 59, 255, 255, 255,
	)

	for i, s := range b.states() {
		var r, g, b int = 255, 255, 255
		text = fmt.Sprint(s.Value)
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
		pos := 800 - 5 - render.SizeOfString(text)
		render.DrawUIText(
			text,
			pos, 59+18*(1+float64(i)), r, g, b,
		)
		text = fmt.Sprintf("%s=", s.Key)
		pos -= render.SizeOfString(text) + 2
		render.DrawUIText(
			text,
			pos, 59+18*(1+float64(i)), 255, 255, 255,
		)
	}
}
