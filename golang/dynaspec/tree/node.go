package tree

import (
	"fmt"
	"strconv"
	"strings"
)

type NodeType string

const (
	NodeTypeGate      = "gate"
	NodeTypeCondition = "condition"
)

type Node struct {
	Operand  string   `json:"operand,omitempty"`
	NodeType NodeType `json:"type,omitempty"`
	Operator string   `json:"operator,omitempty"`
	Children []Node   `json:"children,omitempty"`
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

func (n Node) IsSatisfied(value interface{}) (bool, error) {
	o, ok := reg.Get(n.Operator)
	if !ok {
		return false, fmt.Errorf("operator %s not found", n.Operator)
	}
	return o.IsSatisfied(n, value)
}

func (n Node) IsEmpty() bool {
	return n.Operator == ""
}

// cast converts a string to a specified generic type T.
// Supported types: string, int, int64, uint32, uint64, float64, bool.
func Operand[T any](n Node) (T, error) {
	var zero T
	if n.Operand == "" {
		return zero, fmt.Errorf("operand for node [operator:%s, nodeType:%s] is empty", n.Operator, n.NodeType)
	}

	switch any(zero).(type) {
	case string:
		return any(n.Operand).(T), nil

	case int:
		v, err := strconv.Atoi(n.Operand)
		if err != nil {
			return zero, fmt.Errorf("operand for node [operator:%s, nodeType:%s] must be int [actual: %s]", n.Operator, n.NodeType, n.Operand)
		}
		return any(v).(T), nil

	case int64:
		v, err := strconv.ParseInt(n.Operand, 10, 64)
		if err != nil {
			return zero, fmt.Errorf("operand for node [operator:%s, nodeType:%s] must be int64 [actual: %s]", n.Operator, n.NodeType, n.Operand)
		}
		return any(v).(T), nil

	case uint32:
		v, err := strconv.ParseUint(n.Operand, 10, 32)
		if err != nil {
			return zero, fmt.Errorf("operand for node [operator:%s, nodeType:%s] must be uint32 [actual: %s]", n.Operator, n.NodeType, n.Operand)
		}
		return any(uint32(v)).(T), nil

	case uint64:
		v, err := strconv.ParseUint(n.Operand, 10, 64)
		if err != nil {
			return zero, fmt.Errorf("operand for node [operator:%s, nodeType:%s] must be uint64 [actual: %s]", n.Operator, n.NodeType, n.Operand)
		}
		return any(v).(T), nil

	case float64:
		v, err := strconv.ParseFloat(n.Operand, 64)
		if err != nil {
			return zero, fmt.Errorf("operand for node [operator:%s, nodeType:%s] must be float64 [actual: %s]", n.Operator, n.NodeType, n.Operand)
		}
		return any(v).(T), nil

	case bool:
		v, err := strconv.ParseBool(n.Operand)
		if err != nil {
			return zero, fmt.Errorf("operand for node [operator:%s, nodeType:%s] must be bool [actual: %s]", n.Operator, n.NodeType, n.Operand)
		}
		return any(v).(T), nil

	case []int:
		oprInt := make([]int, 0)
		oprStr := strings.Split(n.Operand, ",")
		for _, operand := range oprStr {
			v, err := strconv.ParseInt(operand, 10, 64)
			if err != nil {
				return zero, fmt.Errorf("operand for node [operator:%s, nodeType:%s] must be []int [actual: %s]", n.Operator, n.NodeType, n.Operand)
			}
			oprInt = append(oprInt, int(v))
		}
		return any(oprInt).(T), nil
		
	case []string:
		oprStr := strings.Split(n.Operand, ",")
		return any(oprStr).(T), nil

	default:
		return zero, fmt.Errorf("operand for node [operator:%s, nodeType:%s] must be %T [actual: %s]", n.Operator, n.NodeType, zero, n.Operand)
	}
}
