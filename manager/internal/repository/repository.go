package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type GenericRepository[T any] struct {
	db *gorm.DB
}

func NewGenericRepository[T any](db *gorm.DB) CRUD[T] {
	return &GenericRepository[T]{db: db}
}

func (r *GenericRepository[T]) Create(ctx context.Context, entity *T) (*T, error) {
	result := r.db.WithContext(ctx).Create(entity)
	return entity, result.Error
}

func (r *GenericRepository[T]) GetAll(ctx context.Context) ([]T, error) {
	var entities []T
	result := r.db.WithContext(ctx).Find(&entities)
	return entities, result.Error
}

func (r *GenericRepository[T]) GetByID(ctx context.Context, id string) (*T, error) {
	var entity T
	result := r.db.WithContext(ctx).First(&entity, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFoundByID(id)
	}
	return &entity, result.Error
}

func (r *GenericRepository[T]) Update(ctx context.Context, entity *T) (*T, error) {
	result := r.db.WithContext(ctx).Updates(entity)
	return entity, result.Error
}

func (r *GenericRepository[T]) Delete(ctx context.Context, id string) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}
