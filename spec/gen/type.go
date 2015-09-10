package main

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
