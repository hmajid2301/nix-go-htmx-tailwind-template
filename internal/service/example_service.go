package service

import (
	"context"
)

type ExampleService struct {
	store Storer
}

type Storer interface {
	AddExample(ctx context.Context, field string) error
}

func NewLobbyService(store Storer) *ExampleService {
	return &ExampleService{store: store}
}

func (r *ExampleService) Add(ctx context.Context, field string) error {
	err := r.store.AddExample(ctx, field)
	if err != nil {
		return err
	}

	return nil
}
