package or

import "github.com/MuhammedJavad/codepatch/dynaspec/tree"

const OrGate = "or"

type Or struct{}

func (o Or) IsSatisfied(n tree.Node, value interface{}) bool {
	next := n.Next()
	for {
		node, ok := next()
		if !ok {
			return false
		}
		if node.IsSatisfied(value) {
			return true
		}
	}
}