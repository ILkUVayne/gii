package srpc

import (
	"context"
	"fmt"
	"gii/srpc/codec"
	"net"
	"strings"
	"testing"
	"time"
)

func _assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf("assertion failed: "+msg, v...))
	}
}

func TestClient_dialTimeout(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")

	t.Run("timeout", func(t *testing.T) {
		proto := ParseProtocol(
			EncType(codec.GobType),
			SetTimeout(ConnectTimeout, Second, 1),
		)
		_ = Dial("tcp", l.Addr().String(), GetRProto(proto))
	})
	t.Run("0", func(t *testing.T) {
		proto := ParseProtocol(
			EncType(codec.GobType),
			SetTimeout(ConnectTimeout, Second, 0),
		)
		_ = Dial("tcp", l.Addr().String(), GetRProto(proto))
	})
}

type Bar int

func (b Bar) Timeout(argv int, reply *int) error {
	time.Sleep(time.Second * 2)
	return nil
}

func startServer(addr chan string) {
	Register(new(Bar))
	// pick a free port
	l, _ := net.Listen("tcp", ":0")
	addr <- l.Addr().String()
	Accept(l)
}

func TestClient_Call(t *testing.T) {
	t.Parallel()
	addrCh := make(chan string)
	go startServer(addrCh)
	addr := <-addrCh
	time.Sleep(time.Second)
	t.Run("client timeout", func(t *testing.T) {
		client := DefaultDial("tcp", addr)
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		var reply int
		err := client.CallTimeout(ctx, "Bar.Timeout", 1, &reply)
		_assert(err != nil && strings.Contains(err.Error(), ctx.Err().Error()), "expect a timeout error")
	})
	t.Run("server handle timeout", func(t *testing.T) {
		proto := ParseProtocol(
			EncType(codec.GobType),
			SetTimeout(HandleTimeout, Second, 1),
		)
		client := Dial("tcp", addr, GetRProto(proto))
		var reply int
		err := client.CallTimeout(context.Background(), "Bar.Timeout", 1, &reply)
		_assert(err != nil && strings.Contains(err.Error(), "handle timeout"), "expect a timeout error")
	})
}
