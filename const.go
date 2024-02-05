package mongorm

// high-level operators
const (
	EQ CondOperator = 0
	NE CondOperator = 1
)

type CondOperator int

func (o CondOperator) Int() int {
	return int(o)
}
