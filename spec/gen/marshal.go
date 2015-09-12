package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"h12.me/gengo"
)

func fromBNFToGoFuncs() {
	w := os.Stdout
	goFile := fromBNFToGoFile()
	fpl(w, "package proto")
	fpl(w, "import (")
	fpl(w, `"io"`)
	fpl(w, ")")
	for _, decl := range goFile.TypeDecls {
		genMarshalFunc(w, decl)
		fpl(w, "")
	}
}

func genMarshalFunc(w io.Writer, decl *gengo.TypeDecl) {
	t := decl.Type
	switch t.Kind {
	case gengo.StructKind:
		fpl(w, "func (t *%s) Marshal(w io.Writer) error {", decl.Name)
		for i, f := range t.Fields {
			genFieldMarshal(w, i, f)
		}
		fpl(w, "return nil")
		fpl(w, "}")
		return
	case gengo.ArrayKind:
		fpl(w, "func (t %s) Marshal(w io.Writer) error {", decl.Name)
		marshalValue(w, "t", t.Kind, t.Ident)
		fpl(w, "return nil")
		fpl(w, "}")
		return
	case gengo.IdentKind:
		if t.Ident == "T" {
			fpl(w, "// %s T", decl.Name)
		} else {
			fpl(w, "func (t %s) Marshal(w io.Writer) error {", decl.Name)
			marshalValue(w, "t", t.Kind, t.Ident)
			fpl(w, "return nil")
			fpl(w, "}")
		}
		return
	}
	fpl(w, "// type %s, %v", decl.Name, decl.Type)
}

func marshalBytes(w io.Writer, bytes string) {
	fpl(w, "if _, err := w.Write(%s); err != nil {", bytes)
	fpl(w, "return err")
	fpl(w, "}")
}

func marshalMarshaler(w io.Writer, marshaler string) {
	fpl(w, "if err := %s.Marshal(w); err != nil {", marshaler)
	fpl(w, "return err")
	fpl(w, "}")
}

func d0ToN(n int) string {
	ss := make([]string, n)
	for i := range ss {
		ss[i] = "b" + strconv.Itoa(i)
	}
	return strings.Join(ss, ", ")
}

func genFieldMarshal(w io.Writer, i int, f *gengo.Field) {
	fName := "t." + f.Name
	if f.Name == "" {
		fName = "t." + f.Type.Ident
	}
	marshalValue(w, fName, f.Type.Kind, f.Type.Ident)
}

func marshalValue(w io.Writer, name string, kind gengo.Kind, typ string) {
	switch kind {
	case gengo.IdentKind:
		switch typ {
		case "int64":
			marshalBytes(w, marshalBigEndian(name, 64))
		case "int32":
			marshalBytes(w, marshalBigEndian(name, 32))
		case "int16":
			marshalBytes(w, marshalBigEndian(name, 16))
		case "int8":
			marshalBytes(w, marshalBigEndian(name, 8))
		case "string":
			fpl(w, "{")
			fpl(w, "l := int16(len(%s))", name)
			marshalBytes(w, marshalBigEndian("l", 16))
			marshalBytes(w, fmt.Sprintf("[]byte(%s)", name))
			fpl(w, "}")
		case "[]byte":
			fpl(w, "{")
			fpl(w, "l := int32(len(%s))", name)
			marshalBytes(w, marshalBigEndian("l", 32))
			marshalBytes(w, fmt.Sprintf("%s", name))
			fpl(w, "}")
		default:
			marshalMarshaler(w, name)
		}
	case gengo.ArrayKind:
		fpl(w, "{")
		fpl(w, "l := int32(len(%s))", name)
		marshalBytes(w, marshalBigEndian("l", 32))
		fpl(w, "for i := range %s {", name)
		switch typ {
		case "int8", "int16", "int32", "int64", "string":
			marshalValue(w, name+"[i]", gengo.IdentKind, typ)
		default:
			marshalMarshaler(w, name+"[i]")
		}
		fpl(w, "}")
		fpl(w, "}")
	default:
		fpl(w, "// value %s %v", name, kind)
	}
}

func marshalBigEndian(name string, bit int) string {
	var buf bytes.Buffer
	buf.WriteString("[]byte{")
	if bit >= 64 {
		buf.WriteString(fmt.Sprintf("byte(%s>>56), ", name))
		buf.WriteString(fmt.Sprintf("byte(%s>>48), ", name))
		buf.WriteString(fmt.Sprintf("byte(%s>>32), ", name))
	}
	if bit >= 32 {
		buf.WriteString(fmt.Sprintf("byte(%s>>24), ", name))
		buf.WriteString(fmt.Sprintf("byte(%s>>16), ", name))
	}
	if bit >= 16 {
		buf.WriteString(fmt.Sprintf("byte(%s>>8), ", name))
	}
	buf.WriteString(fmt.Sprintf("byte(%s)}", name))
	return buf.String()
}
