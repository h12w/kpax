package main

import (
	"bytes"
	"fmt"
	"io"

	"h12.me/gs"
	"h12.me/html-query"
	"h12.me/html-query/expr"
)

func fromHTMLToBNF(htmlFile string, w io.Writer) {
	root := gs.WebPage{}.Load(htmlFile).Parse()
	var (
		td    = expr.Td
		class = expr.Class
		code  = expr.Code
		not   = expr.Not
	)
	root.Descendants(td, class("code")).For(func(def *query.Node) {
		for i, line := range def.Descendants(code, not(class("spaces"))).All() {
			if i != 0 {
				w.Write([]byte{'\t'})
			}
			if txt := line.Text(); txt != nil {
				w.Write(bytes.TrimSpace([]byte(*txt)))
			} else {
				fmt.Sprintf("%#v", line)
			}
			w.Write([]byte{'\n'})
		}
		w.Write([]byte{'\n'})
	})
}
