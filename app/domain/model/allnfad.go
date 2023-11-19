package model

import "time"

type AllNFAD struct {
	Date       time.Time
	Norm       uint32
	Floodplain uint32
	Adverse    uint32
	Dangerous  uint32
}
