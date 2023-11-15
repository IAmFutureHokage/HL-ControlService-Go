package dto

import (
	"github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
)

type GetResponse struct {
	Page    uint32
	MaxPage uint32
	Data    []model.NFAD
}
