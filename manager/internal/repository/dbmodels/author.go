package dbmodels

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Author struct {
	ID    string `gorm:"primaryKey;type:uuid;unique"`
	Name  string `gorm:"not null"`
	Books []Book `gorm:"foreignKey:author_id"`
}

func (Author) TableName() string {
	return "authors"
}

func (a *Author) BeforeCreate(_ *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}
