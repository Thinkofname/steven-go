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

package scene

import (
	"encoding/xml"
	"reflect"
	"strconv"
	"strings"

	"github.com/thinkofdeath/steven/resource"
	"github.com/thinkofdeath/steven/ui"
)

func LoadScene(plugin, name string, handler interface{}) *Type {
	s := New(false)
	r, err := resource.Open(plugin, "scene/"+name+".xml")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	v := reflect.ValueOf(handler)
	d := xml.NewDecoder(r)
	for {
		t, err := d.Token()
		if t == nil {
			break
		}
		if err != nil {
			panic(err)
		}
		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Local != "scene" {
				panic("Excepted scene got " + t.Name.Local)
			}
			decodeScene(s, d, nil, v)
		}
	}
	return s
}

func decodeScene(s *Type, d *xml.Decoder, parent ui.Drawable, handler reflect.Value) (children []ui.Drawable) {
	for {
		t, err := d.Token()
		if t == nil {
			return
		}
		if err != nil {
			panic(err)
		}
		var dec decTag
		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "button":
				dec = &buttonTag{}
			case "text":
				dec = &textTag{}
			case "textbox":
				dec = &textboxTag{}
			default:
				panic("Unknown tag " + t.Name.Local)
			}
			if dec == nil {
				panic("Invalid state")
			}
			dec.parse(t.Attr)
		case xml.EndElement:
			return
		default:
			continue
		}
		draw := dec.decode(handler)
		if dec.ID() != "" {
			ty := handler.Type()
			h := handler
			if ty.Kind() == reflect.Ptr {
				ty = ty.Elem()
				h = h.Elem()
			}
			for i := 0; i < ty.NumField(); i++ {
				f := ty.Field(i)
				if id := f.Tag.Get("uiID"); id == dec.ID() {
					h.Field(i).Set(reflect.ValueOf(draw))
				}
			}
		}
		children = append(children, draw)
		draw.AttachTo(parent)
		s.AddDrawable(draw)
		sub := decodeScene(s, d, draw, handler)
		dec.post(draw, sub)
		children = append(children, sub...)
	}
}

type tag struct {
	id     string
	VAlign ui.AttachPoint
	HAlign ui.AttachPoint
	X      float64
	Y      float64
}

func (t *tag) parse(x []xml.Attr) {
	var err error
	for _, attr := range x {
		switch attr.Name.Local {
		case "x":
			t.X, err = strconv.ParseFloat(attr.Value, 64)
		case "y":
			t.Y, err = strconv.ParseFloat(attr.Value, 64)
		case "valign":
			t.VAlign = decodeVAlign(attr.Value)
		case "halign":
			t.HAlign = decodeHAlign(attr.Value)
		case "id":
			t.id = attr.Value
		}
		if err != nil {
			panic(err)
		}
	}
}

func (t tag) ID() string { return t.id }

func (tag) post(d ui.Drawable, c []ui.Drawable) {}

type decTag interface {
	ID() string
	parse([]xml.Attr)
	decode(handler reflect.Value) ui.Drawable
	post(d ui.Drawable, c []ui.Drawable)
}

type buttonTag struct {
	tag
	Width    float64
	Height   float64
	Disabled bool
	Click    string
}

func (b *buttonTag) parse(x []xml.Attr) {
	b.tag.parse(x)
	var err error
	for _, attr := range x {
		switch attr.Name.Local {
		case "width":
			b.Width, err = strconv.ParseFloat(attr.Value, 64)
		case "height":
			b.Height, err = strconv.ParseFloat(attr.Value, 64)
		case "click":
			b.Click = attr.Value
		case "disabled":
			b.Disabled, err = strconv.ParseBool(attr.Value)
		}
		if err != nil {
			panic(err)
		}
	}
}

var ClickSound func()

func (b buttonTag) decode(handler reflect.Value) ui.Drawable {
	btn := ui.NewButton(b.X, b.Y, b.Width, b.Height)
	if ClickSound != nil {
		btn.AddClick(ClickSound)
	}
	if b.Click != "" {
		if strings.HasPrefix(b.Click, ".") {
			click := b.Click[1:]
			m := handler.MethodByName(click)
			btn.AddClick(m.Interface().(func()))
		}
	}
	btn.SetDisabled(b.Disabled)
	return btn.Attach(
		b.VAlign, b.HAlign,
	)
}
func (b buttonTag) post(d ui.Drawable, c []ui.Drawable) {
	btn := d.(*ui.Button)
	for _, draw := range c {
		if txt, ok := draw.(*ui.Text); ok {
			oldB := txt.B()
			newB := int(float64(txt.B()) * 0.63)
			btn.AddHover(func(over bool) {
				if over && !btn.Disabled() {
					txt.SetB(newB)
				} else {
					txt.SetB(oldB)
				}
			})
		}
	}
}

type textTag struct {
	tag
	Value   string
	R, G, B int
}

func (t *textTag) parse(x []xml.Attr) {
	t.tag.parse(x)
	var err error
	t.R = 255
	t.G = 255
	t.B = 255
	for _, attr := range x {
		switch attr.Name.Local {
		case "value":
			t.Value = attr.Value
		case "r":
			t.R, err = strconv.Atoi(attr.Value)
		case "g":
			t.G, err = strconv.Atoi(attr.Value)
		case "b":
			t.B, err = strconv.Atoi(attr.Value)
		}
		if err != nil {
			panic(err)
		}
	}
}

func (t textTag) decode(handler reflect.Value) ui.Drawable {
	txt := ui.NewText(t.Value, t.X, t.Y, t.R, t.G, t.B)
	return txt.Attach(
		t.VAlign, t.HAlign,
	)
}

func decodeHAlign(a string) ui.AttachPoint {
	switch a {
	case "center":
		return ui.Center
	case "right":
		return ui.Right
	default:
		return ui.Left
	}
}

func decodeVAlign(a string) ui.AttachPoint {
	switch a {
	case "middle":
		return ui.Middle
	case "bottom":
		return ui.Bottom
	default:
		return ui.Top
	}
}

type textboxTag struct {
	tag
	Width    float64
	Height   float64
	Password bool
	Submit   string
}

func (t *textboxTag) parse(x []xml.Attr) {
	t.tag.parse(x)
	var err error
	for _, attr := range x {
		switch attr.Name.Local {
		case "width":
			t.Width, err = strconv.ParseFloat(attr.Value, 64)
		case "height":
			t.Height, err = strconv.ParseFloat(attr.Value, 64)
		case "password":
			t.Password, err = strconv.ParseBool(attr.Value)
		case "submit":
			t.Submit = attr.Value
		}
		if err != nil {
			panic(err)
		}
	}
}
func (t textboxTag) decode(handler reflect.Value) ui.Drawable {
	txt := ui.NewTextBox(t.X, t.Y, t.Width, t.Height)
	txt.SetPassword(t.Password)
	if t.Submit != "" {
		if strings.HasPrefix(t.Submit, ".") {
			submit := t.Submit[1:]
			m := handler.MethodByName(submit)
			txt.SubmitFunc = m.Interface().(func())
		}
	}
	return txt.Attach(
		t.VAlign, t.HAlign,
	)
}
