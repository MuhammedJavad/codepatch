package registrar

import (
	"github.com/MuhammedJavad/codepatch/dynaspec/gates/and"
	"github.com/MuhammedJavad/codepatch/dynaspec/gates/or"
	"github.com/MuhammedJavad/codepatch/dynaspec/gates/xor"
	"github.com/MuhammedJavad/codepatch/dynaspec/tree"
)

type (
	registry struct {
		conditions map[string]tree.Operator
		gates      map[string]tree.Operator
	}
	Spec     struct {
		name     string
		operator tree.Operator
	}
)

func (r *registry) Get(name string) (tree.Operator, bool) {
	o, ok := r.gates[name]
	if !ok {
		o, ok = r.conditions[name]
	}
	return o, ok
}

func (r *registry) IsGate(name string) bool {
	_, ok := r.gates[name]
	return ok
}

func (r *registry) IsCondition(name string) bool {
	_, ok := r.conditions[name]
	return ok
}

func (r *registry) ConditionOperators() []string {
	keys := make([]string, 0, len(r.conditions))
	for k := range r.conditions {
		keys = append(keys, k)
	}
	return keys
}

func (r *registry) GateOperators() []string {
	keys := make([]string, 0, len(r.gates))
	for k := range r.gates {
		keys = append(keys, k)
	}
	return keys
}

var reg = registry{
	gates: map[string]tree.Operator{
		and.AndGate:  and.And{},
		and.NandGate: and.Nand{},
		or.OrGate:    or.Or{},
		or.NorGate:   or.Nor{},
		xor.XorGate:  xor.Xor{},
		xor.XnorGate: xor.Xnor{},
	},
	conditions: make(map[string]tree.Operator, 0),
}

func Use(name string, operator tree.Operator) Spec {
	return Spec{
		name:     name,
		operator: operator,
	}
}

func Register(specs ...Spec) {
	reg.conditions = make(map[string]tree.Operator, len(specs))
	
	for _, spec := range specs {
		reg.conditions[spec.name] = spec.operator
	}

	tree.UseRegistry(&reg)
}
