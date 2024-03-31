package mongorm

type CondOperator uint8

const (
	_ CondOperator = iota
	EQ
	NE
	GT
	GTE
	LT
	LTE
	IN
)

func (o CondOperator) Uint8() uint8 {
	return uint8(o)
}
