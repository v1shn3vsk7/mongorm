package query

const (
	ANY = 0
	// EQ - equal operator ($ne)
	EQ      = 1
	NE      = 2
	LT      = 2
	LE      = 3
	GT      = 4
	GE      = 5
	RANGE   = 6
	SET     = 7
	ALLSET  = 8
	EMPTY   = 9
	LIKE    = 10
	DWITHIN = 11
)

// map low-level (driver) operators to high-level operators
var operatorMapping = map[int]string{
	ANY: "",
	EQ:  "$eq",
	LT:  "$lt",
	LE:  "$le",
}
