package model

import (
	"time"
)

type ControlValue struct {
	ID        string           `gorm:"primarykey;not null"`
	PostCode  string           `gorm:"not null"`
	Type      ControlValueType `gorm:"not null"`
	DateStart time.Time        `gorm:"not null"`
	Value     uint32           `gorm:"not null"`
}
