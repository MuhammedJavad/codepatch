package tree

type (
	Operator interface {
		IsSatisfied(Node, interface{}) (bool, error)
	}
	registrar interface {
		Get(name string) (Operator, bool)
		IsGate(name string) bool
		IsCondition(name string) bool
		ConditionOperators() []string
		GateOperators() []string
	}
)

var reg registrar

func UseRegistry(r registrar) {
	reg = r
}
