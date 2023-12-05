package srpc

import (
	"reflect"
	"testing"
)

type Cat struct {
}

type Args struct {
	Number1, Number2 int
}

func (c *Cat) Add(arg Args, reply *int) error {
	*reply = arg.Number1 + arg.Number2
	return nil
}

func (c *Cat) add(arg Args, reply *int) error {
	*reply = arg.Number1 + arg.Number2
	return nil
}

func TestNewServer(t *testing.T) {
	s := newService(new(Cat))
	if len(s.methods) != 1 {
		t.Errorf("service num is wrong %d", len(s.methods))
	}
	if s.name != "Cat" {
		t.Errorf("service name is wrong %s", s.name)
	}
	if s.typ != reflect.TypeOf(&Cat{}) {
		t.Errorf("service typ is wrong %T", s.typ)
	}
}

func TestServer_Call(t *testing.T) {
	s := newService(new(Cat))
	mt := s.methods["Add"]
	arg := mt.NewArg()
	reply := mt.NewReply()
	arg.Set(reflect.ValueOf(Args{Number1: 1, Number2: 2}))
	err := s.call(mt, arg, reply)
	if err != nil {
		t.Errorf("service call error: %s", err.Error())
	}
	if *reply.Interface().(*int) != 3 {
		t.Errorf("service call reply worng: %d", *reply.Interface().(*int))
	}
	if mt.NumCalls() != 1 {
		t.Errorf("service call numcalls is not update: %d", mt.NumCalls())
	}
}
