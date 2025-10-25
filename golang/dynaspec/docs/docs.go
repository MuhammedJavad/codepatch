package docs

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/MuhammedJavad/codepatch/dynaspec"

	acc "github.com/MuhammedJavad/codepatch/dynaspec/accessor/mysql"
	"github.com/MuhammedJavad/codepatch/dynaspec/tree"
)

type CustomerIdMustBeInCondition struct{}
type QuantityMustBeGreaterThanCondition struct{}
type ProductIdMustBeInCondition struct{}

func (c CustomerIdMustBeInCondition) IsSatisfied(n tree.Node, value interface{}) bool {
	expectedIds := strings.Split(n.Operand, ",")

	for _, op := range expectedIds {
		expected, err := strconv.Atoi(strings.TrimSpace(op))
		if err != nil {
			continue
		}
		if value.(Order).CustomerID == uint(expected) {
			return true
		}
	}

	return false
}

func (c QuantityMustBeGreaterThanCondition) IsSatisfied(n tree.Node, value interface{}) bool {
	expected, err := strconv.Atoi(n.Operand)
	if err != nil {
		return false
	}
	return value.(Order).Quantity > uint(expected)
}

func (c ProductIdMustBeInCondition) IsSatisfied(n tree.Node, value interface{}) bool {
	expectedIds := strings.Split(n.Operand, ",")
	for _, op := range expectedIds {
		expected, err := strconv.Atoi(strings.TrimSpace(op))
		if err != nil {
			continue
		}
		if value.(Order).ProductID == uint(expected) {
			return true
		}
	}
	return false
}

type Order struct {
	CustomerID uint
	ProductID  uint
	Quantity   uint
}

func main() {
	dynaspec.Register(
		dynaspec.Use("customerIdMustBeIn", CustomerIdMustBeInCondition{}),
		dynaspec.Use("quantityMustBeGreaterThan", QuantityMustBeGreaterThanCondition{}),
		dynaspec.Use("productIdMustBeIn", ProductIdMustBeInCondition{}))

	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	accessor := acc.New(db)
	// get the tree from the accessor
	tree, err := accessor.Get(context.Background(), 1)
	if err != nil {
		log.Fatal(err)
	}

	order := Order{
		CustomerID: 1234,
		ProductID:  1234,
		Quantity:   1,
	}
	if tree.Traverse(order) {
		fmt.Printf("Tree is satisfied. [result: %v]\n", tree.Result)
	} else {
		fmt.Println("Tree is not satisfied")
	}
}

func builderExample() {
	_ = tree.Builder("example", nil).
		WithStartTime(time.Now()).
		WithEndTime(time.Now().Add(time.Hour * 24)).
		WithRoot(func(nb *tree.NodeBuilder) {
			nb.AsGate("and", func(gb *tree.GateBuilder) {
				gb.AddCondition("quantityMustBeGreaterThan", "10")
				gb.AddGate("or", func(gb *tree.GateBuilder) {
					gb.AddCondition("productIdMustBeIn", "1234")
					gb.AddCondition("productIdMustBeIn", "1234")
				})
			})
		})
}
