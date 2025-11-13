package or

import "github.com/MuhammedJavad/codepatch/dynaspec/tree"

const OrGate = "or"

type Or struct{}

func (o Or) IsSatisfied(n tree.Node, value interface{}) (bool, error) {
	next := n.Next()
	for {
		node, ok := next()
		if !ok {
			return false, nil
		}
		satisfied, err := node.IsSatisfied(value)
		if err != nil {
			return false, err
		}
		if satisfied {
			return true, nil
		}
	}
}