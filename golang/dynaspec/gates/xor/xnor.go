package xor

import "github.com/MuhammedJavad/codepatch/dynaspec/tree"

const XnorGate = "xnor"

type Xnor struct{}

func (x Xnor) IsSatisfied(n tree.Node, value interface{}) (bool, error) {
	next := n.Next()
	trueCount := 0
	for {
		node, ok := next()
		if !ok {
			break
		}
		satisfied, err := node.IsSatisfied(value)
		if err != nil {
			return false, err
		}
		if satisfied {
			trueCount++
		}
	}
	// XNOR over multiple operands: true if an even number of children are true
	return trueCount%2 == 0, nil
}