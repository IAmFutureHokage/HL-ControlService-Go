package model

import (
	"time"

	"github.com/google/uuid"
)

type Waterlevel struct {
	Id         uuid.UUID
	PostCode   string
	Date       time.Time
	Waterlevel uint32
}

type KafkaMessage struct {
	PostCode   string    `json:"PostCode"`
	Date       time.Time `json:"Date"`
	WaterLevel uint32    `json:"WaterLevel"`
}
