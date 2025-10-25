package dynaspec

import (
	"github.com/MuhammedJavad/codepatch/dynaspec/tree"
	"strconv"
	"testing"
)

// ----------- Model -----------
type Value struct {
	Flag   bool
	Number int
	Text   string
}

// ----------- Operators -----------
type flagMustBeAsExpectedCondition struct{}
type numberMustBeAsExpectedCondition struct{}
type textMustBeAsExpectedCondition struct{}

func (c flagMustBeAsExpectedCondition) IsSatisfied(t tree.Node, value interface{}) bool {
	expected, err := strconv.ParseBool(t.Operand)
	if err != nil {
		return false
	}
	return value.(Value).Flag == expected
}

func (c numberMustBeAsExpectedCondition) IsSatisfied(t tree.Node, value interface{}) bool {
	expected, err := strconv.Atoi(t.Operand)
	if err != nil {
		return false
	}
	return value.(Value).Number == expected
}

func (c textMustBeAsExpectedCondition) IsSatisfied(t tree.Node, value interface{}) bool {
	return value.(Value).Text == t.Operand
}

func registerOperators() {
	Register(
		Use("flagMustBeAsExpected", flagMustBeAsExpectedCondition{}),
		Use("numberMustBeAsExpected", numberMustBeAsExpectedCondition{}),
		Use("textMustBeAsExpected", textMustBeAsExpectedCondition{}),
	)
}

// createTree creates a test node with specified depth
func createTree(depth int) tree.Tree {
	if depth == 0 {
		return tree.Tree{
			ID:     1,
			Name:   "test_tree",
			Active: true,
			Root: tree.Node{
				Operand:  "true",
				NodeType: tree.NodeTypeCondition,
				Operator: "flagMustBeAsExpected",
				Children: []tree.Node{},
			},
		}
	}

	// TODO; use Builder here

	return tree.Tree{}
}

// BenchmarkNodeByValue tests performance when Node is passed by value
func BenchmarkNodeByValue(b *testing.B) {
	// // Register operator
	// Register(
	// 	Use("value_ref_operator", ValueRefOperator{result: true}),
	// )

	// // Create a tree with moderate depth
	// root := createTree(4) // 3^4 = 81 leaf nodes
	// tree := tree.Tree{
	// 	ID:     1,
	// 	Name:   "value_test_tree",
	// 	Active: true,
	// 	Root:   root,
	// }

	// b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	// 	tree.Traverse("test_value")
	// }
}

// BenchmarkNodeByValueDeep tests performance with deeper tree
func BenchmarkNodeByValueDeep(b *testing.B) {
	// // Register operator
	// Register(
	// 	Use("value_ref_operator", ValueRefOperator{result: true}),
	// )

	// // Create a deeper tree
	// root := createTree(5) // 3^5 = 243 leaf nodes
	// tree := tree.Tree{
	// 	ID:     1,
	// 	Name:   "value_deep_test_tree",
	// 	Active: true,
	// 	Root:   root,
	// }

	// b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	// 	tree.Traverse("test_value")
	// }
}

// BenchmarkNodeByValueWide tests performance with wide tree
func BenchmarkNodeByValueWide(b *testing.B) {
	// // Register operator
	// Register(
	// 	Use("value_ref_operator", ValueRefOperator{result: true}),
	// )

	// // Create a wide tree (many children at root level)
	// children := make([]tree.Node, 100)
	// for i := 0; i < 100; i++ {
	// 	children[i] = tree.Node{
	// 		Operand:  "condition",
	// 		NodeType: tree.NodeTypeCondition,
	// 		Operator: "value_ref_operator",
	// 		Children: []tree.Node{},
	// 	}
	// }

	// root := tree.Node{
	// 	Operand:  "",
	// 	NodeType: tree.NodeTypeGate,
	// 	Operator: "and",
	// 	Children: children,
	// }

	// tree := tree.Tree{
	// 	ID:     1,
	// 	Name:   "value_wide_test_tree",
	// 	Active: true,
	// 	Root:   root,
	// }

	// b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	// 	tree.Traverse("test_value")
	// }
}

// BenchmarkNodeByValueWithFalseConditions tests performance with early termination
func BenchmarkNodeByValueWithFalseConditions(b *testing.B) {
	// // Register operators
	// Register(
	// 	Use("value_ref_operator", ValueRefOperator{result: true}),
	// 	Use("value_ref_false", ValueRefOperator{result: false}),
	// )

	// // Create a tree where first condition is false (tests short-circuit)
	// root := tree.Node{
	// 	Operand:  "",
	// 	NodeType: tree.NodeTypeGate,
	// 	Operator: "and",
	// 	Children: []tree.Node{
	// 		{
	// 			Operand:  "false_condition",
	// 			NodeType: tree.NodeTypeCondition,
	// 			Operator: "value_ref_false",
	// 			Children: []tree.Node{},
	// 		},
	// 		// These won't be evaluated due to short-circuit
	// 		createTree(4),
	// 		createTree(4),
	// 	},
	// }

	// tree := tree.Tree{
	// 	ID:     1,
	// 	Name:   "value_short_circuit_tree",
	// 	Active: true,
	// 	Root:   root,
	// }

	// b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	// 	tree.Traverse("test_value")
	// }
}

// BenchmarkNodeByValueComplexStructure tests performance with complex nested structure
func BenchmarkNodeByValueComplexStructure(b *testing.B) {
	// // Register operator
	// Register(
	// 	Use("value_ref_operator", ValueRefOperator{result: true}),
	// )

	// // Create a complex tree with mixed AND/OR gates
	// root := tree.Node{
	// 	Operand:  "",
	// 	NodeType: tree.NodeTypeGate,
	// 	Operator: "and",
	// 	Children: []tree.Node{
	// 		{
	// 			Operand:  "",
	// 			NodeType: tree.NodeTypeGate,
	// 			Operator: "or",
	// 			Children: []tree.Node{
	// 				createTree(2),
	// 				createTree(2),
	// 				createTree(2),
	// 			},
	// 		},
	// 		{
	// 			Operand:  "",
	// 			NodeType: tree.NodeTypeGate,
	// 			Operator: "and",
	// 			Children: []tree.Node{
	// 				createTree(2),
	// 				createTree(2),
	// 			},
	// 		},
	// 	},
	// }

	// tree := tree.Tree{
	// 	ID:     1,
	// 	Name:   "value_complex_tree",
	// 	Active: true,
	// 	Root:   root,
	// }

	// b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	// 	tree.Traverse("test_value")
	// }
}

// BenchmarkNodeMemoryAllocation tests memory allocation patterns
func BenchmarkNodeMemoryAllocation(b *testing.B) {
	// // Register operator
	// Register(
	// 	Use("value_ref_operator", ValueRefOperator{result: true}),
	// )

	// // Create a tree
	// root := createTree(4)
	// tree := tree.Tree{
	// 	ID:     1,
	// 	Name:   "memory_test_tree",
	// 	Active: true,
	// 	Root:   root,
	// }

	// b.ResetTimer()
	// b.ReportAllocs()
	// for i := 0; i < b.N; i++ {
	// 	tree.Traverse("test_value")
	// }
}

// BenchmarkNodeCopyOverhead tests the overhead of copying Node structs
func BenchmarkNodeCopyOverhead(b *testing.B) {
	// // Register operator
	// Register(
	// 	Use("value_ref_operator", ValueRefOperator{result: true}),
	// )

	// // Create a large node to test copy overhead
	// largeNode := createTree(4)

	// b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	// 	// Simulate copying the node (this happens in the current implementation)
	// 	copiedNode := largeNode
	// 	_ = copiedNode
	// }
}

// BenchmarkNodeSliceOperations tests performance of slice operations on children
func BenchmarkNodeSliceOperations(b *testing.B) {
	// // Register operator
	// Register(
	// 	Use("value_ref_operator", ValueRefOperator{result: true}),
	// )

	// // Create a node with many children
	// children := make([]tree.Node, 1000)
	// for i := 0; i < 1000; i++ {
	// 	children[i] = tree.Node{
	// 		Operand:  "condition",
	// 		NodeType: tree.NodeTypeCondition,
	// 		Operator: "value_ref_operator",
	// 		Children: []tree.Node{},
	// 	}
	// }

	// root := tree.Node{
	// 	Operand:  "",
	// 	NodeType: tree.NodeTypeGate,
	// 	Operator: "and",
	// 	Children: children,
	// }

	// tree := tree.Tree{
	// 	ID:     1,
	// 	Name:   "slice_ops_tree",
	// 	Active: true,
	// 	Root:   root,
	// }

	// b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	// 	tree.Traverse("test_value")
	// }
}
