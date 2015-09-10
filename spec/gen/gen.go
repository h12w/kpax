package main

import (
	"fmt"
	"io"
	"strings"
)

func (n *Node) simple() bool {
	if len(n.Child) != 1 || n.Child[0].Type != LeafNode {
		return false
	}
	switch n.Child[0].Name {
	case "string", "int16", "int32", "int64":
		return true
	}
	return false
}

func (n *Node) GenStruct(w io.Writer, m map[string]*Node) {
	if n.Type == SeqNode {
		fp(w, "type %s struct {", n.Name)
		for _, c := range n.Child {
			c.GenField(w, m)
		}
		fp(w, "}")
	} else {
		fp(w, "type %s -", n.Name)
	}
}

func (n *Node) GenField(w io.Writer, m map[string]*Node) {
	name := n.Name
	// if name == "" {
	// 	if n.Type == OrNode {
	// 		ss := make([]string, len(n.Child))
	// 		for i := range ss {
	// 			ss[i] = n.Child[i].Name
	// 		}
	// 		name = strings.Join(ss, "Or")
	// 	}
	// }
	if name == "" {
		name = "-"
	}
	typ := name
	if typNode, ok := m[n.Name]; ok {
		if len(typNode.Child) == 1 && typNode.Child[0].Type == LeafNode {
			typ = typNode.Child[0].Name
		}
	}
	fp(w, "%s %s", goName(name), typ)
}

func fp(w io.Writer, format string, v ...interface{}) {
	fmt.Fprintf(w, format, v...)
	fmt.Fprintln(w)
}

func goName(s string) string {
	if strings.HasSuffix(s, "Id") {
		return strings.TrimSuffix(s, "Id") + "ID"
	}
	return s
}
