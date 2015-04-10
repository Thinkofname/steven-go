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

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/thinkofdeath/steven/chat"
	"github.com/thinkofdeath/steven/protocol"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/resource/locale"
)

const (
	chatHistoryLines = 10
	maxLineWidth     = 500
)

type ChatUI struct {
	Elements []*chatUIElement

	dirty    bool
	Lines    [chatHistoryLines]chat.AnyComponent
	lineFade [chatHistoryLines]float64

	lineLength float64

	enteringText bool
	inputLine    []rune
	cursorTick   float64
	first        bool
}

type chatUIElement struct {
	text   *render.UIText
	offset int
	line   int
}

func (c *ChatUI) render(delta float64) {
	// Always redraw in input mode because of the cursor
	if c.enteringText {
		c.dirty = true
	}
	if c.dirty {
		c.dirty = false
		// Clear the existing elements.
		// They will be reused if needed (handled by the
		// ui system)
		for _, e := range c.Elements {
			if e.text != nil {
				e.text.Free()
			}
		}
		c.Elements = c.Elements[:0]

		for i, line := range c.Lines {
			c.newLine()
			if line.Value == nil {
				continue
			}
			c.lineLength = 0
			c.renderComponent(i, line.Value, nil)
		}

		if c.enteringText {
			// Shift all the lines up
			c.newLine()
			c.lineLength = 0

			color := chat.White
			gc := func() chat.Color { return color }
			line := c.inputLine
			// Make it clear that a command is being typed
			if len(line) != 0 && line[0] == '/' {
				color = chat.Gold
				c.renderText(len(c.Lines), line[:1], gc)
				color = chat.Yellow
				line = line[1:]
			}
			c.renderText(len(c.Lines), line, gc)
			c.cursorTick += delta
			// Add on our cursor
			if int(c.cursorTick/30)%2 == 0 {
				c.renderText(len(c.Lines), []rune{'|'}, gc)
			}
			// Lazy way of preventing rounding errors buiding up over time
			if c.cursorTick > 0xFFFFFF {
				c.cursorTick = 0
			}
		}
	}
	// Slowly fade out each line
	for i := range c.lineFade {
		c.lineFade[i] -= 0.005 * delta
		if c.lineFade[i] < 0 {
			c.lineFade[i] = 0
		}
	}
	for _, e := range c.Elements {
		if e.text != nil {
			// If entering text show every line
			if c.enteringText {
				e.text.Alpha(1.0)
			} else {
				e.text.Alpha(c.lineFade[e.line])
			}
		}
	}
}

func (c *ChatUI) handleKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if (key == glfw.KeyEscape || key == glfw.KeyEnter) && action == glfw.Release {
		if key == glfw.KeyEnter && len(c.inputLine) != 0 {
			writeChan <- &protocol.ChatMessage{string(c.inputLine)}
		}
		// Return control back to the default
		c.enteringText = false
		c.inputLine = c.inputLine[:0]
		lockMouse = true
		w.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
		c.dirty = true
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

func (c *ChatUI) renderComponent(line int, co interface{}, color chatGetColorFunc) {
	switch co := co.(type) {
	case *chat.TextComponent:
		getColor := chatGetColor(&co.Component, color)
		c.renderText(line, []rune(co.Text), getColor)
		for _, e := range co.Extra {
			c.renderComponent(line, e.Value, getColor)
		}
	case *chat.TranslateComponent:
		getColor := chatGetColor(&co.Component, color)
		for _, part := range locale.Get(co.Translate) {
			switch part := part.(type) {
			case string:
				c.renderText(line, []rune(part), getColor)
			case int:
				if part < 0 || part >= len(co.With) {
					continue
				}
				c.renderComponent(line, co.With[part].Value, getColor)
			}
		}
		for _, e := range co.Extra {
			c.renderComponent(line, e.Value, getColor)
		}
	default:
		fmt.Printf("Can't handle %T\n", co)
	}
}

func (c *ChatUI) renderText(line int, runes []rune, getColor chatGetColorFunc) {
	width := 0.0
	r, g, b := chatColorRGB(getColor())
	for i := 0; i < len(runes); i++ {
		size := float64(render.SizeOfCharacter(runes[i]))
		if c.lineLength+width+size > maxLineWidth {
			c.appendText(line, string(runes[:i]), r, g, b)
			c.lineLength = 0
			runes = runes[i:]
			i = 0
			width = 0
			c.newLine()
		}
		width += size
	}
	c.lineLength += c.appendText(line, string(runes), r, g, b)
}

func (c *ChatUI) appendText(line int, str string, r, g, b int) float64 {
	e := &chatUIElement{
		text:   render.AddUIText(str, 2+c.lineLength, 480-18, r, g, b),
		offset: 0,
		line:   line,
	}
	c.Elements = append(c.Elements, e)
	return e.text.Width + 2
}

type chatGetColorFunc func() chat.Color

func chatGetColor(c *chat.Component, parent chatGetColorFunc) chatGetColorFunc {
	return func() chat.Color {
		if c.Color != "" {
			return c.Color
		}
		if parent != nil {
			return parent()
		}
		return chat.White
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

func (c *ChatUI) newLine() {
	for _, e := range c.Elements {
		if e.text == nil {
			continue
		}
		e.offset++
		if e.offset > chatHistoryLines {
			e.text.Free()
			e.text = nil
			continue
		}
		e.text.Shift(0, -18)
	}
}

func (c *ChatUI) Add(msg chat.AnyComponent) {
	copy(c.Lines[0:chatHistoryLines-1], c.Lines[1:])
	copy(c.lineFade[0:chatHistoryLines-1], c.lineFade[1:])
	c.Lines[chatHistoryLines-1] = msg
	c.lineFade[chatHistoryLines-1] = 3.0
	c.dirty = true
}
