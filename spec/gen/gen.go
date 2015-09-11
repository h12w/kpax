package main

import (
	"fmt"
	"io"
	"strings"
)

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

func (n *Node) GenType(w io.Writer, m map[string]*Decl) {
	switch n.NodeType {
	case LeafNode:
		if decl, ok := m[n.Value]; ok && decl.Type.simple() {
			decl.Type.GenType(w, m)
		} else {
			fp(w, goType(n.Value))
		}
	case SeqNode:
		if len(n.Child) == 1 && n.Child[0].NodeType == LeafNode {
			n.Child[0].GenType(w, m)
		} else {
			fpl(w, "struct {")
			for _, c := range n.Child {
				c.GenField(w, m)
			}
			fp(w, "}")
		}
	case ZeroOrMoreNode:
		fp(w, "[]")
		if len(n.Child) == 1 {
			n.Child[0].GenType(w, m)
		} else {
			(&Node{NodeType: SeqNode, Child: n.Child}).GenType(w, m)
		}
	case OrNode:
		fp(w, "T")
	default:
		fp(w, "-")
	}
}

func (d *Decl) GenDecl(w io.Writer, m map[string]*Decl) {
	fp(w, "type %s ", d.Name)
	d.Type.GenType(w, m)
	fpl(w, "")
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

func (n *Node) GenField(w io.Writer, m map[string]*Decl) {
	var typ *Decl
	if decl, ok := m[n.Value]; ok {
		typ = decl
	}
	name := n.Value
	if name == "" && n.NodeType == ZeroOrMoreNode && len(n.Child) == 1 {
		name = n.Child[0].Value + "s"
	}
	fp(w, goName(name)+" ")
	if typ != nil {
		fp(w, typ.typeName())
	} else {
		n.GenType(w, m)
	}
	fpl(w, "")
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
