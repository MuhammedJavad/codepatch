package tree

import (
	"encoding/json"
	"time"
)

type Tree struct {
	ID     uint            `json:"id"`
	Name   string          `json:"name"`
	Start  *time.Time      `json:"start"`
	End    *time.Time      `json:"end"`
	Active bool            `json:"active"`
	Result json.RawMessage `json:"result"`
	Root   Node            `json:"root"`
}

func (t Tree) Traverse(value interface{}) bool {
	if !t.Active {
		return false
	}
	if t.Start != nil && t.Start.Before(time.Now()) {
		return false
	}
	if t.End != nil && t.End.After(time.Now()) {
		return false
	}
	return t.Root.IsSatisfied(value)
}

func NewTree(id uint, name string, start, end *time.Time, active bool, result json.RawMessage, root Node) Tree {
	return Tree{
		ID:     id,
		Name:   name,
		Start:  start,
		End:    end,
		Active: active,
		Result: result,
		Root:   root,
	}
}
