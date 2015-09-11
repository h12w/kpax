package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func parseBNF() []*Decl {
	f, err := os.Open("bnf.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var decls []*Decl
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if line != "" {
			decls = append(decls, parseLine(line))
		}
	}
	return decls
}

func fromBNFToJSON() {
	nodes := parseBNF()
	fmt.Println(toJSON(nodes))
}

func fromBNFToGoJSON() {
	goFile := fromBNFToGoFile()
	fmt.Println(goFile.JSON())
}

func fromBNFToGo() {
	goFile := fromBNFToGoFile()
	goFile.Fprint(os.Stdout)
}

func toJSON(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "", "    ")
	return string(buf)
}

func parseLine(line string) *Decl {
	m := strings.Split(line, "=>")
	name, spec := strings.TrimSpace(m[0]), strings.TrimSpace(m[1])
	n := parseSpec(spec)
	return &Decl{Name: name, Type: n}
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

func parseSpec(def string) (node *Node) {
	tokens := tokenize(def)
	s := nodeStack{}
	s.push(&Node{NodeType: SeqNode})
	for _, token := range tokens {
		top := s.top()
		switch token {
		case "(":
			n := &Node{NodeType: SeqNode}
			top.Child = append(top.Child, n)
			s.push(n)
		case ")":
			s.pop()
		case "[":
			n := &Node{NodeType: ZeroOrMoreNode}
			top.Child = append(top.Child, n)
			s.push(n)
		case "]":
			s.pop()
		case "|":
			top.NodeType = OrNode
		default:
			top.Child = append(top.Child, &Node{Value: token, NodeType: LeafNode})
		}
	}
	if s.count() == 0 {
		fmt.Println(toJSON(tokens))
	}
	top := s.top()
	if len(top.Child) == 1 && top.Child[0].Value == "" {
		top.Child[0].Value = top.Value
		return top.Child[0]
	}
	return s.top()
}
