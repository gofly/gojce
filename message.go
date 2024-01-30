package gojce

import (
	"bytes"
	"io"
)

// Message gojce消息体接口， 类似protobuf
type Message interface {
	Encode(w io.Writer) error
	Decode(r io.Reader) error
	ClassName() string
	MD5() string
	ResetDefautlt()
}

// Marshal gojce 打包函数 与标准包 json xml proto 保持一致
func Marshal(m Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := m.Encode(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal gojce 解包
func Unmarshal(data []byte, m Message) error {
	buf := bytes.NewBuffer(data)
	return m.Decode(buf)
}
