package keyvaluestore

type Command struct {
	Operation string
	Operand1  interface{}
	Operand2  interface{}
}
