package main

import (
	"html"
	"log"
	"os"
	"strings"

	"h12.me/gs"
	"h12.me/html-query"
	"h12.me/html-query/expr"
)

type Type struct {
	Name string
	Spec
}

type Spec interface{}

type BasicSpec string

const (
	String BasicSpec = "string"
	Int16  BasicSpec = "int16"
	Int32  BasicSpec = "int32"
	Int64  BasicSpec = "int64"
)

type ArraySpec struct {
	Types []*Type
}

type OrSpec struct {
	Types []*Type
}

var (
	class = expr.Class
	code  = expr.Code
	h4    = expr.H4
	h6    = expr.H6
	id    = expr.Id
)

func main() {
	root := gs.WebPage{}.Load("spec.html").Parse()

	f, err := os.Create("bnf.auto.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	root.Descendants(class("code")).For(func(item *query.Node) {
		//fmt.Println(*h.Text())
		bnf := cleanBNFCode(*item.Script().Text())
		f.Write([]byte(bnf))
		f.Write([]byte("\n\n"))
	})
}

func cleanBNFCode(code string) string {
	code = html.UnescapeString(code)
	code = strings.TrimSpace(code)
	code = strings.TrimPrefix(code, "<![CDATA[")
	code = strings.TrimSuffix(code, "]]>")
	code = strings.TrimSpace(code)
	ss := strings.Split(code, "\n")
	for i := range ss {
		ss[i] = strings.TrimPrefix(ss[i], "                                ")
	}
	return strings.Join(ss, "\n")
}

type typeMap map[string]*Type

func (m typeMap) ParseBNF(bnf string) *Type {
	parts := strings.Split(bnf, "=>")
	if len(parts) == 2 {
		name, spec := parts[0], parts[1]
		name = strings.TrimSpace(name)
		spec = strings.TrimSpace(spec)
		return &Type{Name: name, Spec: m.ParseSpec(spec)}
	}
	return nil
}

func (m typeMap) ParseSpec(spec string) Spec {
	return nil
}
