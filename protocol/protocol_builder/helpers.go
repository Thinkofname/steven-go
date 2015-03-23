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
