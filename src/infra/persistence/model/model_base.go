package model

import (
	"time"

	"gorm.io/gorm"
)

// Base .
type Base struct {
	ID        uint           `gorm:"column:id" json:"-"`
	CreatedAt time.Time      `gorm:"index:,sort:desc" json:"-"`
	UpdatedAt time.Time      `gorm:"index,sort:desc" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
