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
	"bytes"
	"encoding/json"
	"errors"
)

var null = []byte("null")

// AnyComponent can respresent any valid chat component. This
// is designed to be used as a target for json unmarshalling
// where the type of the component isn't known.
//
// This also has the benefit of allowing plain strings or
// arrays to be converted into text components simplifing
// handling of those cases.
type AnyComponent struct {
	Value interface{}
}

// Wrap wraps the passed component with an AnyComponent
func Wrap(i interface{}) AnyComponent {
	return AnyComponent{Value: i}
}

// Type returns the type of the component this contains.
// It is genernally prefered to type switch over the result
// of Value().
func (s AnyComponent) Type() Type {
	switch s.Value.(type) {
	case *TextComponent:
		return Text
	case *TranslateComponent:
		return Translate
	case *ScoreComponent:
		return Score
	case *SelectorComponent:
		return Selector
	}
	return Invalid
}

// MarshalJSON marshals the contained value to json.
func (s *AnyComponent) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Value)
}

var (
	errInvalid = errors.New("invalid component")
)

// UnmarshalJSON unmarshals the passed json as a component.
// This allows for objects, arrays and strings.
func (s *AnyComponent) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errInvalid
	}
	var err error
	switch b[0] {
	case '{':
		var m map[string]json.RawMessage
		err = json.Unmarshal(b, &m)
		if err != nil {
			return err
		}
		if _, ok := m["text"]; ok {
			s.Value = &TextComponent{}
		} else if _, ok := m["translate"]; ok {
			s.Value = &TranslateComponent{}
		} else if _, ok := m["score"]; ok {
			s.Value = &ScoreComponent{}
		} else if _, ok := m["selector"]; ok {
			s.Value = &SelectorComponent{}
		}
		if s.Value == nil {
			return errInvalid
		}
		err = mapStruct(s.Value, m)
	case '[':
		t := &TextComponent{
			Text: "",
		}
		err = json.Unmarshal(b, &t.Extra)
		s.Value = t
	case '"':
		var v string
		err = json.Unmarshal(b, &v)
		s.Value = &TextComponent{Text: v}
	case 'n':
		if bytes.Equal(b, null) {
			break
		}
		fallthrough
	default:
		return errInvalid
	}
	return err
}

//go:generate stringer -type=Type

// Type represents a type of a component.
type Type int

const (
	// Invalid is returned for unknown component types.
	Invalid Type = iota
	// Text is a simple text only component.
	Text
	// Translate is a translatable component.
	Translate
	// Score is a component whos contents depends on a scoreboard value.
	Score
	// Selector is a component whos contents follows Minecraft's
	// selector rules.
	Selector
)

type (
	// Component is the base for any chat component.
	Component struct {
		Extra         []AnyComponent `json:"extra,omitempty"`
		Bold          *bool          `json:"bold,omitempty"`
		Italic        *bool          `json:"italic,omitempty"`
		Underlined    *bool          `json:"underlined,omitempty"`
		Strikethrough *bool          `json:"strikethrough,omitempty"`
		Obfuscated    *bool          `json:"obfuscated,omitempty"`
		Color         Color          `json:"color,omitempty"`

		ClickEvent *ClickEvent `json:"clickEvent,omitempty"`
		HoverEvent *HoverEvent `json:"hoverEvent,omitempty"`
		Insertion  string      `json:"insertion,omitempty"`
	}
	// TextComponent is a component with a plain text value.
	TextComponent struct {
		Text string `json:"text,omitempty"`
		Component
	}
	// TranslateComponent is a component whos value is loaded from
	// a locale file on the client based on the Translate key and
	// the client's language, substituting in values from the With
	// slice.
	TranslateComponent struct {
		Translate string         `json:"translate,omitempty"`
		With      []AnyComponent `json:"with,omitempty"`
		Component
	}
	// ScoreComponent is a component whos value is based on the
	// contents of the scoreboard.
	ScoreComponent struct {
		Score struct {
			Name      string `json:"name,omitempty"`
			Objective string `json:"objective,omitempty"`
		} `json:"score"`
		Component
	}
	// SelectorComponent is a component whos value is based on
	// evalutating the contained selector.
	SelectorComponent struct {
		Selector string `json:"selector,omitempty"`
		Component
	}
)

type (
	// Color represents one of the 16 valid colors in Minecraft.
	Color string
	// ClickAction is an action that will be preformed on clicking.
	ClickAction string
	// HoverAction is an action that will be preformed on hovering.
	HoverAction string
)

// Valid colors.
const (
	Black       Color = "black"
	DarkBlue    Color = "dark_blue"
	DarkGreen   Color = "dark_green"
	DarkAqua    Color = "dark_aqua"
	DarkRed     Color = "dark_red"
	DarkPurple  Color = "dark_purple"
	Gold        Color = "gold"
	Gray        Color = "gray"
	DarkGray    Color = "dark_gray"
	Blue        Color = "blue"
	Green       Color = "green"
	Aqua        Color = "aqua"
	Red         Color = "red"
	LightPurple Color = "light_purple"
	Yellow      Color = "yellow"
	White       Color = "white"
)

// Optional boolean values
var (
	vTrue         = true
	True          = &vTrue
	vFalse        = false
	False         = &vFalse
	Missing *bool = nil
)

// Valid ClickActions.
const (
	OpenURL        ClickAction = "open_url"
	OpenFile       ClickAction = "open_file"
	RunCommand     ClickAction = "run_command"
	SuggestCommand ClickAction = "suggest_command"
)

// Valid HoverActions.
const (
	ShowText        HoverAction = "show_text"
	ShowAchievement HoverAction = "show_achievement"
	ShowItem        HoverAction = "show_item"
)

// ClickEvent is an event which will be preformed when the
// area of text is clicked.
type ClickEvent struct {
	Action ClickAction `json:"action"`
	Value  string      `json:"value"`
}

// HoverEvent is an event which will be preformed when the
// area of text is hovered over.
type HoverEvent struct {
	Action HoverAction `json:"action"`
	Value  string      `json:"value"`
}
