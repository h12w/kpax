package main

import (
	"fmt"
	"io"
	"os"

	"h12.me/gengo"
)

func fromBNFToGoFuncs() {
	w := os.Stdout
	goFile := fromBNFToGoFile()
	fpl(w, "package proto")
	fpl(w, "import (")
	fpl(w, `"hash/crc32"`)
	fpl(w, ")")
	for _, decl := range goFile.TypeDecls {
		genMarshalFunc(w, decl)
		fpl(w, "")
		genUnmarshalFunc(w, decl)
		fpl(w, "")
	}
}

func genMarshalFunc(w io.Writer, decl *gengo.TypeDecl) {
	t := decl.Type
	if t.Kind == gengo.IdentKind && t.Ident == "T" {
		return
	}
	fpl(w, "func (t *%s) Marshal(w *Writer) {", decl.Name)
	switch t.Kind {
	case gengo.StructKind:
		if f0 := t.Fields[0]; len(t.Fields) == 2 &&
			(f0.Name == "Size" || f0.Name == "CRC") {
			fpl(w, "offset := len(w.B)")
			marshalField(w, t.Fields[0])
			fpl(w, "start := len(w.B)")
			marshalField(w, t.Fields[1])
			switch f0.Name {
			case "Size":
				fpl(w, "w.SetInt32(offset, int32(len(w.B)-start))")
			case "CRC":
				fpl(w, "w.SetUint32(offset, crc32.ChecksumIEEE(w.B[start:]))")
			}
		} else {
			for _, f := range t.Fields {
				marshalField(w, f)
			}
		}
	case gengo.ArrayKind:
		marshalValue(w, "(*t)", t.Kind, t.Ident)
	case gengo.IdentKind:
		marshalValue(w, "(*t)", t.Kind, t.Ident)
	default:
		fpl(w, "// type %s, %v", decl.Name, decl.Type)
	}
	fpl(w, "}")
}

func genUnmarshalFunc(w io.Writer, decl *gengo.TypeDecl) {
	t := decl.Type
	if t.Kind == gengo.IdentKind && t.Ident == "T" {
		return
	}
	fpl(w, "func (t *%s) Unmarshal(r *Reader) {", decl.Name)
	switch t.Kind {
	case gengo.StructKind:
		for _, f := range t.Fields {
			fName := "t." + f.Name
			if f.Name == "" {
				fName = "t." + f.Type.Ident
			}
			unmarshalValue(w, fName, f.Type.Kind, f.Type.Ident, f.Type.Ident)
		}
	case gengo.ArrayKind:
		unmarshalValue(w, "(*t)", t.Kind, t.Ident, t.Ident)
	case gengo.IdentKind:
		unmarshalValue(w, "(*t)", t.Kind, t.Ident, decl.Name)
	default:
		fpl(w, "// type %s, %v", decl.Name, decl.Type)
	}
	fpl(w, "}")
}

func marshalField(w io.Writer, f *gengo.Field) {
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
			marshalInt(w, name, 64)
		case "int32":
			marshalInt(w, name, 32)
		case "int16":
			marshalInt(w, name, 16)
		case "int8":
			marshalInt(w, name, 8)
		case "string":
			fpl(w, "w.WriteString(string(%s))", name)
		case "[]byte":
			fpl(w, "w.WriteBytes(%s)", name)
		default:
			marshalMarshaler(w, name)
		}
	case gengo.ArrayKind:
		marshalInt(w, fmt.Sprintf("int32(len(%s))", name), 32)
		fpl(w, "for i := range %s {", name)
		switch typ {
		case "int8", "int16", "int32", "int64", "string":
			marshalValue(w, name+"[i]", gengo.IdentKind, typ)
		default:
			marshalMarshaler(w, name+"[i]")
		}
		fpl(w, "}")
	default:
		fpl(w, "// value %s %v", name, kind)
	}
}

func unmarshalValue(w io.Writer, name string, kind gengo.Kind, typ, declType string) {
	switch kind {
	case gengo.IdentKind:
		switch typ {
		case "int64":
			unmarshalInt(w, name, "b", 64)
		case "int32":
			unmarshalInt(w, name, "b", 32)
		case "int16":
			unmarshalInt(w, name, "b", 16)
		case "int8":
			unmarshalInt(w, name, "b", 8)
		case "string":
			fpl(w, "%s = %s(r.ReadString())", name, declType)
		case "[]byte":
			fpl(w, "%s = r.ReadBytes()", name)
		default:
			unmarshalUnmarshaler(w, name, typ)
		}
	case gengo.ArrayKind:
		fpl(w, "%s = make([]%s, int(r.ReadInt32()))", name, typ)
		fpl(w, "for i := range %s {", name)
		switch typ {
		case "int8", "int16", "int32", "int64", "string":
			unmarshalValue(w, name+"[i]", gengo.IdentKind, typ, typ)
		default:
			unmarshalUnmarshaler(w, name+"[i]", typ)
		}
		fpl(w, "}")
	default:
		fpl(w, "// value %s %v", name, kind)
	}
}

func marshalMarshaler(w io.Writer, marshaler string) {
	fpl(w, "%s.Marshal(w)", marshaler)
}

func unmarshalUnmarshaler(w io.Writer, unmarshaler, typ string) {
	fpl(w, "%s.Unmarshal(r)", unmarshaler)
}

func marshalInt(w io.Writer, name string, bit int) {
	fpl(w, "w.WriteInt%d(%s)", bit, name)
}

func unmarshalInt(w io.Writer, name string, bufName string, bit int) {
	fpl(w, "%s=r.ReadInt%d()", name, bit)
}

func writeBytes(w io.Writer, bytes string) {
	fpl(w, "if _, err := w.Write(%s); err != nil {", bytes)
	fpl(w, "return err")
	fpl(w, "}")
}

func readBytes(w io.Writer, name string, l string) {
	fpl(w, "%s := make([]byte, %s)", name, l)
	fpl(w, "if _, err := r.Read(%s); err != nil {", name)
	fpl(w, "return err")
	fpl(w, "}")
}

func readByteArray(w io.Writer, name string, bit int) {
	fpl(w, "var %s [%d]byte", name, bit/8)
	fpl(w, "if _, err := r.Read(%s[:]); err != nil {", name)
	fpl(w, "return err")
	fpl(w, "}")
}

func intToBytes(name string, bit int) string {
	switch bit {
	case 64:
		return fmt.Sprintf("[]byte{byte(%[1]s>>56), byte(%[1]s>>48), byte(%[1]s>>40), byte(%[1]s)>>32, "+
			"byte(%[1]s>>24), byte(%[1]s>>16), byte(%[1]s>>8), byte(%[1]s)}", name)
	case 32:
		return fmt.Sprintf("[]byte{byte(%[1]s>>24), byte(%[1]s>>16), byte(%[1]s>>8), byte(%[1]s)}", name)
	case 16:
		return fmt.Sprintf("[]byte{byte(%[1]s>>8), byte(%[1]s)}", name)
	case 8:
		return fmt.Sprintf("[]byte{byte(%[1]s)}", name)
	}
	return ""
}

func bytesToInt(name string, bit int) string {
	switch bit {
	case 64:
		return fmt.Sprintf("int64(%[1]s[0])<<56 | int64(%[1]s[1])<<48 | int64(%[1]s[2])<<40 | int64(%[1]s[3])<<32 | "+
			"int64(%[1]s[4])<<24 | int64(%[1]s[5])<<16 | int64(%[1]s[6])<<8 | int64(%[1]s[7])", name)
	case 32:
		return fmt.Sprintf("int32(%[1]s[0])<<24 | int32(%[1]s[1])<<16 | int32(%[1]s[2])<<8 | int32(%[1]s[3])", name)
	case 16:
		return fmt.Sprintf("int16(%[1]s[0])<<8 | int16(%[1]s[1])", name)
	case 8:
		return fmt.Sprintf("int8(%[1]s[0])", name)
	}
	return ""
}
