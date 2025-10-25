package tree

type (
	Operator interface {
		IsSatisfied(Node, interface{}) bool
	}

	registry interface {
		Get(name string) (Operator, bool)
	}
)

var reg registry

func UseRegistry(r registry) {
	reg = r
}
