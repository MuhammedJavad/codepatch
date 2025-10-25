package tree

type NodeType string

const (
	NodeTypeGate      = "gate"
	NodeTypeCondition = "condition"
)

type Node struct {
	Operand  string   `json:"operand"`
	NodeType NodeType `json:"type"`
	Operator string   `json:"operator"`
	Children []Node   `json:"children"`
}

// NewNode creates a new Node
func NewNode(operand string, nodeType NodeType, operator string, children []Node) Node {
	return Node{
		Operand:  operand,
		NodeType: nodeType,
		Operator: operator,
		Children: children,
	}
}

func (n Node) Next() func() (Node, bool) {
	visited := -1

	return func() (Node, bool) {
		if len(n.Children) == 0 {
			return Node{}, false
		}
		// increase and check children length
		visited++
		if visited >= len(n.Children) {
			return Node{}, false
		}
		// check for node type - condition type has no children
		if n.NodeType == NodeTypeCondition {
			return Node{}, false
		}
		// return next child
		return n.Children[visited], true
	}
}

func (n Node) IsSatisfied(value interface{}) bool {
	o, ok := reg.Get(n.Operator)
	if !ok {
		return false
	}
	return o.IsSatisfied(n, value)
}

func (n Node) isEmpty() bool {
	return n.Operator == ""
}
