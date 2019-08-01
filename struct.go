package xml2json

import (
	"strings"
)

// Node is a data element on a tree
type Node struct {
	Children              map[string]Nodes
	Data                  string
	ChildrenAlwaysAsArray bool
}

// Nodes is a list of nodes
type Nodes []*Node

// AddChild appends a node to the list of children
func (n *Node) AddChild(s string, c *Node) {
	// Lazy lazy
	if n.Children == nil {
		n.Children = map[string]Nodes{}
	}

	n.Children[s] = append(n.Children[s], c)
}

// IsComplex returns whether it is a complex type (has children)
func (n *Node) IsComplex() bool {
	return len(n.Children) > 0
}

// GetChild returns child by path if exists. Path looks like "grandparent.parent.child.grandchild"
func (n *Node) GetChild(path string) *Node {
	result := n
	names := strings.Split(path, ".")
	for _, name := range names {
		children, exists := result.Children[name]
		if !exists {
			return nil
		}
		if len(children) == 0 {
			return nil
		}
		result = children[0]
	}
	return result
}

// GetChildren returns a slice of Children for a given path
func (n *Node) GetChildren(path string) []*Node {
	pathtok := strings.Split(path, ".")
	return n.getChildren(pathtok)
}

func (n *Node) getChildren(pt []string) []*Node {
	var result []*Node
	name := pt[0]

	children, exists := n.Children[name]
	if !exists {
		return nil
	}
	if len(children) == 0 {
		return nil
	}
	if len(pt) > 1 {
		for _, child := range children {
			result = append(result, child.getChildren(pt[1:])...)
		}
	} else {
		result = append(result, children...)
	}

	return result
}
