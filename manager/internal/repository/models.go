package repository

import "context"

type CRUD[T any] interface {
	Create(context.Context, *T) (*T, error)
	GetAll(context.Context) ([]T, error)
	GetByID(context.Context, string) (*T, error)
	Update(context.Context, *T) (*T, error)
	Delete(context.Context, string) error
}

const bookClause = "Books"
