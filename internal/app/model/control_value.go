package model

import (
	"time"
)

type ControlValue struct {
	ID        string
	PostCode  string
	Type      ControlValueType
	DateStart time.Time
	Value     uint32
}
