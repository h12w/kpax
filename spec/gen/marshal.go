package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"h12.me/gengo"
)

func fromBNFToGoFuncs() {
	w := os.Stdout
	goFile := fromBNFToGoFile()
	fpl(w, "package proto")
	for _, decl := range goFile.TypeDecls {
		genTypeMarshal(w, decl)
		fpl(w, "")
	}
}

func genTypeMarshal(w io.Writer, decl *gengo.TypeDecl) {
	if decl.Type.Kind == gengo.StructKind {
		fpl(w, "func (t *%s) MarshalBinary() (data []byte, err error) {", decl.Name)
		for i, f := range decl.Type.Fields {
			genFieldMarshal(w, i, f)
		}
		fpl(w, "}")
	} else {
		fpl(w, "// type %s", decl.Name)
	}
}

func genFieldMarshal(w io.Writer, i int, f *gengo.Field) {
	switch f.Type.Kind {
	case gengo.IdentKind:
		switch f.Type.Ident {
		case "int64":
			fpl(w, "d%d := %s", i, marshalBigEndian(f.Name, 64))
		case "int32":
			fpl(w, "d%d := %s", i, marshalBigEndian(f.Name, 32))
		case "int16":
			fpl(w, "d%d := %s", i, marshalBigEndian(f.Name, 16))
		case "int8":
			fpl(w, "d%d := %s", i, marshalBigEndian(f.Name, 8))
		}
	default:
		fpl(w, "// field %v", f.Type.Kind)
	}
}

func marshalBigEndian(name string, bit int) string {
	var buf bytes.Buffer
	buf.WriteString("[]byte{")
	if bit >= 64 {
		buf.WriteString(fmt.Sprintf("byte(t.%s>>56), ", name))
		buf.WriteString(fmt.Sprintf("byte(t.%s>>48), ", name))
		buf.WriteString(fmt.Sprintf("byte(t.%s>>32), ", name))
	}
	if bit >= 32 {
		buf.WriteString(fmt.Sprintf("byte(t.%s>>24), ", name))
		buf.WriteString(fmt.Sprintf("byte(t.%s>>16), ", name))
	}
	if bit >= 16 {
		buf.WriteString(fmt.Sprintf("byte(t.%s>>8), ", name))
	}
	buf.WriteString(fmt.Sprintf("byte(t.%s)}", name))
	return buf.String()
}
