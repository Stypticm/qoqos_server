package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockMovement struct {
	ID         string    `gorm:"primaryKey;type:text;column:id" json:"id"`
	PartName   string    `gorm:"column:partName" json:"partName"`
	PartModel  string    `gorm:"column:partModel" json:"partModel"`
	Delta      int       `gorm:"column:delta" json:"delta"`
	Reason     string    `gorm:"column:reason" json:"reason"`
	ReferenceId *string  `gorm:"column:referenceId" json:"referenceId"`
	CreatedAt  time.Time `gorm:"column:createdAt" json:"createdAt"`
}

func (StockMovement) TableName() string {
	return "StockMovement"
}

func (m *StockMovement) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return nil
}
