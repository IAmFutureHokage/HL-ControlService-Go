package model

import (
	"time"

	"github.com/google/uuid"
)

type ControlValue struct {
	ID        uuid.UUID
	PostCode  string
	Type      ControlValueType
	DateStart time.Time
	Value     uint32
}
