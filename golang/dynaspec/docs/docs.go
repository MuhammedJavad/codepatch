package docs

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	dynar "github.com/MuhammedJavad/codepatch/dynaspec/registrar"
	dyna "github.com/MuhammedJavad/codepatch/dynaspec/tree"
)

type CustomerIdMustBeInCondition struct{}
type QuantityMustBeGreaterThanCondition struct{}
type ProductIdMustBeInCondition struct{}

func (c CustomerIdMustBeInCondition) IsSatisfied(n tree.Node, value interface{}) (bool, error) {
	expectedIds := strings.Split(n.Operand, ",")

	for _, op := range expectedIds {
		expected, err := strconv.Atoi(strings.TrimSpace(op))
		if err != nil {
			continue
		}
		if value.(Order).CustomerID == uint(expected) {
			return true, nil
		}
	}

	return false, nil
}

func (c QuantityMustBeGreaterThanCondition) IsSatisfied(n tree.Node, value interface{}) (bool, error) {
	expected, err := strconv.Atoi(n.Operand)
	if err != nil {
		return false, err
	}
	return value.(Order).Quantity > uint(expected), nil
}

func (c ProductIdMustBeInCondition) IsSatisfied(n tree.Node, value interface{}) (bool, error) {
	expectedIds := strings.Split(n.Operand, ",")
	for _, op := range expectedIds {
		expected, err := strconv.Atoi(strings.TrimSpace(op))
		if err != nil {
			continue
		}
		if value.(Order).ProductID == uint(expected) {
			return true, nil
		}
	}
	return false, nil
}

type Order struct {
	CustomerID uint
	ProductID  uint
	Quantity   uint
}

func main() {
	dynar.Register(
		dynar.Use("customerIdMustBeIn", CustomerIdMustBeInCondition{}),
		dynar.Use("quantityMustBeGreaterThan", QuantityMustBeGreaterThanCondition{}),
		dynar.Use("productIdMustBeIn", ProductIdMustBeInCondition{}))

	now := time.Now()
	end := now.Add(time.Hour * 24)

	tree, err := dyna.Builder(1.0).
		WithStartTime(&now).
		WithEndTime(&end).
		WithRoot(func(nb *dyna.NodeBuilder) {
			nb.AsGate("and", func(gb *dyna.GateBuilder) {
				gb.AddCondition("quantityMustBeGreaterThan", "10")
			})
		}).
		Build()
	if err != nil {
		fmt.Printf("Error building tree: %v\n", err)
		return
	}

	order := Order{
		CustomerID: 1234,
		ProductID:  1234,
		Quantity:   1,
	}
	if satisfied, err := tree.Traverse(order); err != nil {
		fmt.Printf("Error traversing tree: %v\n", err)
	} else if satisfied {
		fmt.Printf("Tree is satisfied. [result: %v]\n", tree.Result)
	}

}

func builderExample() {
	now := time.Now()
	end := now.Add(time.Hour * 24)

	_, _ = dyna.Builder(1.0).
		WithStartTime(&now).
		WithEndTime(&end).
		WithRoot(func(nb *dyna.NodeBuilder) {
			nb.AsGate("and", func(gb *dyna.GateBuilder) {
				gb.AddCondition("quantityMustBeGreaterThan", "10")
				gb.AddGate("or", func(gb *dyna.GateBuilder) {
					gb.AddCondition("productIdMustBeIn", "1234")
					gb.AddCondition("productIdMustBeIn", "1234")
				})
			})
		}).
		Build()
}
