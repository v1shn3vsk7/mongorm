package action

type actionType uint8

const (
	_ = iota
	Where
	And
	Or
)
