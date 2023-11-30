package model

import (
	"time"
)

type NFAD struct {
	ID        string           `gorm:"primarykey;not null"`
	PostCode  uint32           `gorm:"not null"`
	Type      ControlValueType `gorm:"not null"`
	DateStart time.Time        `gorm:"not null"`
	Value     uint32           `gorm:"not null"`
}
