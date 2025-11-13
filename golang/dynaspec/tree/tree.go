package tree

import 	"time"


type Tree struct {
	ID      uint       `json:"id"`
	Start   *time.Time `json:"start"`
	End     *time.Time `json:"end"`
	Active  bool       `json:"active"`
	Result  float32    `json:"result"`
	Root    Node       `json:"root"`
}

func (t Tree) Traverse(value interface{}) (bool, error) {
	if !t.Active {
		return false, nil
	}
	if t.Start != nil && time.Now().Before(*t.Start) {
		return false, nil
	}
	if t.End != nil && time.Now().After(*t.End) {
		return false, nil
	}
	return t.Root.IsSatisfied(value)
}