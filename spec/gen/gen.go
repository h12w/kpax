package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"os"
	"strings"

	"h12.me/gs"
	"h12.me/html-query"
	"h12.me/html-query/expr"
)

type NodeType int

const (
	LeafNode NodeType = iota
	SeqNode
	OrNode
	ZeroOrMoreNode
)

func (t NodeType) MarshalText() ([]byte, error) {
	switch t {
	case LeafNode:
		return []byte(""), nil
	case SeqNode:
		return []byte(" "), nil
	case OrNode:
		return []byte("|"), nil
	case ZeroOrMoreNode:
		return []byte("*"), nil
	}
	return nil, nil
}

type Node struct {
	Name  string   `json:"name,omitempty"`
	Type  NodeType `json:"type,omitempty"`
	Child []*Node  `json:"child,omitempty"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("go run gen.go (bnf | json)")
		return
	}
	switch os.Args[1] {
	case "bnf":
		genBNF()
	case "json":
		genJSON()
	}
}

func genJSON() {
	f, err := os.Open("bnf.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if line != "" {
			fmt.Println(toJSON(parseLine(line)))
		}
	}
}

func toJSON(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "", "    ")
	return string(buf)
}

func parseLine(line string) Node {
	m := strings.Split(line, "=>")
	name, spec := strings.TrimSpace(m[0]), strings.TrimSpace(m[1])
	return Node{
		Name:  name,
		Type:  SeqNode,
		Child: parseSpec(spec),
	}
}

func tokenize(def string) []string {
	var tokens []string
	var token []rune
	for _, c := range def {
		switch c {
		case '(', ')', '[', ']', '|', ' ':
			if len(token) > 0 {
				tokens = append(tokens, string(token))
				token = nil
			}
			if c != ' ' {
				tokens = append(tokens, string(c))
			}
		default:
			token = append(token, c)
		}
	}
	if len(token) > 0 {
		tokens = append(tokens, string(token))
	}
	return tokens
}

type nodeStack struct {
	a []*Node
}

func (s *nodeStack) top() *Node {
	return s.a[len(s.a)-1]
}

func (s *nodeStack) push(n *Node) {
	s.a = append(s.a, n)
}

func (s *nodeStack) count() int {
	return len(s.a)
}

func (s *nodeStack) pop() (n *Node) {
	u := len(s.a) - 1
	s.a, n = s.a[:u], s.a[u]
	return
}

func parseSpec(def string) (nodes []*Node) {
	tokens := tokenize(def)
	s := nodeStack{}
	s.push(&Node{Type: SeqNode})
	for _, token := range tokens {
		top := s.top()
		switch token {
		case "(":
			n := &Node{Type: SeqNode}
			top.Child = append(top.Child, n)
			s.push(n)
		case ")":
			s.pop()
		case "[":
			n := &Node{Type: ZeroOrMoreNode}
			top.Child = append(top.Child, n)
			s.push(n)
		case "]":
			s.pop()
		case "|":
			top.Type = OrNode
		default:
			top.Child = append(top.Child, &Node{Name: token, Type: LeafNode})
		}
	}
	if s.count() == 0 {
		fmt.Println(toJSON(tokens))
	}
	return s.top().Child
}

func genBNF() {
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
