package box

import (
	"context"
)

type (
	Box interface {
		Save(ctx context.Context, message Message) error
		Iterator(ctx context.Context, source, event string, limit int) (Iterator, error)
	}

	// Accessor defines the interface for storing and retrieving messages
	Accessor interface {
		Insert(ctx context.Context, message Message) error
		Update(ctx context.Context, message Message) error
		Retrieve(ctx context.Context, source, event string, limit int) ([]Message, error)
	}

	accessor = Accessor // internal alias for backward compatibility

	impl struct {
		accessor
	}
)

func New(accessor Accessor) *impl {
	return &impl{accessor: accessor}
}

func (i *impl) Save(ctx context.Context, message Message) error {
	return i.accessor.Insert(ctx, message)
}

func (i *impl) Iterator(ctx context.Context, source, event string, limit int) (Iterator, error) {
	messages, err := i.accessor.Retrieve(ctx, source, event, limit)
	if err != nil {
		return nil, err
	}
	return newIterator(i.accessor, messages), nil
}
