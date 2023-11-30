package srpc

import (
	"encoding/binary"
	"gii/glog"
	"gii/srpc/codec"
)

// rpc协议
// 固定2字节
//
// 前四位表示编解码方式，可以表示15种
// GobEnc  00010000 00000000
// JsonEnc 00100000 00000000
//
// todo:后面的bit先预留
const (
	GobEnc  = 1 << 12
	JsonEnc = 1 << 13
)

type Option func(s *rpcProto)

type rpcProto struct {
	encType codec.Type
	// todo 其他协议预留
}

type RpcProto [2]byte

func CheckEnc(b RpcProto) (t codec.Type) {
	_ = b[1]
	mask := binary.BigEndian.Uint16(b[0:])
	switch {
	case mask&GobEnc == GobEnc:
		t = codec.GobType
	case mask&JsonEnc == JsonEnc:
		t = codec.JsonType
	default:
		glog.Error("rpc: invalid enc type")
	}
	return
}

func DefaultProtocol() *RpcProto {
	return ParseProtocol(EncType(codec.GobType))
}

func ParseProtocol(opts ...Option) *RpcProto {
	var mask uint16
	rp := new(rpcProto)
	for _, v := range opts {
		v(rp)
	}
	// 编解码
	switch rp.encType {
	case codec.GobType:
		mask |= GobEnc
	case codec.JsonType:
		mask |= JsonEnc
	}

	var Rp RpcProto
	binary.BigEndian.PutUint16(Rp[0:], mask)
	return &Rp
}

func EncType(t codec.Type) Option {
	return func(s *rpcProto) {
		s.encType = t
	}
}
