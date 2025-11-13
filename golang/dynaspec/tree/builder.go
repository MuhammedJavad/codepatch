package tree

import (
	"errors"
	"fmt"
	"time"
)

// builder provides a fluent interface for building trees
type builder struct {
	tree Tree
	err  error
}
type TreeModel struct {
	Result    float32    `json:"result"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Root      NodeModel  `json:"root"`
}
type NodeModel struct {
	Operator string      `json:"operator"`
	Operand  string      `json:"operand"`
	Children []NodeModel `json:"children"`
}

func BuildFromModel(model TreeModel) (Tree, error) {
	return Builder(model.Result).
		WithStartTime(model.StartTime).
		WithEndTime(model.EndTime).
		WithRoot(func(nb *NodeBuilder) {
			buildRoot(nb, model.Root)
		}).
		Build()
}

func buildRoot(nb *NodeBuilder, n NodeModel) {
	if len(n.Children) <= 0 {
		nb.AsCondition(n.Operator, n.Operand)
		return
	}

	nb.AsGate(n.Operator, func(gb *GateBuilder) {
		for _, child := range n.Children {
			buildChild(gb, child)
		}
	})
}

func buildChild(gb *GateBuilder, c NodeModel) {
	if len(c.Children) <= 0 {
		gb.AddCondition(c.Operator, c.Operand)
		return
	}

	gb.AddGate(c.Operator, func(childGb *GateBuilder) {
		for _, grandChild := range c.Children {
			buildChild(childGb, grandChild)
		}
	})
}

// Builder creates a new tree builder
func Builder(result float32) *builder {
	return &builder{
		tree: Tree{
			Active: true,
			Result: result,
		},
	}
}

// WithStartTime sets the start time for the tree
func (tb *builder) WithStartTime(start *time.Time) *builder {
	if start == nil {
		return tb
	}
	tb.tree.Start = start
	return tb
}

// WithEndTime sets the end time for the tree
func (tb *builder) WithEndTime(end *time.Time) *builder {
	if end == nil {
		return tb
	}
	tb.tree.End = end
	return tb
}

// AddRoot adds a root node using a builder function (legacy method)
func (tb *builder) WithRoot(builder func(*NodeBuilder)) *builder {
	nb := &NodeBuilder{}
	builder(nb)
	root, err := nb.build()
	if err != nil {
		tb.err = err
		return tb
	}
	tb.tree.Root = root
	return tb
}

// Build validates and returns the built tree
func (tb *builder) Build() (Tree, error) {
	if tb.err != nil {
		return Tree{}, tb.err
	}

	if tb.tree.Root.IsEmpty() {
		return Tree{}, errors.New("root node is required")
	}

	// Validate time range
	if tb.tree.Start != nil && tb.tree.End != nil && tb.tree.Start.After(*tb.tree.End) {
		return Tree{}, errors.New("start time cannot be after end time")
	}

	if tb.tree.End != nil {
		tb.tree.Active = true
		// first check the end time
		// it determines that tree is active or not
	} else if tb.tree.Start != nil {
		tb.tree.Active = false
		// having start time means that tree is not active yet
		// an scheduler should check the date periodically
		// and update the active status
	} else {
		tb.tree.Active = true
		// no start time and end time means that tree is active fir ever
	}

	return tb.tree, nil
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

	if !reg.IsCondition(operator) {
		nb.err = fmt.Errorf("invalid condition operator: %s. Valid conditions are: %v", operator, reg.ConditionOperators())
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

// AsGate creates a gate node with children
// Gate creates a gate node with children (legacy method)
func (nb *NodeBuilder) AsGate(g string, builder func(*GateBuilder)) *NodeBuilder {
	if !reg.IsGate(g) {
		nb.err = fmt.Errorf("invalid gate operator: %s. Valid gates are: %v", g, reg.GateOperators())
		return nb
	}
	gb := &GateBuilder{
		gate: g,
	}
	builder(gb)
	node, err := gb.build()
	if err != nil {
		nb.err = err
		return nb
	}
	nb.node = node
	return nb
}

// build returns the built node
func (nb *NodeBuilder) build() (Node, error) {
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
	node, err := nb.build()
	if err != nil {
		gb.internalErr = err
		return
	}
	gb.children = append(gb.children, node)
}

func (gb *GateBuilder) AddGate(g string, builder func(*GateBuilder)) {
	nb := &NodeBuilder{}
	nb.AsGate(g, builder)
	node, err := nb.build()
	if err != nil {
		gb.internalErr = err
		return
	}
	gb.children = append(gb.children, node)
}

// Build creates the gate node
func (gb *GateBuilder) build() (Node, error) {
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
