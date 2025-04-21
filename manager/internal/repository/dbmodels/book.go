package dbmodels

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	ID         string `gorm:"primaryKey"`
	Title      string
	AuthorID   string    `gorm:"column:author_id"`
	Author     Author    `gorm:"constraint:OnDelete:CASCADE;"`
	BorrowerID *string   `gorm:"column:borrower_id"`
	Borrower   *Borrower `gorm:"constraint:OnDelete:SET NULL;"`
}

func (Book) TableName() string {
	return "books"
}

func (b *Book) BeforeCreate(_ *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}
