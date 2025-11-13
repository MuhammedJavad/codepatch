package xor

import "github.com/MuhammedJavad/codepatch/dynaspec/tree"

const XorGate = "xor"

type Xor struct{}

func (x Xor) IsSatisfied(n tree.Node, value interface{}) (bool, error) {
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
	// XOR over multiple operands: true if an odd number of children are true
	return trueCount%2 == 1, nil
}