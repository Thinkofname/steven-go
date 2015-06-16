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

package console

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/thinkofdeath/steven/format"
)

var (
	// stdout and a log file combind
	w io.Writer

	// For the possible option of scrolling in the
	// future
	historyBuffer [200]format.AnyComponent
	lock          sync.Mutex

	defaultRegistry registry
)

func checkInit() {
	if w != nil {
		return
	}
	f, err := os.Create("steven-log.txt")
	if err != nil {
		panic(err)
	}
	w = io.MultiWriter(f, os.Stdout)
}

// Text appends the passed formatted string plus a new line to
// the log buffer. The formatting uses the same rules as fmt.
func Text(f string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}
	file = file[strings.LastIndex(file, "/")+1:]

	msg := &format.TextComponent{
		Text: fmt.Sprintf("[%s:%d] ", file, line),
		Component: format.Component{
			Color: format.Aqua,
		},
	}
	msg.Extra = append(msg.Extra, format.Wrap(&format.TextComponent{
		Text: fmt.Sprintf(f, args...),
		Component: format.Component{
			Color: format.White,
		},
	}))
	Component(format.Wrap(msg))
}

// Component appends the component to the log buffer.
func Component(c format.AnyComponent) {
	lock.Lock()
	defer lock.Unlock()
	checkInit()
	io.WriteString(w, c.String()+"\n")
	copy(historyBuffer[1:], historyBuffer[:])
	historyBuffer[0] = c
}

// History returns up to the requested number of lines from the log
// buffer.
//
// As a special case -1 will return the whole buffer.
func History(lines int) []format.AnyComponent {
	lock.Lock()
	defer lock.Unlock()
	if lines == -1 || lines > len(historyBuffer) {
		return historyBuffer[:]
	}
	return historyBuffer[:lines]
}
