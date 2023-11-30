package srpc

import (
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
	GobEnc  = 1 << 4
	JsonEnc = 1 << 5
)

type RpcProto []byte

func CheckEnc(b []byte) (t codec.Type) {
	_ = b[0]
	switch {
	case uint8(b[0])&GobEnc == GobEnc:
		t = codec.GobType
	case uint8(b[0])&JsonEnc == JsonEnc:
		t = codec.JsonType
	default:
		glog.Error("rpc: invalid enc type")
	}
	return
}

func DefaultProtocol() RpcProto {
	return RpcProto{16, 0}
}
