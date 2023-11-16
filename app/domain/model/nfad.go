package model

import (
	"time"
)

type NFAD struct {
	ID        string `gorm:"primarykey"`
	PostCode  uint32
	Type      ControlType
	DateStart time.Time
	DateEnd   time.Time
	Value     uint32
}
