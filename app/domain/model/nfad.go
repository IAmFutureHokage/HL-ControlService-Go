package model

import (
	"time"
)

type NFAD struct {
	ID        string `gorm:"primarykey"`
	PostCode  int32
	Type      ControlType
	DateStart time.Time
	DateEnd   time.Time
	Value     int32
}
