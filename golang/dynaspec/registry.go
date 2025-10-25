package dynaspec

import (
	"github.com/MuhammedJavad/codepatch/dynaspec/gates/and"
	"github.com/MuhammedJavad/codepatch/dynaspec/gates/or"
	"github.com/MuhammedJavad/codepatch/dynaspec/gates/xor"
	"github.com/MuhammedJavad/codepatch/dynaspec/tree"
)

type (
	registry map[string]tree.Operator
	spec struct {
		name     string
		operator tree.Operator
	}
)

func (r registry) Get(name string) (tree.Operator, bool) {
	o, ok := r[name]
	return o, ok
}

func Use(name string, operator tree.Operator) spec {
	return spec{
		name:     name,
		operator: operator,
	}
}

func Register(specs ...spec) {
	r := make(registry)
	r[and.AndGate] = and.And{}
	r[and.NandGate] = and.Nand{}
	r[or.OrGate] = or.Or{}
	r[or.NorGate] = or.Nor{}
	r[xor.XorGate] = xor.Xor{}
	r[xor.XnorGate] = xor.Xnor{}

	for _, spec := range specs {
		r[spec.name] = spec.operator
	}
	
	tree.UseRegistry(r)
}
