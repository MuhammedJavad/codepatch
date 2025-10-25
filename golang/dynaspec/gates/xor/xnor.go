package xor

import "github.com/MuhammedJavad/codepatch/dynaspec/tree"

const XnorGate = "xnor"

type Xnor struct{}

func (x Xnor) IsSatisfied(n tree.Node, value interface{}) bool {
	next := n.Next()
	trueCount := 0
	for {
		node, ok := next()
		if !ok {
			break
		}
		if node.IsSatisfied(value) {
			trueCount++
		}
	}
	// XNOR over multiple operands: true if an even number of children are true
	return trueCount%2 == 0
}