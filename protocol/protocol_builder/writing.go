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
	"go/ast"
	"reflect"
)

type writing struct {
	out *bytes.Buffer
	buf bytes.Buffer

	scratchSize int
	tmpCount    int
	base        string
}

func (w *writing) writeStruct(spec *ast.StructType, name string) {
	var lastCondition conditions
	for _, field := range spec.Fields.List {
		tag := reflect.StructTag("")
		if field.Tag != nil {
			tag = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
		}

		var condition conditions
		if ifTag := tag.Get("if"); ifTag != "" {
			condition = parseCondition(ifTag)
		}

		if !lastCondition.equals(condition) {
			if lastCondition != nil {
				w.buf.WriteString("}\n")
			}
			if condition != nil {
				condition.print(name, &w.buf)
			}
		}
		lastCondition = condition

		for _, n := range field.Names {
			w.writeType(field.Type, fmt.Sprintf("%s.%s", name, n), tag)
		}
	}
	if lastCondition != nil {
		w.buf.WriteString("}\n")
	}
}

func (w *writing) writeType(e ast.Expr, name string, tag reflect.StructTag) {
	switch e := e.(type) {
	case *ast.StructType:
		w.writeStruct(e, name)
	case *ast.SelectorExpr:
		pck := e.X.(*ast.Ident).Name
		s := e.Sel.Name
		w.writeNamed(pck+"."+s, name, tag)
	case *ast.StarExpr:
		w.writeType(e.X, name, tag)
	case *ast.Ident:
		w.writeNamed(e.Name, name, tag)
	case *ast.ArrayType:
		lT := tag.Get("length")
		if lT != "remaining" && lT != "" && lT[0] != '@' {
			w.writeNamed(lT, fmt.Sprintf("%s(len(%s))", lT, name), "")
		}
		if i, ok := e.Elt.(*ast.Ident); ok && (i.Name == "byte" || i.Name == "uint8") {
			fmt.Fprintf(&w.buf, "if _, err = ww.Write(%s); err != nil { return }\n", name)
		} else {
			iVar := w.tmp()
			fmt.Fprintf(&w.buf, "for %s := range %s {\n", iVar, name)
			w.writeType(e.Elt, fmt.Sprintf("%s[%s]", name, iVar), tag)
			w.buf.WriteString("}\n")
		}
	default:
		fmt.Fprintf(&w.buf, "// Unhandled %#v\n", e)
	}
}

func (w *writing) writeNamed(t, name string, tag reflect.StructTag) {
	as := tag.Get("as")
	if as != "" {
		switch as {
		case "json":
			imports["encoding/json"] = struct{}{}
			tmp1 := w.tmp()
			tmp2 := w.tmp()
			fmt.Fprintf(&w.buf, `var %[1]s []byte
				if %[1]s, err = json.Marshal(&%[2]s); err != nil { return }
				%[3]s := string(%[1]s)
				if err = writeString(ww, %[3]s); err != nil { return  }
			`, tmp1, name, tmp2)
		case "raw":
			fmt.Fprintf(&w.buf, "if err = %s.Serialize(ww); err != nil { return  }\n", name)
		default:
			fmt.Fprintf(&w.buf, "// Can't 'as' %s\n", as)
		}
		return
	}
	if s, ok := structs[t]; ok {
		w.writeStruct(s.Type.(*ast.StructType), name)
		return
	}
	funcName := ""

	switch t {
	case "VarInt":
		funcName = "writeVarInt"
	case "VarLong":
		funcName = "writeVarLong"
	case "string":
		funcName = "writeString"
	case "bool":
		funcName = "writeBool"
	case "Metadata":
		funcName = "writeMetadata"
	case "nbt.Compound":
		funcName = "writeNBT"
	case "int8", "uint8", "byte":
		w.scratch(1)
		generateNumberWrite(&w.buf, name, t, 1, t[0] != 'i')
	case "int16", "uint16":
		w.scratch(2)
		generateNumberWrite(&w.buf, name, t, 2, t[0] == 'u')
	case "float32":
		imports["math"] = struct{}{}
		orig := name
		name = w.tmp()
		fmt.Fprintf(&w.buf, "%s := math.Float32bits(%s)\n", name, orig)
		t = "uint32"
		fallthrough
	case "int32", "uint32":
		w.scratch(4)
		generateNumberWrite(&w.buf, name, t, 4, t[0] == 'u')
	case "float64":
		imports["math"] = struct{}{}
		orig := name
		name = w.tmp()
		fmt.Fprintf(&w.buf, "%s := math.Float64bits(%s)\n", name, orig)
		t = "uint64"
		fallthrough
	case "int64", "uint64", "Position":
		w.scratch(8)
		generateNumberWrite(&w.buf, name, t, 8, t[0] != 'i')
	default:
		fmt.Fprintf(&w.buf, "// TODO write to %s type %s\n", name, t)
	}
	if len(funcName) != 0 {
		if notProtocol {
			funcName = "protocol." + funcName
			imports["github.com/thinkofdeath/steven/protocol"] = struct{}{}
		}
		fmt.Fprintf(&w.buf, "if err = %s(ww, %s); err != nil { return  }\n", funcName, name)
	}
}

func (w *writing) tmp() string {
	w.tmpCount++
	return fmt.Sprintf("tmp%d", w.tmpCount-1)
}

func (w *writing) scratch(size int) {
	if size > w.scratchSize {
		w.scratchSize = size
	}
}

func (w *writing) flush() {
	if w.scratchSize > 0 {
		fmt.Fprintf(w.out, "var tmp [%d]byte\n", w.scratchSize)
	}

	w.buf.WriteTo(w.out)
}
