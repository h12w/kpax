package main

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"h12.me/gengo"
)

func fromBNFToGoFile() *gengo.File {
	decls := parseBNF()
	declMap := make(map[string]*Decl)
	for _, decl := range decls {
		if foundDecl, ok := declMap[decl.Name]; ok && !reflect.DeepEqual(foundDecl, decl) {
			log.Fatalf("conflict name %s:\n%#v\n%#v", decl.Name, foundDecl, decl)
		}
		declMap[decl.Name] = decl
	}
	return genGoFile(decls, declMap)
}

func genGoFile(decls []*Decl, declMap map[string]*Decl) *gengo.File {
	var goDecls []*gengo.TypeDecl
	for _, decl := range decls {
		if !decl.Type.simple() {
			goDecls = append(goDecls, decl.GenDecl(declMap))
		}
	}
	return &gengo.File{
		PackageName: "proto",
		TypeDecls:   goDecls,
	}
}

func (n *Node) simple() bool {
	if len(n.Child) != 1 || n.Child[0].NodeType != LeafNode {
		return false
	}
	switch n.Child[0].Value {
	case "string", "int8", "int16", "int32", "int64", "bytes":
		return true
	}
	return false
}

func (n *Node) GenType(m map[string]*Decl) *gengo.Type {
	switch n.NodeType {
	case LeafNode:
		if decl, ok := m[n.Value]; ok && decl.Type.simple() {
			return decl.Type.GenType(m)
		} else {
			return &gengo.Type{Kind: gengo.IdentKind, Ident: goType(n.Value)}
		}
	case SeqNode:
		if len(n.Child) == 1 && n.Child[0].NodeType == LeafNode {
			return n.Child[0].GenType(m)
		} else {
			var fields []*gengo.Field
			for _, c := range n.Child {
				fields = append(fields, c.GenField(m))
			}
			return &gengo.Type{
				Kind:   gengo.StructKind,
				Fields: fields,
			}
		}
	case ZeroOrMoreNode:
		var t *gengo.Type
		if len(n.Child) == 1 {
			t = n.Child[0].GenType(m)
		} else {
			t = (&Node{NodeType: SeqNode, Child: n.Child}).GenType(m)
		}
		t.Kind = gengo.ArrayKind
		return t
	case OrNode:
		return &gengo.Type{
			Kind:  gengo.IdentKind,
			Ident: "T",
		}
	default:
		return &gengo.Type{
			Kind:  gengo.IdentKind,
			Ident: "-",
		}
	}
}

func (d *Decl) GenDecl(m map[string]*Decl) *gengo.TypeDecl {
	ut := d.Type.Value
	for ut != "" {
		if t, ok := m[ut]; ok && t.Type.Value != "" {
			ut = t.Type.Value
		} else {
			break
		}
	}
	return &gengo.TypeDecl{
		Name:           goName(d.Name),
		Type:           *d.Type.GenType(m),
		UnderlyingType: ut,
	}
}

func (d *Decl) typeName() string {
	typ := ""
	if d.Type.simple() {
		typ = d.Type.Child[0].Value
	} else {
		typ = d.Name
	}
	return goType(typ)
}

func (n *Node) GenField(m map[string]*Decl) *gengo.Field {
	name := n.Value
	decl, _ := m[name]
	if name == "" && n.NodeType == ZeroOrMoreNode && len(n.Child) == 1 {
		name = n.Child[0].Value + "s"
	}
	if decl != nil && n.NodeType == ZeroOrMoreNode {
		goType := decl.Type.GenType(m)
		if n.NodeType == ZeroOrMoreNode {
			goType.Kind = gengo.ArrayKind
		}
		return &gengo.Field{
			Name: goName(name),
			Type: *goType,
		}
	}
	return &gengo.Field{
		Name: goName(name),
		Type: *n.GenType(m),
	}
}

func fp(w io.Writer, format string, v ...interface{}) {
	fmt.Fprintf(w, format, v...)
}

func fpl(w io.Writer, format string, v ...interface{}) {
	fmt.Fprintf(w, format, v...)
	fmt.Fprintln(w)
}

func goName(s string) string {
	if strings.HasSuffix(s, "Id") {
		return strings.TrimSuffix(s, "Id") + "ID"
	}
	s = strings.Replace(s, "Crc", "CRC", -1)
	s = strings.Replace(s, "Api", "API", -1)
	s = strings.Replace(s, "Isr", "ISR", -1)
	return s
}

func goType(s string) string {
	switch s {
	case "bytes":
		return "[]byte"
	}
	return goName(s)
}
