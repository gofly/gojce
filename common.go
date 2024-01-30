package gojce

import (
	"errors"
	"reflect"
)

type JceTag uint8
type JceEncodeType uint8

const (
	Int8 JceEncodeType = iota
	Int16
	Int32
	Int64
	Float32
	Float64
	String1
	String4
	Map
	List
	StructBegin
	StructEnd
	Zero
	SimpleList

	Int1 = Int8
)

func (t JceEncodeType) Size() int {
	switch t {
	case Int8:
		return 1
	case Int16:
		return 2
	case Int32:
		return 4
	case Int64:
		return 8
	case Float32:
		return 4
	case Float64:
		return 8
	}
	return -1
}

func (t JceEncodeType) String() string {
	switch t {
	case Int8:
		return "Int8"
	case Int16:
		return "Int16"
	case Int32:
		return "Int32"
	case Int64:
		return "Int64"
	case Float32:
		return "Float32"
	case Float64:
		return "Float64"
	case String1:
		return "String1"
	case String4:
		return "String4"
	case Map:
		return "Map"
	case List:
		return "List"
	case StructBegin:
		return "StructBegin"
	case StructEnd:
		return "StructEnd"
	case Zero:
		return "Zero"
	case SimpleList:
		return "SimpleList"
	}
	return "Unknown"
}

var ErrBufferPeekOverflow = errors.New("buffer overflow when peekBuf")
var ErrJceDecodeRequireNotExist = errors.New("require field not exist")
var ErrNotStruct = errors.New("invalid 'Struct' value")

var structType = reflect.TypeOf(new(Struct)).Elem()
