package model

type ControlValueType byte

const (
	None ControlValueType = iota
	Norm
	Floodplain
	Adverse
	Dangerous
)

func (ct ControlValueType) ToByte() byte {
	return byte(ct)
}

func FromByte(b byte) ControlValueType {
	return ControlValueType(b)
}
