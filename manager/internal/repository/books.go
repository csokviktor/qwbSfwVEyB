package repository

import (
	"github.com/csokviktor/qwbSfwVEyB/manager/internal/repository/dbmodels"
	"gorm.io/gorm"
)

type books struct {
	db *gorm.DB
}

//go:generate mockgen -source books.go -destination mocks/books.go --package mocks

type Books interface {
	CRUD[dbmodels.Book]
}

// NewBooksRepository creates a new book repository
func NewBooks(db *gorm.DB) Books {
	return NewGenericRepository[dbmodels.Book](db)
}
