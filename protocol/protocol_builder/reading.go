package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"reflect"
)

type reading struct {
	out *bytes.Buffer
	buf bytes.Buffer

	scratchSize int
	tmpCount    int
}

func (r *reading) readStruct(s *ast.TypeSpec, name string) {
	spec := s.Type.(*ast.StructType)
	for _, field := range spec.Fields.List {
		tag := reflect.StructTag("")
		if field.Tag != nil {
			tag = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
		}
		for _, n := range field.Names {
			r.readType(field.Type, fmt.Sprintf("%s.%s", name, n), tag)
		}
	}
}

func (r *reading) readType(e ast.Expr, name string, tag reflect.StructTag) {
	switch e := e.(type) {
	case *ast.SelectorExpr:
		pck := e.X.(*ast.Ident).Name
		s := e.Sel.Name
		if tag.Get("as") != "" {
			r.readNamed(pck+"."+s, name, tag)
		}
	case *ast.StarExpr:
		if tag.Get("as") == "" {
			panic("pointers can only be 'as'd")
		}
		ty := e.X.(*ast.Ident).Name
		fmt.Fprintf(&r.buf, "%s = new(%s)\n", name, ty)
		r.readNamed(ty, name, tag)
	case *ast.Ident:
		r.readNamed(e.Name, name, tag)
	case *ast.ArrayType:
		lT := tag.Get("length")
		lenVar := r.tmp()
		fmt.Fprintf(&r.buf, "var %s %s\n", lenVar, lT)
		r.readNamed(lT, lenVar, "")
		// Allocate the slice

		fmt.Fprintf(&r.buf, "%s = make([]%s, %s)\n", name, e.Elt, lenVar)
		if i, ok := e.Elt.(*ast.Ident); ok && (i.Name == "byte" || i.Name == "uint8") {
			fmt.Fprintf(&r.buf, "if _, err = rr.Read(%s); err != nil { return }\n", name)
		} else {
			iVar := r.tmp()
			fmt.Fprintf(&r.buf, "for %s := range %s {\n", iVar, name)
			r.readType(e.Elt, fmt.Sprintf("%s[%s]", name, iVar), "")
			r.buf.WriteString("}\n")
		}
	default:
		fmt.Fprintf(&r.buf, "// Unhandled %#v\n", e)
	}
}

func (r *reading) readNamed(t, name string, tag reflect.StructTag) {
	as := tag.Get("as")
	if as != "" {
		switch as {
		case "json":
			imports["encoding/json"] = struct{}{}
			tmp1 := r.tmp()
			fmt.Fprintf(&r.buf, `var %[1]s string 
			if %[1]s, err = readString(rr); err != nil { return err }
			if err = json.Unmarshal([]byte(%[1]s), &%[2]s); err != nil { return err }
			`, tmp1, name)
		case "raw":
			fmt.Fprintf(&r.buf, "if err = %s.Deserialize(rr); err != nil { return err }\n", name)
		default:
			fmt.Fprintf(&r.buf, "// Can't 'as' %s\n", as)
		}
		return
	}
	if s, ok := structs[t]; ok {
		r.readStruct(s, name)
		return
	}
	funcName := ""
	origName := name
	origT := t

	switch t {
	case "VarInt":
		funcName = "readVarInt"
	case "string":
		funcName = "readString"
	case "bool":
		funcName = "readBool"
	case "int8", "uint8", "byte":
		r.scratch(1)
		generateNumberRead(&r.buf, name, t, 1, t[0] != 'i')
	case "int16", "uint16":
		r.scratch(2)
		generateNumberRead(&r.buf, name, t, 2, t[0] == 'u')
	case "float32":
		imports["math"] = struct{}{}
		name = r.tmp()
		t = "uint32"
		fmt.Fprintf(&r.buf, "var %s uint32\n", name)
		fallthrough
	case "int32", "uint32":
		r.scratch(4)
		generateNumberRead(&r.buf, name, t, 4, t[0] == 'u')
		if origT == "float32" {
			fmt.Fprintf(&r.buf, "%s = math.Float32frombits(%s)\n", origName, name)
		}
	case "float64":
		imports["math"] = struct{}{}
		name = r.tmp()
		t = "uint64"
		fmt.Fprintf(&r.buf, "var %s uint64\n", name)
		fallthrough
	case "int64", "uint64":
		r.scratch(8)
		generateNumberRead(&r.buf, name, t, 8, t[0] == 'u')
		if origT == "float64" {
			fmt.Fprintf(&r.buf, "%s = math.Float64frombits(%s)\n", origName, name)
		}
	default:
		fmt.Fprintf(&r.buf, "// TODO read to %s type %s\n", name, t)
	}
	if len(funcName) != 0 {
		fmt.Fprintf(&r.buf, "if %s, err = %s(rr); err != nil { return err }\n", name, funcName)
	}
}

func (r *reading) tmp() string {
	r.tmpCount++
	return fmt.Sprintf("tmp%d", r.tmpCount-1)
}

func (r *reading) scratch(size int) {
	if size > r.scratchSize {
		r.scratchSize = size
	}
}

func (r *reading) flush() {
	if r.scratchSize > 0 {
		fmt.Fprintf(r.out, "var tmp [%d]byte\n", r.scratchSize)
	}

	r.buf.WriteTo(r.out)
}
