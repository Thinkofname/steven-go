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

// Package json proxies requests to encoding/json but allows for 'broken'
// json
package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"unicode"
)

func Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		fmt.Printf("Invalid json: '%s' trying to work around it\n", err)
		data, err = clean(data)
		if err == nil {
			err = json.Unmarshal(data, v)
		}
	}
	return err
}

type mode int

const (
	mAny mode = iota
	mEOF
	mObject
	mObjectKey
	mObjectNext
	mArray
	mArrayNext
)

type state struct {
	data   []rune
	out    *bytes.Buffer
	offset int
	mode   mode

	stack []mode

	err error
}

func clean(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	s := &state{
		data: []rune(string(data)),
		out:  &buf,
	}
	for s.Run() {

	}
	return buf.Bytes(), s.Err()
}

func (s *state) Run() bool {
	if s.offset >= len(s.data) {
		s.mode = mEOF
		// TODO Error if unexpected
	}
	if s.err != nil || s.mode == mEOF {
		return false
	}
	switch s.mode {
	case mAny:
		s.mode = s.any()
	case mObject:
		s.mode = s.object()
	case mObjectKey:
		s.mode = s.objectKey()
	case mObjectNext:
		s.mode = s.objectNext()
	case mArray:
		s.mode = s.array()
	case mArrayNext:
		s.mode = s.arrayNext()
	default:
		panic("unhandled case")
	}
	return true
}

func (s *state) any() mode {
	r := s.data[s.offset]
	s.offset++
	switch {
	case r == '{': // Object
		s.out.WriteRune('{')
		s.stack = append(s.stack, mObject)
		return mObject
	case r == '"':
		s.out.WriteRune('"')
		s.parseString()
		return s.endValue()
	case r == '[': // Array
		s.out.WriteRune('[')
		s.stack = append(s.stack, mArray)
		return mArray
	case r == 'f': // False
		s.offset--
		for _, c := range "false" {
			if s.data[s.offset] != c {
				s.err = fmt.Errorf("any: unexpected rune %c", s.data[s.offset])
				return 0
			}
			s.out.WriteRune(s.data[s.offset])
			s.offset++
		}
		return s.endValue()
	case r == 't':
		s.offset--
		for _, c := range "true" {
			if s.data[s.offset] != c {
				s.err = fmt.Errorf("any: unexpected rune %c", s.data[s.offset])
				return 0
			}
			s.out.WriteRune(s.data[s.offset])
			s.offset++
		}
		return s.endValue()
	case unicode.IsDigit(r), r == '-':
		s.out.WriteRune(r)
		s.parseNumber()
		return s.endValue()
	case unicode.IsSpace(r):
		return mAny
	case r == '/':
		s.parseComment()
		return mAny
	}
	s.out.WriteRune('"')
	s.out.WriteRune(r)
	for s.offset < len(s.data) {
		r := s.data[s.offset]
		s.offset++
		if !unicode.IsLetter(r) {
			s.offset--
			break
		}
		s.out.WriteRune(r)
	}
	s.out.WriteRune('"')
	return s.endValue()
}

func (s *state) endValue() mode {
	if len(s.stack) == 0 {
		return mEOF
	}
	m := s.stack[len(s.stack)-1]
	switch m {
	case mObject:
		return mObjectNext
	case mArray:
		return mArrayNext
	default:
		panic(m)
	}
}

func (s *state) array() mode {
	r := s.data[s.offset]
	s.offset++
	switch {
	case r == ']':
		s.out.WriteRune(']')
		s.stack = s.stack[:len(s.stack)-1]
		return s.endValue()
	case unicode.IsSpace(r):
		return mArray
	}
	s.offset--
	return s.any()
}

func (s *state) arrayNext() mode {
	r := s.data[s.offset]
	s.offset++
	switch {
	case r == '/':
		s.parseComment()
		return mArrayNext
	case r == ',':
		s.out.WriteRune(',')
		return mAny
	case r == ']':
		s.out.WriteRune(']')
		s.stack = s.stack[:len(s.stack)-1]
		return s.endValue()
	case unicode.IsSpace(r):
		return mArrayNext
	}
	s.err = fmt.Errorf("arrayNext: unexpected rune %c", r)
	return 0
}

func (s *state) object() mode {
	// Expecting keys
	r := s.data[s.offset]
	s.offset++
	switch {
	case r == '/':
		s.parseComment()
		return mObject
	case r == '"': // Key start
		s.out.WriteRune('"')
		s.parseString()
		return mObjectKey
	case unicode.IsLetter(r): // Broken key
		s.out.WriteRune('"')
		s.out.WriteRune(r)
		for s.offset < len(s.data) {
			r := s.data[s.offset]
			s.offset++
			if !unicode.IsLetter(r) {
				s.offset--
				break
			}
			s.out.WriteRune(r)
		}
		s.out.WriteRune('"')
		return mObjectKey
	case r == '}':
		s.out.WriteRune('}')
		s.stack = s.stack[:len(s.stack)-1]
		return s.endValue()
	case unicode.IsSpace(r):
		return mObject
	}
	s.err = fmt.Errorf("object: unexpected rune %c", r)
	return 0
}

func (s *state) objectKey() mode {
	r := s.data[s.offset]
	s.offset++
	switch {
	case r == ':':
		s.out.WriteRune(':')
		return mAny
	case unicode.IsSpace(r):
		return mObjectKey
	}
	s.err = fmt.Errorf("objectKey: unexpected rune %c", r)
	return 0
}

func (s *state) objectNext() mode {
	r := s.data[s.offset]
	s.offset++
	switch {
	case r == '/':
		s.parseComment()
		return mObjectNext
	case r == ',':
		s.out.WriteRune(',')
		return mObject
	case r == '}':
		s.out.WriteRune('}')
		s.stack = s.stack[:len(s.stack)-1]
		return s.endValue()
	case unicode.IsSpace(r):
		return mObjectNext
	}
	s.err = fmt.Errorf("objectNext: unexpected rune %c", r)
	return 0
}

func (s *state) parseNumber() {
	for s.offset < len(s.data) {
		r := s.data[s.offset]
		s.offset++
		switch {
		case unicode.IsDigit(r), r == '.':
			s.out.WriteRune(r)
		default:
			s.offset--
			return
		}

	}
}

func (s *state) parseString() {
	for s.offset < len(s.data) {
		r := s.data[s.offset]
		s.offset++
		switch r {
		case '"':
			s.out.WriteRune(r)
			return
		case '\\':
			s.out.WriteRune(r)
			r = s.data[s.offset]
			s.offset++
			fallthrough
		default:
			s.out.WriteRune(r)
		}

	}
}

func (s *state) parseComment() {
	t := s.data[s.offset]
	s.offset++
	for s.offset < len(s.data) {
		r := s.data[s.offset]
		s.offset++
		switch r {
		case '\n':
			if t == '/' {
				return
			}
		case '*':
			if s.data[s.offset+1] == '/' || t == '*' {
				s.offset++
				return
			}
		}
	}
}

func (s *state) Err() error {
	return s.err
}
