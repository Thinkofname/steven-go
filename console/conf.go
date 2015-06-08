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
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/thinkofdeath/steven/chat"
)

// ExecConf executes/loads the passed config file
func ExecConf(path string) {
	defer func() {
		if path == "conf.cfg" {
			saveConf()
		}
	}()
	f, err := os.Open(path)
	if err != nil {
		Component(chat.Build("Failed to open ").
			Append(path).
			Append(": ").
			Append(err.Error()).
			Create(),
		)
		return
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		txt := strings.TrimSpace(s.Text())
		if strings.HasPrefix(txt, "#") || txt == "" {
			continue
		}
		if err := Execute(txt); err != nil {
			Component(chat.Build("Error: " + err.Error()).Color(chat.Red).Create())
		}
	}
}

func saveConf() {
	f, err := os.Create("conf.cfg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	sort.Sort(cvars)
	for _, c := range cvars {
		if c.props().is(Serializable) {
			fmt.Fprintln(f, c)
			fmt.Fprintln(f)
		}
	}
}

type cvarList []cvarI

func (c cvarList) Len() int           { return len(c) }
func (c cvarList) Swap(a, b int)      { c[a], c[b] = c[b], c[a] }
func (c cvarList) Less(a, b int) bool { return c[a].Name() < c[b].Name() }
