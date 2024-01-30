package gojce

// XXX: benchmark encode and decode

import (
	"bytes"
	"encoding/hex"
	"io"
	"reflect"
	"testing"
)

type RequestPacket struct {
	IVersion     int16             `tag:"1"  required:"true"`
	CPacketType  byte              `tag:"2"  required:"true"`
	IMessageType int32             `tag:"3"  required:"true"`
	IRequestId   int64             `tag:"4"  required:"true"`
	SServantName string            `tag:"5"  required:"true"`
	SFuncName    string            `tag:"6"  required:"true"`
	SBuffer      []byte            `tag:"7"  required:"true"`
	ITimeout     int32             `tag:"8"  required:"true"`
	Context      map[string]string `tag:"9"  required:"true"`
	Status       map[string]string `tag:"10"  required:"true"`
}

func (p *RequestPacket) Encode(w io.Writer) error {
	var err error
	encoder := NewEncoder(w)
	err = encoder.WriteInt16(p.IVersion, 1)
	if nil != err {
		return err
	}
	err = encoder.WriteByte(p.CPacketType, 2)
	if nil != err {
		return err
	}
	err = encoder.WriteInt32(p.IMessageType, 3)
	if nil != err {
		return err
	}
	err = encoder.WriteInt64(p.IRequestId, 4)
	if nil != err {
		return err
	}
	err = encoder.WriteString(p.SServantName, 5)
	if nil != err {
		return err
	}
	err = encoder.WriteString(p.SFuncName, 6)
	if nil != err {
		return err
	}
	err = encoder.WriteBytes(p.SBuffer, 7)
	if nil != err {
		return err
	}
	err = encoder.WriteInt32(p.ITimeout, 8)
	if nil != err {
		return err
	}
	err = encoder.WriteMap(p.Context, 9)
	if nil != err {
		return err
	}
	err = encoder.WriteMap(p.Status, 10)
	if nil != err {
		return err
	}
	return encoder.Flush()
}

func (p *RequestPacket) Decode(r io.Reader) error {
	var err error
	decoder := NewDecoder(r)
	err = decoder.ReadInt16(&p.IVersion, 1, true)
	if nil != err {
		return err
	}
	err = decoder.ReadByte(&p.CPacketType, 2, true)
	if nil != err {
		return err
	}
	err = decoder.ReadInt32(&p.IMessageType, 3, true)
	if nil != err {
		return err
	}
	err = decoder.ReadInt64(&p.IRequestId, 4, true)
	if nil != err {
		return err
	}
	err = decoder.ReadString(&p.SServantName, 5, true)
	if nil != err {
		return err
	}
	err = decoder.ReadString(&p.SFuncName, 6, true)
	if nil != err {
		return err
	}
	err = decoder.ReadBytes(&p.SBuffer, 7, true)
	if nil != err {
		return err
	}
	err = decoder.ReadInt32(&p.ITimeout, 8, true)
	if nil != err {
		return err
	}
	err = decoder.ReadMap(&p.Context, 9, true)
	if nil != err {
		return err
	}
	err = decoder.ReadMap(&p.Status, 10, true)
	if nil != err {
		return err
	}
	return err
}

func TestCodec(t *testing.T) {
	v1 := &RequestPacket{}
	v1.IVersion = 256
	v1.SFuncName = "helloww"
	v1.IMessageType = 12456
	v1.ITimeout = 10101
	v1.SServantName = "343242342$$"
	v1.Context = make(map[string]string)
	v1.Context["AAA"] = "BBB"
	v1.SBuffer = []byte("#######")
	var buf bytes.Buffer
	err := v1.Encode(&buf)
	if nil != err {
		t.Fatalf("###%v", err)
	}
	t.Logf("####%v", buf.Len())

	v2 := &RequestPacket{}
	err = v2.Decode(&buf)
	if nil != err {
		t.Fatalf("###%v %v", err, v2)
	}
	t.Logf("####%v", v2)
}

func TestStringCodec(t *testing.T) {
	var err error

	// string test
	val1 := "hello, adam!"
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 string
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != val1 {
		t.Fatal(val2)
	}

	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.WriteString(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()

	var val3 string
	decoder2 := NewDecoder(&buf2)
	err = decoder2.ReadString(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val3 != val1 {
		t.Fatal(val3)
	}
}

func TestBoolCodec(t *testing.T) {
	var err error

	// string test
	var val1 bool = true
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 bool
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != val1 {
		t.Fatal(val2)
	}

	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.WriteBool(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()

	var val3 bool
	decoder2 := NewDecoder(&buf2)
	err = decoder2.ReadBool(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val3 != val1 {
		t.Fatal(val3)
	}
}

func TestCharCodec(t *testing.T) {
	var err error

	// string test
	var val1 byte = 'c'
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 byte
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != val1 {
		t.Fatal(val2)
	}

	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.WriteByte(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()

	var val3 byte
	decoder2 := NewDecoder(&buf2)
	err = decoder2.ReadByte(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val3 != val1 {
		t.Fatal(val3)
	}
}

func TestFloat32Codec(t *testing.T) {
	var err error

	// string test
	var val1 float32 = 134234.2342
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 float32
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != val1 {
		t.Fatal(val2)
	}

	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.WriteFloat32(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()

	var val3 float32
	decoder2 := NewDecoder(&buf2)
	err = decoder2.ReadFloat32(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val3 != val1 {
		t.Fatal(val3)
	}
}

func TestFloat64Codec(t *testing.T) {
	var err error

	// string test
	var val1 float64 = 134234.2342
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 float64
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != val1 {
		t.Fatal(val2)
	}

	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.WriteFloat64(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()

	var val3 float64
	decoder2 := NewDecoder(&buf2)
	err = decoder2.ReadFloat64(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val3 != val1 {
		t.Fatal(val3)
	}
}

func TestInt8Codec(t *testing.T) {
	var err error

	// string test
	var val1 int8 = 127
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 int8
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != val1 {
		t.Fatal(val2)
	}

	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.WriteInt8(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()

	var val3 int8
	decoder2 := NewDecoder(&buf2)
	err = decoder2.ReadInt8(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val3 != val1 {
		t.Fatal(val3)
	}
}

func TestUInt8Codec(t *testing.T) {
	var err error

	// string test
	var val1 uint8 = 127
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 uint8
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != val1 {
		t.Fatal(val2)
	}

	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.WriteInt8(int8(val1), 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()

	var val3 int8
	decoder2 := NewDecoder(&buf2)
	err = decoder2.ReadInt8(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if uint8(val3) != val1 {
		t.Fatal(val3)
	}
}

func TestInt16Codec(t *testing.T) {
	var err error

	// string test
	var val1 int16 = 22345
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 int16
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != val1 {
		t.Fatal(val2)
	}

	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.WriteInt16(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()

	var val3 int16
	decoder2 := NewDecoder(&buf2)
	err = decoder2.ReadInt16(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val3 != val1 {
		t.Fatal(val3)
	}
}

func TestUInt16Codec(t *testing.T) {
	var err error

	// string test
	var val1 uint16 = 22345
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 uint16
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != val1 {
		t.Fatal(val2)
	}

	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.WriteInt16(int16(val1), 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()

	var val3 int16
	decoder2 := NewDecoder(&buf2)
	err = decoder2.ReadInt16(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if uint16(val3) != val1 {
		t.Fatal(val3)
	}
}

func TestInt32Codec(t *testing.T) {
	var err error

	// string test
	var val1 int32 = -1000
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 int32
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != val1 {
		t.Fatal(val2)
	}

	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.WriteInt32(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()

	var val3 int32
	decoder2 := NewDecoder(&buf2)
	err = decoder2.ReadInt32(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val3 != val1 {
		t.Fatal(val3)
	}
}

func TestUInt32Codec(t *testing.T) {
	var err error

	// string test
	var val1 uint32 = 6553621
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 uint32
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != val1 {
		t.Fatal(val2)
	}

	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.WriteInt32(int32(val1), 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()

	var val3 int32
	decoder2 := NewDecoder(&buf2)
	err = decoder2.ReadInt32(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if uint32(val3) != val1 {
		t.Fatal(val3)
	}
}

func TestInt64Codec(t *testing.T) {
	var err error

	// string test
	var val1 int64 = 41596504421
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 int64
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != val1 {
		t.Fatal(val2)
	}

	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.WriteInt64(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()

	var val3 int64
	decoder2 := NewDecoder(&buf2)
	err = decoder2.ReadInt64(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val3 != val1 {
		t.Fatal(val3)
	}
}

func TestUInt64Codec(t *testing.T) {
	var err error

	// string test
	var val1 uint64 = 41596504421
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 uint64
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != val1 {
		t.Fatal(val2)
	}

	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.WriteInt64(int64(val1), 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()

	var val3 int64
	decoder2 := NewDecoder(&buf2)
	err = decoder2.ReadInt64(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if uint64(val3) != val1 {
		t.Fatal(val3)
	}
}

func TestMapCodec(t *testing.T) {
	var err error

	// string test
	val1 := map[string]map[string][]byte{
		"req_data": {
			"struct": []byte("adam_hello"),
		},
	}
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 map[string]map[string][]byte
	decoder1 := NewDecoder(&buf1)
	err = decoder1.Decode(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// if reflect.DeepEqual(val1, val2) {
	// 	t.Fatal(val2)
	// }
}

func TestBytesCodec(t *testing.T) {
	var err error

	//	// bytes test
	val1 := []byte("adam_hello")
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.WriteBytes(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 []byte
	decoder1 := NewDecoder(&buf1)
	err = decoder1.ReadBytes(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(val1, val2) {
		t.Fatal(val2)
	}

	// bytes test
	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.Encode(&val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var val3 []byte
	decoder2 := NewDecoder(&buf2)
	err = decoder2.Decode(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(val1, val3) {
		t.Fatal(val3)
	}
}

func TestInt16ArrayCodec(t *testing.T) {
	var err error

	// int16 array test
	val1 := []int16{233, 3234, 23223, 15}
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.WriteVector(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 []int16
	decoder1 := NewDecoder(&buf1)
	err = decoder1.ReadVector(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(val2) {
		t.Fatal(len(val2))
	}
	for i := range val1 {
		if val1[i] != val2[i] {
			t.Fatal("not equal")
		}
	}

	// int16 array test
	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var val3 []int16
	decoder2 := NewDecoder(&buf2)
	err = decoder2.Decode(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(val3) {
		t.Fatal(len(val3))
	}
	for i := range val1 {
		if val1[i] != val3[i] {
			t.Fatal("not equal")
		}
	}
}

func TestUInt16ArrayCodec(t *testing.T) {
	var err error

	// uint16 array test
	val1 := []uint16{233, 3234, 23223, 15}
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.WriteVector(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 []uint16
	decoder1 := NewDecoder(&buf1)
	err = decoder1.ReadVector(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(val2) {
		t.Fatal(len(val2))
	}
	for i := range val1 {
		if val1[i] != val2[i] {
			t.Fatal("not equal")
		}
	}

	// uint16 array test
	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var val3 []uint16
	decoder2 := NewDecoder(&buf2)
	err = decoder2.Decode(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(val3) {
		t.Fatal(len(val3))
	}
	for i := range val1 {
		if val1[i] != val3[i] {
			t.Fatal("not equal")
		}
	}
}

func TestUInt32ArrayCodec(t *testing.T) {
	var err error

	// uint32 array test
	val1 := []uint32{324234233, 32342234, 23223, 15}
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.WriteVector(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 []uint32
	decoder1 := NewDecoder(&buf1)
	err = decoder1.ReadVector(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(val2) {
		t.Fatal(len(val2))
	}
	for i := range val1 {
		if val1[i] != val2[i] {
			t.Fatal("not equal")
		}
	}

	// uint32 array test
	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var val3 []uint32
	decoder2 := NewDecoder(&buf2)
	err = decoder2.Decode(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(val3) {
		t.Fatal(len(val3))
	}
	for i := range val1 {
		if val1[i] != val3[i] {
			t.Fatal("not equal")
		}
	}
}

func TestInt8ArrayCodec(t *testing.T) {
	var err error

	// int8 array test
	val1 := []int8{3, 34, 23, 15}
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.WriteVector(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 []int8
	decoder1 := NewDecoder(&buf1)
	err = decoder1.ReadVector(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(val2) {
		t.Fatal(len(val2))
	}
	for i := range val1 {
		if val1[i] != val2[i] {
			t.Fatal("not equal")
		}
	}

	// int8 array test
	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var val3 []int8
	decoder2 := NewDecoder(&buf2)
	err = decoder2.Decode(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(val3) {
		t.Fatal(len(val3))
	}
	for i := range val1 {
		if val1[i] != val3[i] {
			t.Fatal("not equal")
		}
	}
}

func TestUInt8ArrayCodec(t *testing.T) {
	var err error

	// uint8 array test
	val1 := []uint8{3, 34, 23, 15}
	var buf1 bytes.Buffer
	// 注意这里必须要用 WriteBytes
	encoder1 := NewEncoder(&buf1)
	err = encoder1.WriteBytes(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 []uint8
	decoder1 := NewDecoder(&buf1)
	err = decoder1.ReadVector(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(val2) {
		t.Fatal(len(val2))
	}
	for i := range val1 {
		if val1[i] != val2[i] {
			t.Fatal("not equal")
		}
	}

	// uint8 array test
	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var val3 []uint8
	decoder2 := NewDecoder(&buf2)
	err = decoder2.Decode(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(val3) {
		t.Fatal(len(val3))
	}
	for i := range val1 {
		if val1[i] != val3[i] {
			t.Fatal("not equal")
		}
	}
}

func TestInt32ArrayCodec(t *testing.T) {
	var err error

	// int32 array test
	val1 := []int32{24234414, 42432, 4423423, 4234}
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.WriteVector(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 []int32
	decoder1 := NewDecoder(&buf1)
	err = decoder1.ReadVector(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(val2) {
		t.Fatal(len(val2))
	}
	for i := range val1 {
		if val1[i] != val2[i] {
			t.Fatal("not equal")
		}
	}

	// int32 array test
	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var val3 []int32
	decoder2 := NewDecoder(&buf2)
	err = decoder2.Decode(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(val3) {
		t.Fatal(len(val3))
	}
	for i := range val1 {
		if val1[i] != val3[i] {
			t.Fatal("not equal")
		}
	}
}

func TestInt64ArrayCodec(t *testing.T) {
	var err error

	// int64 array test
	val1 := []int64{2234234234414, 422342323432, 4423423, 4234}
	//	val1 := []int64{}
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.WriteVector(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 []int64
	decoder1 := NewDecoder(&buf1)
	err = decoder1.ReadVector(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(val2) {
		t.Fatal(len(val2))
	}
	for i := range val1 {
		if val1[i] != val2[i] {
			t.Fatal("not equal")
		}
	}

	// int64 array test
	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var val3 []int64
	decoder2 := NewDecoder(&buf2)
	err = decoder2.Decode(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(val3) {
		t.Fatal(len(val3))
	}
	for i := range val1 {
		if val1[i] != val3[i] {
			t.Fatal("not equal")
		}
	}
}

func TestUInt64ArrayCodec(t *testing.T) {
	var err error

	// uint64 array test
	val1 := []uint64{2234234234414, 422342323432, 4423423, 4234}
	//	val1 := []uint64{}
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.WriteVector(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 []uint64
	decoder1 := NewDecoder(&buf1)
	err = decoder1.ReadVector(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(val1) != len(val2) {
		t.Fatal(len(val2))
	}
	for i := range val1 {
		if val1[i] != val2[i] {
			t.Fatal("not equal")
		}
	}

	// uint64 array test
	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var val3 []uint64
	decoder2 := NewDecoder(&buf2)
	err = decoder2.Decode(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	// XXX: 用反射比较
	if len(val1) != len(val3) {
		t.Fatal(len(val3))
	}
	for i := range val1 {
		if val1[i] != val3[i] {
			t.Fatal("not equal")
		}
	}
}

func TestStructArrayCodec(t *testing.T) {
	var err error

	// uint64 array test
	val1 := []RequestPacket{
		{
			SFuncName: "hello",
			ITimeout:  10101,
		},
	}
	var buf1 bytes.Buffer
	encoder1 := NewEncoder(&buf1)
	err = encoder1.WriteVector(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder1.Flush()
	t.Log(hex.EncodeToString(buf1.Bytes()))

	var val2 []RequestPacket
	decoder1 := NewDecoder(&buf1)
	err = decoder1.ReadVector(&val2, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(val1, val2) {
		t.Fatal(val2)
	}

	// uint64 array test
	var buf2 bytes.Buffer
	encoder2 := NewEncoder(&buf2)
	err = encoder2.Encode(val1, 0)
	if err != nil {
		t.Fatal(err)
	}
	encoder2.Flush()
	t.Log(hex.EncodeToString(buf2.Bytes()))

	var val3 []RequestPacket
	decoder2 := NewDecoder(&buf2)
	err = decoder2.Decode(&val3, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(val1, val3) {
		t.Fatal(val3)
	}
}
