package dbmodels

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Borrower struct {
	ID    string `gorm:"primaryKey;type:uuid;unique"`
	Name  string `gorm:"not null"`
	Books []Book `gorm:"foreignKey:borrower_id"`
}

func (Borrower) TableName() string {
	return "borrowers"
}

func (b *Borrower) BeforeCreate(_ *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}
