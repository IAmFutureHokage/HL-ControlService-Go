package model

import (
	"time"
)

type NFAD struct {
	ID        string      `gorm:"primarykey"`
	PostCode  uint32      `gorm:"not null"`
	Type      ControlType `gorm:"not null"`
	DateStart time.Time   `gorm:"not null"`
	DateEnd   time.Time
	Value     uint32 `gorm:"not null"`
}
