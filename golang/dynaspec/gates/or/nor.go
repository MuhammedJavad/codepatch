package or

import "github.com/MuhammedJavad/codepatch/dynaspec/tree"

const NorGate = "nor"

// Nor implements logical NOR: true only if all children are not satisfied
type Nor struct{}

func (n Nor) IsSatisfied(node tree.Node, value interface{}) (bool, error) {
	next := node.Next()
	for {
		child, ok := next()
		if !ok {
			return true, nil
		}
		satisfied, err := child.IsSatisfied(value)
		if err != nil {
			return false, err
		}
		if !satisfied {
			return true, nil
		}
	}
}
