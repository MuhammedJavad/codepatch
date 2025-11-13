package and

import "github.com/MuhammedJavad/codepatch/dynaspec/tree"

const NandGate = "nand"

// Nand implements logical AND over negated children
type Nand struct{}

func (a Nand) IsSatisfied(n tree.Node, value interface{}) (bool, error) {
	next := n.Next()
	for {
		node, ok := next()
		if !ok {
			return true, nil
		}
		satisfied, err := node.IsSatisfied(value)
		if err != nil {
			return false, err
		}
		if satisfied {
			return false, nil
		}
	}
}
