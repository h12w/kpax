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

func fromHTMLToBNF() {
	root := gs.WebPage{}.Load("spec.html").Parse()

	f, err := os.Create("bnf.auto.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var (
		class = expr.Class
	)

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
