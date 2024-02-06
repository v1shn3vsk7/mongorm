package mongorm

// high-level operators
const (
	EQ  CondOperator = 0
	NE  CondOperator = 1
	GT  CondOperator = 2
	GTE CondOperator = 3
	LT  CondOperator = 4
	LTE CondOperator = 5
	In  CondOperator = 6
)

type CondOperator int

func (o CondOperator) Int() int {
	return int(o)
}
