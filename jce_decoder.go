package gojce

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type Decoder struct {
	reader *bufio.Reader
	order  binary.ByteOrder
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		reader: bufio.NewReader(r),
		order:  binary.BigEndian,
	}
}

func (d *Decoder) readByte() (byte, error) {
	b := make([]byte, 1)
	n, err := d.reader.Read(b)
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, ErrBufferPeekOverflow
	}
	return b[0], nil
}

func (d *Decoder) readNBytes(n int, peek ...bool) (b []byte, err error) {
	var count int
	if len(peek) > 0 && peek[0] {
		b, err = d.reader.Peek(n)
		count = len(b)
	} else {
		b = make([]byte, n)
		count, err = d.reader.Read(b)
	}
	if err != nil {
		return nil, err
	}
	if count != n {
		return nil, ErrBufferPeekOverflow
	}
	return b, nil
}

func (d *Decoder) decodeBool(tag JceTag, required bool) (bool, error) {
	v, err := d.decodeInteger(tag, required, Int8)
	if err != nil {
		return false, err
	}
	if v > 0 {
		return true, nil
	}
	return false, nil
}

func (d *Decoder) decodeInt8(tag JceTag, required bool) (int8, error) {
	v, err := d.decodeInteger(tag, required, Int8)
	return int8(v), err
}
func (d *Decoder) decodeUint8(tag JceTag, required bool) (uint8, error) {
	v, err := d.decodeInteger(tag, required, Int16)
	return uint8(v), err
}

func (d *Decoder) decodeInt16(tag JceTag, required bool) (int16, error) {
	v, err := d.decodeInteger(tag, required, Int16)
	return int16(v), err
}
func (d *Decoder) decodeUint16(tag JceTag, required bool) (uint16, error) {
	v, err := d.decodeInteger(tag, required, Int32)
	return uint16(v), err
}
func (d *Decoder) decodeInt32(tag JceTag, required bool) (int32, error) {
	v, err := d.decodeInteger(tag, required, Int32)
	return int32(v), err
}
func (d *Decoder) decodeUint32(tag JceTag, required bool) (uint32, error) {
	v, err := d.decodeInteger(tag, required, Int64)
	return uint32(v), err
}
func (d *Decoder) decodeInt64(tag JceTag, required bool) (int64, error) {
	return d.decodeInteger(tag, required, Int64)
}

func (d *Decoder) decodeInteger(tag JceTag, required bool, typeValue JceEncodeType) (int64, error) {
	flag, headType, _, err := d.skipToTag(tag)
	if err != nil {
		return 0, err
	}
	if flag {
		if headType > typeValue && headType != Zero {
			return 0, fmt.Errorf("read 'Integer' type mismatch, tag: %d, get type: %d", tag, headType)
		}
		switch headType {
		case Zero:
			return 0, nil
		case Int8:
			next, err := d.readByte()
			if err != nil {
				return 0, err
			}
			v := int8(next)
			return int64(v), nil
		case Int16:
			v := int16(0)
			err := binary.Read(d.reader, d.order, &v)
			return int64(v), err
		case Int32:
			v := int32(0)
			err := binary.Read(d.reader, d.order, &v)
			return int64(v), err
		case Int64:
			v := int64(0)
			err := binary.Read(d.reader, d.order, &v)
			return v, err
		default:
			return 0, fmt.Errorf("read 'Integer' type mismatch, tag: %d, get type: %d", tag, headType)
		}
	} else {
		if required {
			return 0, fmt.Errorf("'%d' require field not exist, tag:%d", typeValue, tag)
		}
	}
	return 0, nil
}
func (d *Decoder) decodeFloat32(tag JceTag, required bool) (float32, error) {
	v, err := d.decodeFloatDouble(tag, required, Float32)
	return float32(v), err
}
func (d *Decoder) decodeFloat64(tag JceTag, required bool) (float64, error) {
	return d.decodeFloatDouble(tag, required, Float64)
}

func (d *Decoder) decodeFloatDouble(tag JceTag, required bool, typeValue JceEncodeType) (float64, error) {
	flag, headType, _, err := d.skipToTag(tag)
	if err != nil {
		return 0, err
	}
	if flag {
		if headType > typeValue {
			return 0, fmt.Errorf("read 'FloatDouble' type mismatch, tag: %d, get type: %d", tag, headType)
		}
		switch headType {
		case Zero:
			return 0, nil
		case Float32:
			v := float32(0)
			err := binary.Read(d.reader, d.order, &v)
			return float64(v), err
		case Float64:
			v := float64(0)
			err := binary.Read(d.reader, d.order, &v)
			return v, err
		default:
			return 0, fmt.Errorf("read 'Float32/Float64' type mismatch, tag: %d, get type: %d", tag, headType)
		}
	} else {
		if required {
			return 0, fmt.Errorf("float64 require field not exist, tag: %d", tag)
		}
	}
	return float64(0), nil
}
func (d *Decoder) decodeString(tag JceTag, required bool) (string, error) {
	flag, headType, _, err := d.skipToTag(tag)
	if err != nil {
		return "", err
	}
	if flag {
		strLen := 0
		switch headType {
		case String1:
			b, err := d.readByte()
			if err != nil {
				return "", err
			}
			strLen = int(b)
		case String4:
			len := int32(0)
			binary.Read(d.reader, d.order, &len)
			strLen = int(len)
		default:
			return "", fmt.Errorf("read 'String' type mismatch, tag: %d, get type: %d", tag, headType)
		}
		if strLen < 0 {
			return "", ErrBufferPeekOverflow
		}
		b, err := d.readNBytes(strLen)
		if err != nil {
			return "", err
		}
		return string(b), nil
	} else {
		if required {
			return "", fmt.Errorf("string require field not exist, tag:%d", tag)
		}
	}
	return "", nil
}

func (d *Decoder) decode(tag JceTag, required bool, v *reflect.Value) error {
	switch v.Type().Kind() {
	case reflect.Bool:
		b, err := d.decodeBool(tag, required)
		if err == nil {
			v.SetBool(b)
		} else {
			return err
		}
	case reflect.Int8:
		b, err := d.decodeInt8(tag, required)
		if err == nil {
			v.SetInt(int64(b))
		} else {
			return err
		}
	case reflect.Uint8:
		b, err := d.decodeUint8(tag, required)
		if err == nil {
			v.SetUint(uint64(b))
		} else {
			return err
		}
	case reflect.Int16:
		b, err := d.decodeInt16(tag, required)
		if err == nil {
			v.SetInt(int64(b))
		} else {
			return err
		}
	case reflect.Uint16:
		b, err := d.decodeUint16(tag, required)
		if err == nil {
			v.SetUint(uint64(b))
		} else {
			return err
		}
	case reflect.Int32:
		b, err := d.decodeInt32(tag, required)
		if err == nil {
			v.SetInt(int64(b))
		} else {
			return err
		}
	case reflect.Uint32:
		b, err := d.decodeUint32(tag, required)
		if err == nil {
			v.SetUint(uint64(b))
		} else {
			return err
		}
	case reflect.Int64:
		b, err := d.decodeInt64(tag, required)
		if err == nil {
			v.SetInt(int64(b))
		} else {
			return err
		}
	case reflect.Uint64:
		b, err := d.decodeInt64(tag, required)
		if err == nil {
			v.SetUint(uint64(b))
		} else {
			return err
		}
	case reflect.Float32:
		b, err := d.decodeFloat32(tag, required)
		if err == nil {
			v.SetFloat(float64(b))
		} else {
			return err
		}
	case reflect.Float64:
		b, err := d.decodeFloat64(tag, required)
		if err == nil {
			v.SetFloat(b)
		} else {
			return err
		}
	case reflect.String:
		b, err := d.decodeString(tag, required)
		if err == nil {
			v.SetString(b)
		} else {
			return err
		}
	case reflect.Array, reflect.Slice:
		if v.IsNil() {
			v.Set(reflect.MakeSlice(v.Type(), 0, 0))
		}
		switch v.Type().Elem().Kind() {
		case reflect.Uint8:
			var b []byte
			err := d.ReadBytes(&b, tag, required)
			if err != nil {
				return err
			}
			v.SetBytes(b)
			return nil
		case reflect.String:
			var sv []string
			err := d.ReadStrings(&sv, tag, required)
			if err != nil {
				return err
			}
			v.Set(reflect.ValueOf(sv))
			return nil
		default:
			flag, headType, _, err := d.skipToTag(tag)
			if err != nil {
				return err
			}
			if flag {
				switch headType {
				case List:
					vectorSize, err := d.decodeInt32(0, true)
					if err != nil {
						return err
					}
					sv := reflect.MakeSlice(v.Type(), int(vectorSize), int(vectorSize))
					for i := 0; i < int(vectorSize); i++ {
						iv := sv.Index(i)
						err = d.decode(0, true, &(iv))
						if err != nil {
							return err
						}
					}
					v.Set(sv)
				default:
					return fmt.Errorf("read 'vector' type mismatch, tag: %d, get type: %d", tag, headType)
				}
			} else {
				if required {
					return fmt.Errorf("type require field not exist, tag: %d", tag)
				}
			}
		}
	case reflect.Map:
		flag, headType, _, err := d.skipToTag(tag)
		if err != nil {
			return err
		}
		if flag {
			switch headType {
			case Map:
				mapSize, err := d.decodeInt32(0, true)
				if err != nil {
					return err
				}
				vm := reflect.MakeMap(v.Type())
				for i := 0; i < int(mapSize); i++ {
					kv := reflect.New(v.Type().Key()).Elem()
					vv := reflect.New(v.Type().Elem()).Elem()
					err = d.decode(0, true, &(kv))
					if err != nil {
						return err
					}
					err = d.decode(1, true, &(vv))
					if err != nil {
						return err
					}
					vm.SetMapIndex(kv, vv)
				}
				v.Set(vm)
			default:
				return fmt.Errorf("read 'map' type mismatch, tag: %d, get type: %d", tag, headType)
			}
		} else {
			if required {
				return fmt.Errorf("require field not exist, tag:%d", tag)
			}
		}
	case reflect.Ptr:
		if v.IsNil() {
			return &UnmarshalError{reflect.TypeOf(v)}
		}
		xv := v.Elem()
		return d.decode(tag, required, &xv)
	case reflect.Struct:
		ts, ok := v.Addr().Interface().(Struct)
		if ok {
			return d.ReadStruct(ts, tag, required)
		}
		return &UnmarshalError{reflect.TypeOf(v)}
	default:
		return &UnmarshalError{reflect.TypeOf(v)}
	}
	return nil
}
func (d *Decoder) skipToTag(tag JceTag) (bool, JceEncodeType, JceTag, error) {
	for {
		nextHeadTag, nextHeadType, len, err := d.peekTypeTag()
		if err != nil {
			return false, 0, 0, err
		}
		if nextHeadType == StructEnd || tag < nextHeadTag {
			return false, 0, 0, nil
		}
		_, err = d.readNBytes(len)
		if err != nil {
			return false, 0, 0, err
		}
		if tag == nextHeadTag {
			return true, nextHeadType, nextHeadTag, nil
		}
		d.skipField(nextHeadType)
	}
	// return false, 0, 0, nil
}

func (d *Decoder) peekTypeTag() (JceTag, JceEncodeType, int, error) {
	b, err := d.readNBytes(1, true)
	if err != nil {
		return 0, 0, 0, err
	}
	typeTag := JceTag(b[0])
	tmpTag := typeTag >> 4
	typeValue := JceEncodeType(typeTag & 0x0F)
	if tmpTag == 15 {
		b, err := d.readNBytes(2, true)
		if err != nil {
			return 0, 0, 0, err
		}
		tmpTag = JceTag(b[1])
		return tmpTag, typeValue, 2, nil
	}

	return tmpTag, typeValue, 1, nil
}

func (d *Decoder) skipOneField() error {
	_, headType, length, err := d.peekTypeTag()
	if err != nil {
		return err
	}
	d.readNBytes(length)
	return d.skipField(headType)
}

func (d *Decoder) skipToStructEnd() error {
	for {
		_, headType, len, err := d.peekTypeTag()
		if err != nil {
			return err
		}
		_, err = d.readNBytes(len)
		if err != nil {
			return err
		}
		err = d.skipField(headType)
		if err != nil {
			return err
		}
		if headType == StructEnd {
			break
		}

	}
	return nil
}

func (d *Decoder) skipField(typeValue JceEncodeType) (err error) {
	switch typeValue {
	case Int8, Int16, Int32, Int64, Float32, Float64:
		_, err = d.readNBytes(typeValue.Size())
	case String1:
		var b byte
		b, err = d.readByte()
		if err != nil {
			return
		}
		_, err = d.readNBytes(int(b))
	case String4:
		len := uint32(0)
		err = binary.Read(d.reader, d.order, &len)
		if err != nil {
			return
		}
		_, err = d.readNBytes(int(len))
	case Map:
		var size int32
		size, err = d.decodeInt32(0, true)
		if err != nil {
			return
		}
		for i := int32(0); i < (size * 2); i++ {
			err = d.skipOneField()
			if err != nil {
				return
			}
		}
	case List:
		var size int32
		size, err = d.decodeInt32(0, true)
		if err != nil {
			return
		}
		for i := int32(0); i < size; i++ {
			err = d.skipOneField()
			if err != nil {
				return
			}
		}
	case SimpleList:
		var (
			headType JceEncodeType
			len      int
			size     int32
		)
		_, headType, len, err = d.peekTypeTag()
		if err != nil {
			return
		}
		_, err = d.readNBytes(len)
		if err != nil {
			return
		}
		if headType != Int8 {
			return fmt.Errorf("skipField with invalid type, type value: %d, %d", typeValue, headType)
		}
		size, err = d.decodeInt32(0, true)
		if err != nil {
			return
		}
		_, err = d.readNBytes(int(size))
	case StructBegin:
		err = d.skipToStructEnd()
		if err != nil {
			return
		}
	case StructEnd, Zero:
		break
	default:
		return fmt.Errorf("skipField with invalid type, type value:%d", typeValue)
	}
	return
}

func (d *Decoder) ReadByte(v *byte, tag JceTag, required bool) error {
	tv, err := d.decodeInt8(tag, required)
	if err != nil {
		return err
	}
	*v = byte(tv)
	return nil
}

func (d *Decoder) ReadBool(v *bool, tag JceTag, required bool) error {
	tv, err := d.decodeInt8(tag, required)
	if err != nil {
		return err
	}
	if tv == 0 {
		*v = false
	} else {
		*v = true
	}
	return nil
}

func (d *Decoder) ReadInt8(v *int8, tag JceTag, required bool) error {
	var err error
	*v, err = d.decodeInt8(tag, required)
	return err
}
func (d *Decoder) ReadUint8(v *uint8, tag JceTag, required bool) error {
	var err error
	*v, err = d.decodeUint8(tag, required)
	return err
}
func (d *Decoder) ReadInt16(v *int16, tag JceTag, required bool) error {
	var err error
	*v, err = d.decodeInt16(tag, required)
	return err
}
func (d *Decoder) ReadUint16(v *uint16, tag JceTag, required bool) error {
	var err error
	*v, err = d.decodeUint16(tag, required)
	return err
}
func (d *Decoder) ReadUint32(v *uint32, tag JceTag, required bool) error {
	var err error
	*v, err = d.decodeUint32(tag, required)
	return err
}
func (d *Decoder) ReadInt32(v *int32, tag JceTag, required bool) error {
	var err error
	*v, err = d.decodeInt32(tag, required)
	return err
}
func (d *Decoder) ReadInt64(v *int64, tag JceTag, required bool) error {
	var err error
	*v, err = d.decodeInt64(tag, required)
	return err
}
func (d *Decoder) ReadFloat64(v *float64, tag JceTag, required bool) error {
	var err error
	*v, err = d.decodeFloat64(tag, required)
	return err
}
func (d *Decoder) ReadFloat32(v *float32, tag JceTag, required bool) error {
	var err error
	*v, err = d.decodeFloat32(tag, required)
	return err
}

func (d *Decoder) ReadString(v *string, tag JceTag, required bool) error {
	var err error
	*v, err = d.decodeString(tag, required)
	return err
}

func (d *Decoder) ReadBytes(v *[]byte, tag JceTag, required bool) error {
	flag, headType, _, err := d.skipToTag(tag)
	if err != nil {
		return err
	}
	if !flag {
		if required {
			return fmt.Errorf("require field not exist, tag:%d", tag)
		}
		return nil
	}
	if headType != SimpleList && headType != List {
		return fmt.Errorf("read 'vector<byte>' type mismatch, tag: %d, get type: %d", tag, headType)
	}
	if headType == SimpleList {
		_, cheadType, clen, err := d.peekTypeTag()
		if err != nil {
			return err
		}
		d.readNBytes(clen)
		if cheadType != Int8 {
			return fmt.Errorf("type mismatch, tag: %d, type: %d, %d", tag, headType, cheadType)
		}
	}
	vlen, err := d.decodeInt32(0, true)
	if err != nil {
		return err
	}
	*v, err = d.readNBytes(int(vlen))
	if err != nil {
		return err
	}
	return nil
}
func (d *Decoder) ReadStrings(v *[]string, tag JceTag, required bool) error {
	flag, headType, _, err := d.skipToTag(tag)
	if err != nil {
		return err
	}
	if !flag {
		if required {
			return fmt.Errorf("require field not exist, tag:%d", tag)
		}
		return nil
	}
	if headType != List {
		return fmt.Errorf("read 'vector<string>' type mismatch, tag: %d, get type: %d", tag, headType)
	}
	vlen, err := d.decodeInt32(0, true)
	if err != nil {
		return err
	}
	sv := make([]string, int(vlen))
	*v = sv
	for i := 0; i < int(vlen); i++ {
		err = d.ReadString(&(sv[i]), 0, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Decoder) ReadMap(v interface{}, tag JceTag, required bool) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &UnmarshalError{reflect.TypeOf(v)}
	}
	return d.decode(tag, required, &rv)
}

func (d *Decoder) ReadVector(v interface{}, tag JceTag, required bool) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &UnmarshalError{reflect.TypeOf(v)}
	}
	return d.decode(tag, required, &rv)
}

func (d *Decoder) ReadStruct(v Struct, tag JceTag, required bool) error {
	flag, headType, _, err := d.skipToTag(tag)
	if err != nil {
		return err
	}
	if !flag {
		if required {
			return fmt.Errorf("require field not exist, tag:%d, type %T", tag, v)
		}
		return nil
	}
	if headType != StructBegin {
		return fmt.Errorf("read 'struct' type mismatch, tag: %d, get type: %d", tag, headType)
	}
	err = v.Decode(d.reader)
	if err != nil {
		return err
	}
	d.skipToStructEnd()
	return nil
}

func (d *Decoder) Decode(v interface{}, tag JceTag, required bool) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &UnmarshalError{reflect.TypeOf(v)}
	}
	return d.decode(tag, required, &rv)
}
