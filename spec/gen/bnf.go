package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func parseBNF() []*Node {
	f, err := os.Open("bnf.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var nodes []*Node
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if line != "" {
			node := parseLine(line)
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func fromBNFToJSON() {
	nodes := parseBNF()
	fmt.Println(toJSON(nodes))
}

func fromBNFToGo() {
	nodes := parseBNF()
	nodeMap := make(map[string]*Node)
	for _, n := range nodes {
		nodeMap[n.Name] = n
	}
	for _, n := range nodes {
		if !n.simple() {
			n.GenStruct(os.Stdout, nodeMap)
		}
	}
}

func toJSON(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "", "    ")
	return string(buf)
}

func parseLine(line string) *Node {
	m := strings.Split(line, "=>")
	name, spec := strings.TrimSpace(m[0]), strings.TrimSpace(m[1])
	n := parseSpec(spec)
	n.Name = name
	return n
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
	return s.top()
}
