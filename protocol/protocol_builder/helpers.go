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
	"bytes"
	"fmt"
)

func generateNumberWrite(w *bytes.Buffer, name string, t string, size int, unsigned bool) {
	for i := 0; i < size; i++ {
		fmt.Fprintf(w, "tmp[%d] =  byte(%s >> %d)\n", i, name, (size-1-i)*8)
	}
	fmt.Fprintf(w, "if _, err = ww.Write(tmp[:%d]); err != nil { return }\n", size)
}

func generateNumberRead(w *bytes.Buffer, name, t string, size int, unsigned bool) {
	origT := t
	if !unsigned {
		t = "u" + t
	}
	fmt.Fprintf(w, "if _, err = rr.Read(tmp[:%d]); err != nil { return  }\n", size)
	fmt.Fprintf(w, "%s = ", name)
	if !unsigned {
		fmt.Fprintf(w, "%s(", origT)
	}
	for i := 0; i < size; i++ {
		fmt.Fprintf(w, "(%s(tmp[%d]) << %d)", t, size-1-i, i*8)
		if i != size-1 {
			w.WriteString("|")
		}
	}
	if !unsigned {
		w.WriteString(")")
	}
	w.WriteString("\n")
}
