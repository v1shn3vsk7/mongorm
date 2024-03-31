package action

type Action struct {
	Type     actionType
	Key      string
	Operator uint8
	Value    interface{}
}
