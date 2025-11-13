package and

import "github.com/MuhammedJavad/codepatch/dynaspec/tree"

const AndGate = "and"

type And struct{}

func (a And) IsSatisfied(n tree.Node, value interface{}) (bool, error) {
	next := n.Next()
	for {
		node, ok := next()
		if !ok {
			// end of the tree
			// we return a true flag
			// indicating that all nodes are satisfied
			return true, nil
		}
		satisfied, err := node.IsSatisfied(value)
		if err != nil {
			return false, err
		}
		if !satisfied {
			return false, nil
		}
	}
}
