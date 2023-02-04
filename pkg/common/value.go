package common

import "fmt"

type DataType int

const (
	NUMBER_DT DataType = iota
	STRING_DT
	BOOLEAN_DT
	FUNCTION_DT
	NULL_DT // REDUNDANT: is it still neccessary?
)

type Value struct {
	DataType DataType
	Data     interface{}
}

func (v Value) Stringify() string {
	switch value := v.Data.(type) {
	case string:
		return value
	case float64:
		return fmt.Sprintf("%g", value)
	case bool:
		if value {
			return "true"
		} else {
			return "false"
		}
	case nil:
		return "nil"
	default:
		return ""
	}
}

func NewValue(data interface{}) *Value {
	switch value := data.(type) {
	case string:
		return &Value{STRING_DT, value}
	case float64:
		return &Value{NUMBER_DT, value}
	case bool:
		return &Value{BOOLEAN_DT, value}
	case nil:
		return &Value{NULL_DT, value}
	default:
		panic("Language fatal: Undefined datatype")
	}
}
