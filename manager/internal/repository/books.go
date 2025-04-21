package repository

import (
	"context"
	"errors"

	"github.com/csokviktor/lib_manager/internal/repository/dbmodels"
	"gorm.io/gorm"
)

type books struct {
	db *gorm.DB
}

//go:generate mockgen -source books.go -destination mocks/books.go --package mocks

type Books interface {
	CRUD[dbmodels.Book]
}

func NewBooks(db *gorm.DB) Books {
	return &books{
		db,
	}
}

func (b *books) Create(ctx context.Context, book *dbmodels.Book) (*dbmodels.Book, error) {
	result := b.db.WithContext(ctx).Create(book)
	return book, result.Error
}

func (b *books) GetAll(ctx context.Context) ([]dbmodels.Book, error) {
	var books []dbmodels.Book
	result := b.db.WithContext(ctx).Find(&books)
	return books, result.Error
}

func (b *books) GetByID(ctx context.Context, id string) (*dbmodels.Book, error) {
	var book dbmodels.Book
	result := b.db.WithContext(ctx).First(&book, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFoundByID(id)
	}
	return &book, result.Error
}

func (b *books) Update(ctx context.Context, book *dbmodels.Book) (*dbmodels.Book, error) {
	if _, err := b.GetByID(ctx, book.ID); err != nil {
		return nil, err
	}
	result := b.db.WithContext(ctx).Updates(book)
	return book, result.Error
}

func (b *books) Delete(ctx context.Context, id string) error {
	return b.db.WithContext(ctx).Delete(&dbmodels.Book{}, id).Error
}
