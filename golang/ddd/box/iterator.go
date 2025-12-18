package box

import (
	"context"
	"log/slog"
)

type (
	Iterator interface {
		Next() bool
		Current() Message
		Confirm(ctx context.Context, err error)
	}

	iteratorImpl struct {
		accessor Accessor
		messages []Message
		index    int
	}
)

func newIterator(accessor Accessor, messages []Message) Iterator {
	return &iteratorImpl{accessor: accessor, messages: messages, index: -1}
}

func (i *iteratorImpl) Next() bool {
	if i.index < len(i.messages)-1 {
		i.index++
		return true
	}
	return false
}

func (i *iteratorImpl) Current() Message {
	return i.messages[i.index]
}

func (i *iteratorImpl) Confirm(ctx context.Context, err error) {
	i.messages[i.index].confirm(err)
	if err := i.accessor.Update(ctx, i.messages[i.index]); err != nil {
		slog.ErrorContext(ctx, "failed to confirm message", "error", err)
	}
}
