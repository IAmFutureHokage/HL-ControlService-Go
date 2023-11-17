package dto

import (
	"time"

	"github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
)

type CreateRequest struct {
	PostCode  uint32
	Type      model.ControlType
	DateStart time.Time
	Value     uint32
}
