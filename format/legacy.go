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

package format

import (
	"strings"
	"unicode"
)

const legacyChar = '§'

func ConvertLegacy(c AnyComponent) {
	switch c := c.Value.(type) {
	case *TextComponent:
		for _, e := range c.Extra {
			ConvertLegacy(e)
		}
		if strings.ContainsRune(c.Text, legacyChar) {
			text := []rune(c.Text)
			c.Text = ""
			var parts []AnyComponent

			last := 0
			cur := &TextComponent{}
			for i := 0; i < len(text); i++ {
				if text[i] == legacyChar && i+1 < len(text) {
					i++
					r := unicode.ToLower(text[i])
					cur.Text = string(text[last : i-1])
					last = i + 1
					prev := cur
					parts = append(parts, AnyComponent{cur})
					cur = &TextComponent{}
					if !((r >= 'a' && r <= 'f') || (r >= '0' && r <= '9')) {
						cur.Component = prev.Component
					}
					switch r {
					case '0':
						cur.Color = Black
					case '1':
						cur.Color = DarkBlue
					case '2':
						cur.Color = DarkGreen
					case '3':
						cur.Color = DarkAqua
					case '4':
						cur.Color = DarkRed
					case '5':
						cur.Color = DarkPurple
					case '6':
						cur.Color = Gold
					case '7':
						cur.Color = Gray
					case '8':
						cur.Color = DarkGray
					case '9':
						cur.Color = Blue
					case 'a':
						cur.Color = Green
					case 'b':
						cur.Color = Aqua
					case 'c':
						cur.Color = Red
					case 'd':
						cur.Color = LightPurple
					case 'e':
						cur.Color = Yellow
					case 'f':
						cur.Color = White
					case 'k':
						cur.Obfuscated = True
					case 'l':
						cur.Bold = True
					case 'm':
						cur.Strikethrough = True
					case 'n':
						cur.Underlined = True
					case 'o':
						cur.Italic = True
					case 'r':
					}
				}
			}
			if len(text[last:]) > 0 {
				cur.Text = string(text[last:])
				parts = append(parts, AnyComponent{cur})
			}

			c.Extra = append(parts, c.Extra...)
		}
	case *TranslateComponent:
		for _, w := range c.With {
			ConvertLegacy(w)
		}
		for _, e := range c.Extra {
			ConvertLegacy(e)
		}
	default:
		panic("unhandled component")
	}
}
