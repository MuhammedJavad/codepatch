package xor

import "github.com/MuhammedJavad/codepatch/dynaspec/tree"

const XorGate = "xor"

type Xor struct{}

func (x Xor) IsSatisfied(n tree.Node, value interface{}) bool {
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
	// XOR over multiple operands: true if an odd number of children are true
	return trueCount%2 == 1
}