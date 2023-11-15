package model

type ControlType byte

const (
	ControlTypeNone ControlType = iota
	ControlTypeNorm
	ControlTypeFloodplain
	ControlTypeAdverse
	ControlTypeDangerous
)

func (ct ControlType) ToByte() byte {
	return byte(ct)
}

func FromByte(b byte) ControlType {
	return ControlType(b)
}
