package unitofwork

import (
	"fmt"
	"gorm.io/gorm"
)

type (
	Option         func(opt *AdapterOptions) error
	AdapterOptions struct {
		repos map[string]Factory
	}
	Factory  func(db *gorm.DB) interface{}
	Adapters struct {
		tx    *gorm.DB
		repos map[string]Factory
	}
	Interface interface {
		Transact(f func(adapter Adapters) error) error
	}
	uow struct {
		db             *gorm.DB
		adapterOptions AdapterOptions
	}
)

func (u uow) Transact(f func(adapter Adapters) error) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		adapters := Adapters{
			tx:    tx,
			repos: u.adapterOptions.repos,
		}
		return f(adapters)
	})
}

func (a *Adapters) get(key string) (interface{}, bool) {
	factory, ok := a.repos[key]
	if ok {
		repo := factory(a.tx)
		return repo, true
	}
	return nil, false
}

func Get[TRepository interface{}](adapters Adapters, key string) TRepository {
	if obj, ok := adapters.get(key); ok {
		if repo, casted := obj.(TRepository); casted {
			return repo
		}
	}

	panic(fmt.Sprintf("cannot find %s from adapters", key))
}

func Register(key string, factory Factory) Option {
	return func(opt *AdapterOptions) error {
		opt.repos[key] = factory
		return nil
	}
}

func NewUnitOfWork(db *gorm.DB, opts ...Option) *uow {
	options := AdapterOptions{
		repos: make(map[string]Factory),
	}
	for _, opt := range opts {
		opt(&options)
	}
	return &uow{
		db:             db,
		adapterOptions: options,
	}
}
