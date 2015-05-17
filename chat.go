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
	"math"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/chat"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/ui"
)

const (
	chatHistoryLines = 10
	maxLineWidth     = 500
)

type ChatUI struct {
	Lines [chatHistoryLines]chat.AnyComponent

	container       *ui.Container
	parts           []*chatLine
	input           *ui.Text
	inputBackground *ui.Image

	enteringText    bool
	wasEnteringText bool
	inputLine       []rune
	cursorTick      float64
	first           bool
}

type chatLine struct {
	fade       float64
	text       *ui.Formatted
	background *ui.Image
}

func (c *ChatUI) init() {
	c.container = &ui.Container{
		X: 0,
		Y: 44,
		W: 500,
		H: chatHistoryLines*18 + 2,
	}
	c.container.Attach(ui.Bottom, ui.Left)
	c.input = ui.NewText("", 5, 1, 255, 255, 255).Attach(ui.Bottom, ui.Left)
	c.input.Visible = false
	c.input.Parent = c.container
	c.inputBackground = ui.NewImage(render.GetTexture("solid"), 0, 0, 500, 20, 0, 0, 1, 1, 0, 0, 0).Attach(ui.Bottom, ui.Left)
	c.inputBackground.A = 77
	c.inputBackground.Parent = c.container
	c.inputBackground.Visible = false
	Client.scene.AddDrawable(c.inputBackground)
	Client.scene.AddDrawable(c.input)
}

func (c *ChatUI) Draw(delta float64) {
	if c.wasEnteringText != c.enteringText {
		if c.wasEnteringText {
			c.input.Visible = false
			c.inputBackground.Visible = false
			c.input.Update(string(c.inputLine))
			for _, p := range c.parts {
				p.text.Y += 18
				p.background.Y += 18
			}
		} else {
			for _, p := range c.parts {
				p.text.Y -= 18
				p.background.Y -= 18
			}
		}
		c.wasEnteringText = c.enteringText
	}
	if c.enteringText {
		c.input.Visible = true
		c.inputBackground.Visible = true
		c.cursorTick += delta
		// Add on our cursor
		if int(c.cursorTick/30)%2 == 0 {
			c.input.Update(string(c.inputLine) + "|")
		} else {
			c.input.Update(string(c.inputLine))
		}
		// Lazy way of preventing rounding errors buiding up over time
		if c.cursorTick > 0xFFFFFF {
			c.cursorTick = 0
		}
	}
	parts := c.parts
	offset := 0
	limit := 0.0
	if c.enteringText {
		limit = -18
	}
	for i, p := range parts {
		if p.background.Y < limit {
			c.parts = c.parts[i+1-offset:]
			offset = i + 1
			p.text.Remove()
			p.background.Remove()
		} else {
			p.fade -= 0.005 * delta
			if p.fade < 0 {
				p.fade = 0
			}
			for _, t := range p.text.Text {
				if c.enteringText {
					t.A = 1.0
				} else {
					t.A = p.fade
				}
			}
			ba := 0.3
			if !c.enteringText {
				ba -= (1.0 - p.fade) / 2.0
				ba = math.Min(ba, 0.3)
			}
			p.background.A = int(255 * ba)
			if p.background.A < 0 {
				p.background.A = 0
			}
		}
	}
}

func (c *ChatUI) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if (key == glfw.KeyEscape || key == glfw.KeyEnter) && action == glfw.Release {
		if key == glfw.KeyEnter && len(c.inputLine) != 0 {
			Client.network.Write(&protocol.ChatMessage{string(c.inputLine)})
		}
		// Return control back to the default
		c.enteringText = false
		c.inputLine = c.inputLine[:0]
		lockMouse = true
		w.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		w.SetCharCallback(nil)
		return
	}
	if key == glfw.KeyBackspace && action != glfw.Release {
		if len(c.inputLine) > 0 {
			c.inputLine = c.inputLine[:len(c.inputLine)-1]
		}
	}
}

func (c *ChatUI) handleChar(w *glfw.Window, char rune) {
	if c.first {
		c.first = false
		return
	}
	if len(c.inputLine) < 100 {
		c.inputLine = append(c.inputLine, char)
	}
}

func chatColorRGB(c chat.Color) (r, g, b int) {
	switch c {
	case chat.Black:
		return 0, 0, 0
	case chat.DarkBlue:
		return 0, 0, 170
	case chat.DarkGreen:
		return 0, 170, 0
	case chat.DarkAqua:
		return 0, 170, 170
	case chat.DarkRed:
		return 170, 0, 0
	case chat.DarkPurple:
		return 170, 0, 170
	case chat.Gold:
		return 255, 170, 0
	case chat.Gray:
		return 170, 170, 170
	case chat.DarkGray:
		return 85, 85, 85
	case chat.Blue:
		return 85, 85, 255
	case chat.Green:
		return 85, 255, 85
	case chat.Aqua:
		return 85, 255, 255
	case chat.Red:
		return 255, 85, 85
	case chat.LightPurple:
		return 255, 85, 255
	case chat.Yellow:
		return 255, 255, 85
	case chat.White:
		return 255, 255, 255

	}
	return 255, 255, 255
}

func (c *ChatUI) Add(msg chat.AnyComponent) {
	chat.ConvertLegacy(msg)
	copy(c.Lines[0:chatHistoryLines-1], c.Lines[1:])
	c.Lines[chatHistoryLines-1] = msg
	f := ui.NewFormattedWidth(msg, 5, chatHistoryLines*18+1, 500-10).Attach(ui.Top, ui.Left)
	f.Parent = c.container
	line := &chatLine{
		text:       f,
		fade:       3.0,
		background: ui.NewImage(render.GetTexture("solid"), 0, chatHistoryLines*18, 500, f.Height, 0, 0, 1, 1, 0, 0, 0),
	}
	line.background.Parent = c.container
	line.background.A = 77
	c.parts = append(c.parts, line)
	Client.scene.AddDrawable(line.background)
	Client.scene.AddDrawable(f)
	ff := f
	for _, f := range c.parts {
		f.text.Y -= 18 * float64(ff.Lines)
		f.background.Y -= 18 * float64(ff.Lines)
	}
	if c.enteringText {
		ff.Y -= 18
		line.background.Y -= 18
	}
}
