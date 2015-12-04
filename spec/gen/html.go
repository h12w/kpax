package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"h12.me/gs"
	"h12.me/html-query"
	"h12.me/html-query/expr"
)

var (
	td    = expr.Td
	tr    = expr.Tr
	class = expr.Class
	code  = expr.Code
	not   = expr.Not
	id    = expr.Id
	tbody = expr.Tbody
)

func fromHTMLToBNF(htmlFile string, w io.Writer) {
	root := gs.WebPage{}.Load(htmlFile).Parse()
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

type errorInfo struct {
	name string
	code string
	msg  string
}

func genErrorCodes(htmlFile string, w io.Writer) {
	root := gs.WebPage{}.Load(htmlFile).Parse()
	table := root.Find(id("AGuideToTheKafkaProtocol-ErrorCodes")).FindNext(class("table-wrap")).Find(tbody)
	var errors []errorInfo
	table.Children(tr).For(func(row *query.Node) {
		cols := row.Children(td).All()
		errors = append(errors, errorInfo{
			name: *cols[0].PlainText(),
			code: *cols[1].PlainText(),
			msg:  cleanMsg(*cols[2].PlainText()),
		})
	})
	fpl(w, "package proto")
	fpl(w, "const (")
	for _, e := range errors {
		if e.code != "0" {
			fp(w, "Err")
		}
		fp(w, e.name)
		fp(w, " ErrorCode = ")
		fpl(w, e.code)
	}
	fpl(w, ")")

	fpl(w, "var errTexts = []string{")
	for _, e := range errors {
		if e.code != "-1" {
			fp(w, e.code)
			fp(w, ":")
			fp(w, `"proto(`)
			fp(w, e.code)
			fp(w, "): ")
			fp(w, e.msg)
			fpl(w, `",`)
		}
	}
	fpl(w, "}")
}

func cleanMsg(s string) string {
	s = strings.TrimSuffix(s, ".")
	ss := strings.Split(s, " ")
	if len(ss) > 0 {
		ss[0] = strings.ToLower(ss[0])
	}
	return strings.Join(ss, " ")
}

func fp(w io.Writer, format string, v ...interface{}) {
	fmt.Fprintf(w, format, v...)
}

func fpl(w io.Writer, format string, v ...interface{}) {
	fmt.Fprintf(w, format, v...)
	fmt.Fprintln(w)
}
