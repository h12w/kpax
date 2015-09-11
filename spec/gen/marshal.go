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
	fpl(w, `"bytes"`)
	fpl(w, ")")
	for _, decl := range goFile.TypeDecls {
		genTypeMarshal(w, decl)
		fpl(w, "")
	}
}

func genTypeMarshal(w io.Writer, decl *gengo.TypeDecl) {
	if decl.Type.Kind == gengo.StructKind {
		fpl(w, "func (t *%s) Marshal() []byte {", decl.Name)
		for i, f := range decl.Type.Fields {
			genFieldMarshal(w, i, f)
		}
		fpl(w, "return bytes.Join([][]byte{%s}, nil)", d0ToN(len(decl.Type.Fields)))
		fpl(w, "}")
	} else {
		fpl(w, "// type %s", decl.Name)
	}
}

func d0ToN(n int) string {
	ss := make([]string, n)
	for i := range ss {
		ss[i] = "d" + strconv.Itoa(i)
	}
	return strings.Join(ss, ", ")
}

func genFieldMarshal(w io.Writer, i int, f *gengo.Field) {
	switch f.Type.Kind {
	case gengo.IdentKind:
		switch f.Type.Ident {
		case "int64":
			fpl(w, "d%d := %s", i, marshalBigEndian("t."+f.Name, 64))
		case "int32":
			fpl(w, "d%d := %s", i, marshalBigEndian("t."+f.Name, 32))
		case "int16":
			fpl(w, "d%d := %s", i, marshalBigEndian("t."+f.Name, 16))
		case "int8":
			fpl(w, "d%d := %s", i, marshalBigEndian("t."+f.Name, 8))
		case "string":
			fpl(w, "d%dLen := strlen(t.%s)", i, f.Name)
			fpl(w, "d%d := append(%s, []byte(t.%s)...)", i, marshalBigEndian(fmt.Sprintf("d%dLen", i), 16), f.Name)
		default:
			fpl(w, "d%d := t.Marshal()", i)
		}
	default:
		fpl(w, "// field %v", f.Type.Kind)
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
