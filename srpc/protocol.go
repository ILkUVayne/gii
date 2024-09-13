package srpc

import (
	"encoding/binary"
	"gii/srpc/codec"
	"github.com/ILkUVayne/utlis-go/v2/ulog"
	"time"
)

// rpc协议 [5]rpc
// 固定2字节
// rpc[0] 编解码方式
// GobEnc  00000001
// JsonEnc 00000020
//
// rpc[1:3] 连接超时时间 后3位表示时间单位,中间10位表示超时时间，前三位丢弃
// Nanosecond 00000000 00000001
// Microsecond 00000000 00000010
// Millisecond 00000000 00000011
// Second 00000000 00000100
// Minute 00000000 00000101
// Hour 00000000 00000110
//
// rpc[3:5]- 处理超时时间，存储方式同上
const (
	GobEnc  = 1
	JsonEnc = 1 << 1

	Nanosecond  = 1
	Microsecond = 1 << 1
	Millisecond = 3
	Second      = 1 << 2
	Minute      = 5
	Hour        = 3 << 1

	timeOutMask = 1023 << 3 // 00011111 11111000
)

const (
	ConnectTimeout = iota
	HandleTimeout
)

type Option func(s *deRpcProto)

type RpcProto [5]byte

func GetRProto(b *RpcProto) *RProto {
	_ = b[4]
	var deRProto deRpcProto
	// 获取编解码方式
	switch {
	case b[0]&GobEnc == GobEnc:
		deRProto.EncType = codec.GobType
	case b[0]&JsonEnc == JsonEnc:
		deRProto.EncType = codec.JsonType
	default:
		ulog.Error("rpc: invalid enc type")
	}
	// 获取连接超时时间
	ctoMask := binary.BigEndian.Uint16(b[1:3])
	deRProto.ConnectTimeout = getTimeout(ctoMask)
	// 获取处理超时时间
	hdMask := binary.BigEndian.Uint16(b[3:5])
	deRProto.HandleTimeout = getTimeout(hdMask)

	return &RProto{proto: b, deProto: &deRProto}
}

type deRpcProto struct {
	EncType codec.Type

	ConnectTimeout     time.Duration
	connectTimeout     int
	connectTimeoutType int

	HandleTimeout     time.Duration
	handleTimeout     int
	handleTimeoutType int
}

type RProto struct {
	proto   *RpcProto
	deProto *deRpcProto
}

func DefaultProtocol() *RProto {
	proto := ParseProtocol(
		EncType(codec.GobType),
		SetTimeout(ConnectTimeout, Second, 10),
		SetTimeout(HandleTimeout, Second, 10),
	)
	return GetRProto(proto)
}

func ParseProtocol(opts ...Option) *RpcProto {
	var Rp RpcProto
	rp := new(deRpcProto)
	for _, v := range opts {
		v(rp)
	}
	// 编解码
	var encMask uint8
	switch rp.EncType {
	case codec.GobType:
		encMask |= GobEnc
	case codec.JsonType:
		encMask |= JsonEnc
	}
	Rp[0] = encMask
	// connect timeout
	var ctoMask uint16
	ctoMask |= uint16(rp.connectTimeoutType)
	ctoMask |= uint16(rp.connectTimeout) << 3
	binary.BigEndian.PutUint16(Rp[1:3], ctoMask)
	// handle timeout
	var hdMask uint16
	hdMask |= uint16(rp.handleTimeoutType)
	hdMask |= uint16(rp.handleTimeout) << 3
	binary.BigEndian.PutUint16(Rp[3:5], hdMask)
	return &Rp
}

func EncType(t codec.Type) Option {
	return func(s *deRpcProto) {
		s.EncType = t
	}
}

func SetTimeout(typ int, tt int, t int) Option {
	return func(s *deRpcProto) {
		switch typ {
		case ConnectTimeout:
			s.ConnectTimeout = getTime(tt, t)
			s.connectTimeout = t
			s.connectTimeoutType = tt
		case HandleTimeout:
			s.HandleTimeout = getTime(tt, t)
			s.handleTimeout = t
			s.handleTimeoutType = tt
		}
	}
}

func getTimeout(mask uint16) time.Duration {
	timeout := int(mask&timeOutMask) >> 3
	switch {
	case mask&Nanosecond == Nanosecond:
		return getTime(Nanosecond, timeout)
	case mask&Microsecond == Microsecond:
		return getTime(Microsecond, timeout)
	case mask&Millisecond == Millisecond:
		return getTime(Millisecond, timeout)
	case mask&Second == Second:
		return getTime(Second, timeout)
	case mask&Minute == Minute:
		return getTime(Millisecond, timeout)
	case mask&Hour == Hour:
		return getTime(Hour, timeout)
	default:
		return 0
	}
}

func getTime(tt int, t int) time.Duration {
	var T int
	switch tt {
	case Nanosecond:
		T = t
	case Microsecond:
		T = t * 1000
	case Millisecond:
		T = t * 1000 * 1000
	case Second:
		T = t * 1000 * 1000 * 1000
	case Minute:
		T = t * 1000 * 1000 * 1000 * 60
	case Hour:
		T = t * 1000 * 1000 * 1000 * 60 * 60
	}
	return time.Duration(T)
}
