package common

import "fmt"

type DataType int

const (
	NUMBER_DT DataType = iota
	STRING_DT
	BOOLEAN_DT
	NIL_DT // REDUNDANT: is it still neccessary?
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

func NewValue(dataType DataType, data interface{}) *Value {
	return &Value{
		DataType: dataType,
		Data:     data,
	}
}
