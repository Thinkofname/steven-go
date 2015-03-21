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
}

func (w *writing) writeStruct(s *ast.TypeSpec, name string) {
	spec := s.Type.(*ast.StructType)
	for _, field := range spec.Fields.List {
		tag := reflect.StructTag("")
		if field.Tag != nil {
			tag = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
		}
		for _, n := range field.Names {
			w.writeType(field.Type, fmt.Sprintf("%s.%s", name, n), tag)
		}
	}
}

func (w *writing) writeType(e ast.Expr, name string, tag reflect.StructTag) {
	switch e := e.(type) {
	case *ast.SelectorExpr:
		pck := e.X.(*ast.Ident).Name
		s := e.Sel.Name
		if tag.Get("as") != "" {
			w.writeNamed(pck+"."+s, name, tag)
		}
	case *ast.StarExpr:
		if tag.Get("as") == "" {
			panic("pointers can only be 'as'd")
		}
		w.writeType(e.X, name, tag)
	case *ast.Ident:
		w.writeNamed(e.Name, name, tag)
	case *ast.ArrayType:
		lT := tag.Get("length")
		w.writeNamed(lT, fmt.Sprintf("%s(len(%s))", lT, name), "")
		if i, ok := e.Elt.(*ast.Ident); ok && (i.Name == "byte" || i.Name == "uint8") {
			fmt.Fprintf(&w.buf, "if _, err = ww.Write(%s); err != nil { return }\n", name)
		} else {
			iVar := w.tmp()
			fmt.Fprintf(&w.buf, "for %s := range %s {\n", iVar, name)
			w.writeType(e.Elt, fmt.Sprintf("%s[%s]", name, iVar), "")
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
				if err = writeString(ww, %[3]s); err != nil { return err }
			`, tmp1, name, tmp2)
		case "raw":
			fmt.Fprintf(&w.buf, "if err = %s.Serialize(ww); err != nil { return err }\n", name)
		default:
			fmt.Fprintf(&w.buf, "// Can't 'as' %s\n", as)
		}
		return
	}
	if s, ok := structs[t]; ok {
		w.writeStruct(s, name)
		return
	}
	funcName := ""

	switch t {
	case "VarInt":
		funcName = "writeVarInt"
	case "string":
		funcName = "writeString"
	case "bool":
		funcName = "writeBool"
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
	case "int64", "uint64":
		w.scratch(8)
		generateNumberWrite(&w.buf, name, t, 8, t[0] == 'u')
	default:
		fmt.Fprintf(&w.buf, "// TODO write to %s type %s\n", name, t)
	}
	if len(funcName) != 0 {
		fmt.Fprintf(&w.buf, "if err = %s(ww, %s); err != nil { return err }\n", funcName, name)
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
