package session

import (
	"github.com/ILkUVayne/utlis-go/v2/ulog"
	"reflect"
)

// hooks constants
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

//  Hook map

type hookFunc func(v reflect.Value) func(s *Session) error

var hooks map[string]hookFunc

type IBeforeQuery interface {
	BeforeQuery(s *Session) error
}

type IAfterQuery interface {
	AfterQuery(s *Session) error
}

type IBeforeUpdate interface {
	BeforeUpdate(s *Session) error
}

type IAfterUpdate interface {
	AfterUpdate(s *Session) error
}

type IBeforeDelete interface {
	BeforeDelete(s *Session) error
}

type IAfterDelete interface {
	AfterDelete(s *Session) error
}

type IBeforeInsert interface {
	BeforeInsert(s *Session) error
}

type IAfterInsert interface {
	AfterInsert(s *Session) error
}

func init() {
	hooks = make(map[string]hookFunc)
	hooks[BeforeQuery] = _beforeQuery
	hooks[AfterQuery] = _afterQuery
	hooks[BeforeUpdate] = _beforeUpdate
	hooks[AfterUpdate] = _afterUpdate
	hooks[BeforeDelete] = _beforeDelete
	hooks[AfterDelete] = _afterDelete
	hooks[BeforeInsert] = _beforeInsert
	hooks[AfterInsert] = _afterInsert
}

func (s *Session) CallMethod(method string, v any) {
	hook, ok := hooks[method]
	if !ok {
		return
	}
	value := reflect.ValueOf(s.RefTable().Model)
	if v != nil {
		value = reflect.ValueOf(v)
	}
	f := hook(value)
	if f == nil {
		return
	}
	err := f(s)
	if err != nil {
		ulog.Error(err)
	}
}

func _beforeQuery(v reflect.Value) func(s *Session) error {
	if f, ok := v.Interface().(IBeforeQuery); ok {
		return f.BeforeQuery
	}
	return nil
}

func _afterQuery(v reflect.Value) func(s *Session) error {
	if f, ok := v.Interface().(IAfterQuery); ok {
		return f.AfterQuery
	}
	return nil
}

func _beforeUpdate(v reflect.Value) func(s *Session) error {
	if f, ok := v.Interface().(IBeforeUpdate); ok {
		return f.BeforeUpdate
	}
	return nil
}

func _afterUpdate(v reflect.Value) func(s *Session) error {
	if f, ok := v.Interface().(IAfterUpdate); ok {
		return f.AfterUpdate
	}
	return nil
}

func _beforeDelete(v reflect.Value) func(s *Session) error {
	if f, ok := v.Interface().(IBeforeDelete); ok {
		return f.BeforeDelete
	}
	return nil
}

func _afterDelete(v reflect.Value) func(s *Session) error {
	if f, ok := v.Interface().(IAfterDelete); ok {
		return f.AfterDelete
	}
	return nil
}

func _beforeInsert(v reflect.Value) func(s *Session) error {
	if f, ok := v.Interface().(IBeforeInsert); ok {
		return f.BeforeInsert
	}
	return nil
}

func _afterInsert(v reflect.Value) func(s *Session) error {
	if f, ok := v.Interface().(IAfterInsert); ok {
		return f.AfterInsert
	}
	return nil
}
