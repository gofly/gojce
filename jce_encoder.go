package gojce

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type UnmarshalError struct {
	Type reflect.Type
}

func (e *UnmarshalError) Error() string {
	if e.Type == nil {
		return "Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "Unmarshal(nil " + e.Type.String() + ")"
}

// Encoder 编码器，用于序列化
type Encoder struct {
	w     *bufio.Writer
	order binary.ByteOrder
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:     bufio.NewWriter(w),
		order: binary.BigEndian,
	}
}

func (e *Encoder) Flush() error {
	return e.w.Flush()
}

func (e *Encoder) encodeHeaderTag(tag JceTag, tagType JceEncodeType) {
	if tag < 15 {
		b := byte((uint8(tag) << 4) + uint8(tagType))
		e.w.Write([]byte{b})
	} else {
		b1 := byte(tagType + 240)
		b2 := byte(tag)
		e.w.Write([]byte{b1, b2})
	}
}

func (e *Encoder) encodeTagBoolValue(tag JceTag, bv bool) error {
	if !bv {
		e.encodeHeaderTag(tag, Zero)
	} else {
		e.encodeHeaderTag(tag, Int8)
		e.w.Write([]byte{byte(1)})
	}
	return nil
}
func (e *Encoder) encodeTagInt8Value(tag JceTag, v int8) error {
	if v == 0 {
		e.encodeHeaderTag(tag, Zero)
	} else {
		e.encodeHeaderTag(tag, Int8)
		e.w.Write([]byte{byte(v)})
	}
	return nil
}
func (e *Encoder) encodeTagInt16Value(tag JceTag, v int16) error {
	if v >= (-128) && v <= 127 {
		return e.encodeTagInt8Value(tag, int8(v))
	} else {
		e.encodeHeaderTag(tag, Int16)
		binary.Write(e.w, e.order, v)
	}
	return nil
}
func (e *Encoder) encodeTagInt32Value(tag JceTag, v int32) error {
	if v >= (-32768) && v <= 32767 {
		return e.encodeTagInt16Value(tag, int16(v))
	} else {
		e.encodeHeaderTag(tag, Int32)
		binary.Write(e.w, e.order, v)
	}
	return nil
}
func (e *Encoder) encodeTagInt64Value(tag JceTag, v int64) error {
	if v >= (-2147483647-1) && v <= 2147483647 {
		return e.encodeTagInt32Value(tag, int32(v))
	} else {
		e.encodeHeaderTag(tag, Int64)
		binary.Write(e.w, e.order, v)
	}
	return nil
}
func (e *Encoder) encodeTagFloat32Value(tag JceTag, v float32) error {
	e.encodeHeaderTag(tag, Float32)
	binary.Write(e.w, e.order, v)
	return nil
}
func (e *Encoder) encodeTagFloat64Value(tag JceTag, v float64) error {
	e.encodeHeaderTag(tag, Float64)
	binary.Write(e.w, e.order, v)
	return nil
}

func (e *Encoder) encodeTagStringValue(tag JceTag, str string) error {
	if len(str) > 255 {
		e.encodeHeaderTag(tag, String4)
		slen := uint32(len(str))
		binary.Write(e.w, e.order, slen)
	} else {
		e.encodeHeaderTag(tag, String1)
		e.w.Write([]byte{byte(len(str))})
	}
	e.w.Write([]byte(str))
	return nil
}

func (e *Encoder) encodeValueWithTag(tag JceTag, v *reflect.Value) error {
	switch v.Type().Kind() {
	case reflect.Bool:
		bv := v.Bool()
		return e.encodeTagBoolValue(tag, bv)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.encodeTagInt64Value(tag, v.Int())
	case reflect.Uint8:
		return e.encodeTagInt8Value(tag, int8(v.Uint()))
	case reflect.Int8:
		return e.encodeTagInt8Value(tag, int8(v.Int()))
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.encodeTagInt64Value(tag, int64(v.Uint()))
	case reflect.String:
		str := v.String()
		return e.encodeTagStringValue(tag, str)
	case reflect.Float32:
		return e.encodeTagFloat32Value(tag, float32(v.Float()))
	case reflect.Float64:
		return e.encodeTagFloat64Value(tag, v.Float())
	case reflect.Array, reflect.Slice:
		if v.IsNil() {
			vv := reflect.MakeSlice(v.Type(), 0, 0)
			v = &vv
		}
		if v.Len() != 0 && reflect.Uint8 == v.Index(0).Type().Kind() {
			e.encodeHeaderTag(tag, SimpleList)
			e.encodeHeaderTag(0, Int8)
			e.encodeTagInt32Value(0, int32(v.Len()))
			e.w.Write(v.Bytes())
		} else {
			e.encodeHeaderTag(tag, List)
			e.encodeTagInt32Value(0, int32(v.Len()))
			for i := 0; i < v.Len(); i++ {
				vv := v.Index(i)
				ts, ok := vv.Addr().Interface().(Struct)
				if ok {
					e.WriteStruct(ts, 0)
				} else {
					e.encodeValueWithTag(0, &vv)
				}
			}
		}
		return nil
	case reflect.Map:
		e.encodeHeaderTag(tag, Map)
		if v.IsNil() {
			e.encodeTagInt32Value(0, 0)
		} else {
			ks := v.MapKeys()
			e.encodeTagInt32Value(0, int32(len(ks)))
			for i := 0; i < len(ks); i++ {
				e.encodeValueWithTag(0, &(ks[i]))
				vv := v.MapIndex(ks[i])
				e.encodeValueWithTag(1, &vv)
			}
		}
		return nil
	case reflect.Ptr:
		// XXX: 检查性能
		rv := reflect.Indirect(*v)
		return e.encodeValueWithTag(tag, &rv)
	case reflect.Interface:
		rv := reflect.ValueOf(v.Interface())
		return e.encodeValueWithTag(tag, &rv)
	case reflect.Struct:
		e.encodeHeaderTag(tag, StructBegin)
		if reflect.PtrTo(v.Type()).Implements(structType) {
			if v.CanAddr() {
				ts := v.Addr().Interface().(Struct)
				ts.Encode(e.w)
			} else {
				tmp := reflect.New(v.Type())
				tmp.Elem().Set(*v)
				ts := tmp.Interface().(Struct)
				ts.Encode(e.w)
			}
		} else {
			return fmt.Errorf("invalid type: %v to encode", v.Type())
		}
		// num := v.NumField()
		// for i := 0; i < num; i++ {
		// 	fv := v.Field(i)
		// 	tagstr := v.Type().Field(i).Tag.Get("tag")
		// 	if len(tagstr) > 0 {
		// 		tag, _ := strconv.Atoi(tagstr)
		// 		encodeValueWithTag(writer, uint8(tag), &fv)
		// 	}
		// }
		e.encodeHeaderTag(0, StructEnd)
	}
	return nil
}

type Struct interface {
	Encode(io.Writer) error
	Decode(io.Reader) error
}

func (e *Encoder) WriteStruct(v Struct, tag JceTag) error {
	e.encodeHeaderTag(tag, StructBegin)
	v.Encode(e.w)
	e.encodeHeaderTag(0, StructEnd)
	return nil
}
func (e *Encoder) WriteInt64(v int64, tag JceTag) error {
	e.encodeTagInt64Value(tag, int64(v))
	return nil
}
func (e *Encoder) WriteUint32(v uint32, tag JceTag) error {
	e.encodeTagInt64Value(tag, int64(v))
	return nil
}
func (e *Encoder) WriteInt32(v int32, tag JceTag) error {
	e.encodeTagInt64Value(tag, int64(v))
	return nil
}
func (e *Encoder) WriteUint16(v uint16, tag JceTag) error {
	e.encodeTagInt64Value(tag, int64(v))
	return nil
}
func (e *Encoder) WriteInt16(v int16, tag JceTag) error {
	e.encodeTagInt64Value(tag, int64(v))
	return nil
}
func (e *Encoder) WriteUint8(v uint8, tag JceTag) error {
	e.encodeTagInt64Value(tag, int64(v))
	return nil
}
func (e *Encoder) WriteInt8(v int8, tag JceTag) error {
	e.encodeTagInt8Value(tag, v)
	return nil
}
func (e *Encoder) WriteBool(v bool, tag JceTag) error {
	if v {
		e.encodeTagInt8Value(tag, 1)
	} else {
		e.encodeTagInt8Value(tag, 0)
	}

	return nil
}
func (e *Encoder) WriteFloat32(v float32, tag JceTag) error {
	e.encodeTagFloat32Value(tag, v)
	return nil
}
func (e *Encoder) WriteFloat64(v float64, tag JceTag) error {
	e.encodeTagFloat64Value(tag, v)
	return nil
}
func (e *Encoder) WriteByte(v byte, tag JceTag) error {
	e.encodeTagInt8Value(tag, int8(v))
	return nil
}

func (e *Encoder) WriteBytes(v []uint8, tag JceTag) error {
	e.encodeHeaderTag(tag, SimpleList)
	e.encodeHeaderTag(0, Int8)
	e.WriteInt32(int32(len(v)), 0)
	e.w.Write(v)
	return nil
}

func (e *Encoder) WriteString(v string, tag JceTag) error {
	if len(v) > 255 {
		e.encodeHeaderTag(tag, String4)
		vlen := uint32(len(v))
		binary.Write(e.w, e.order, vlen)
	} else {
		e.encodeHeaderTag(tag, String1)
		e.w.Write([]byte{byte(len(v))})
	}
	e.w.Write([]byte(v))
	return nil
}
func (e *Encoder) WriteStrings(v []string, tag JceTag) error {
	e.encodeHeaderTag(tag, List)
	e.WriteInt32(int32(len(v)), 0)
	for _, s := range v {
		e.WriteString(s, 0)
	}
	return nil
}

// XXX: 暂不支持[]uint8
func (e *Encoder) WriteVector(v interface{}, tag JceTag) error {
	val := reflect.ValueOf(v)
	//structType := reflect.TypeOf((*Struct)(nil)).Elem()
	if val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		e.encodeHeaderTag(tag, List)
		e.WriteInt32(int32(val.Len()), 0)
		for i := 0; i < val.Len(); i++ {
			vv := val.Index(i)
			ts, ok := vv.Addr().Interface().(Struct)
			if ok {
				e.WriteStruct(ts, 0)
			} else {
				e.encodeValueWithTag(0, &vv)
			}
		}
	} else {
		return ErrNotStruct
	}
	return nil
}

func (e *Encoder) WriteMap(v interface{}, tag JceTag) error {
	val := reflect.ValueOf(v)
	e.encodeValueWithTag(tag, &val)
	return nil
}

// XXX: []int8不是SimpleList
func (e *Encoder) Encode(v interface{}, tag JceTag) error {
	val := reflect.ValueOf(v)
	e.encodeValueWithTag(tag, &val)
	return nil
}
