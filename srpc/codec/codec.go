package codec

import (
	"io"
)

type Header struct {
	ServerMethod string // "T.method"
	Seq          uint64
	Error        string
}

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(any) error
	Write(*Header, any) error
}

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

type NewCodecFunc func(io.ReadWriteCloser) Codec

var TypeMaps map[Type]NewCodecFunc

func init() {
	TypeMaps = make(map[Type]NewCodecFunc)
	TypeMaps[GobType] = NewGobCodec
}
