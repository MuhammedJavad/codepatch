package tree

import (
	"encoding/json"
	"errors"
	"time"
)

// builder provides a fluent interface for building trees
type builder struct {
	tree Tree
	err  error
}

// Builder creates a new tree builder
func Builder(name string, result json.RawMessage) *builder {
	return &builder{
		tree: Tree{
			Name:   name,
			Result: result,
			Active: true,
		},
	}
}

// WithStartTime sets the start time for the tree
func (tb *builder) WithStartTime(start time.Time) *builder {
	tb.tree.Start = &start
	return tb
}

// WithEndTime sets the end time for the tree
func (tb *builder) WithEndTime(end time.Time) *builder {
	tb.tree.End = &end
	return tb
}

// AddRoot adds a root node using a builder function (legacy method)
func (tb *builder) WithRoot(builder func(*NodeBuilder)) *builder {
	nb := &NodeBuilder{}
	builder(nb)
	root, err := nb.Build()
	if err != nil {
		tb.err = err
		return tb
	}
	tb.tree.Root = root
	return tb
}

// Build validates and returns the built tree
func (tb *builder) Build() (*Tree, error) {
	if tb.err != nil {
		return nil, tb.err
	}

	if tb.tree.Root.isEmpty() {
		return nil, errors.New("root node is required")
	}

	// Validate time range
	if tb.tree.Start != nil && tb.tree.End != nil && tb.tree.Start.After(*tb.tree.End) {
		return nil, errors.New("start time cannot be after end time")
	}

	return &tb.tree, nil
}

// NodeBuilder provides a fluent interface for building nodes
type NodeBuilder struct {
	node Node
	err  error
}

// AsCondition creates a condition node
func (nb *NodeBuilder) AsCondition(operator, operand string) *NodeBuilder {
	if operator == "" {
		nb.err = errors.New("operator is required for condition")
		return nb
	}

	nb.node = Node{
		NodeType: NodeTypeCondition,
		Operator: operator,
		Operand:  operand,
		Children: nil, // conditions have no children
	}
	return nb
}

var validGates = map[string]bool{
	"and":  true,
	"or":   true,
	"xor":  true,
	"xnor": true,
	"nand": true,
	"nor":  true,
}

// AsGate creates a gate node with children
// Gate creates a gate node with children (legacy method)
func (nb *NodeBuilder) AsGate(g string, builder func(*GateBuilder)) *NodeBuilder {
	if !validGates[g] {
		nb.err = errors.New("invalid gate operator")
		return nb
	}
	gb := &GateBuilder{
		gate: g,
	}
	builder(gb)
	node, err := gb.Build()
	if err != nil {
		nb.err = err
		return nb
	}
	nb.node = node
	return nb
}

// Build returns the built node
func (nb *NodeBuilder) Build() (Node, error) {
	if nb.err != nil {
		return Node{}, nb.err
	}

	if nb.node.NodeType == "" {
		return Node{}, errors.New("node type must be specified. select either condition or gate")
	}

	return nb.node, nil
}

// GateBuilder handles building gates with children
type GateBuilder struct {
	gate        string
	children    []Node
	internalErr error
}

// AddChild adds a child node to the gate
func (gb *GateBuilder) AddCondition(operator, operand string) {
	nb := &NodeBuilder{}
	nb.AsCondition(operator, operand)
	node, err := nb.Build()
	if err != nil {
		gb.internalErr = err
		return
	}
	gb.children = append(gb.children, node)
}

func (gb *GateBuilder) AddGate(g string, builder func(*GateBuilder)) {
	nb := &NodeBuilder{}
	nb.AsGate(g, builder)
	node, err := nb.Build()
	if err != nil {
		gb.internalErr = err
		return
	}
	gb.children = append(gb.children, node)
}

// Build creates the gate node
func (gb *GateBuilder) Build() (Node, error) {
	if gb.internalErr != nil {
		return Node{}, gb.internalErr
	}
	if gb.gate == "" {
		return Node{}, errors.New("gate operator is required")
	}
	if len(gb.children) < 2 {
		return Node{}, errors.New("gate must have at least 2 children")
	}
	return Node{
		NodeType: NodeTypeGate,
		Operator: gb.gate,
		Operand:  "",
		Children: gb.children,
	}, nil
}
