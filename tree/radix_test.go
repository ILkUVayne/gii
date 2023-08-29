package tree

import (
	"reflect"
	"testing"
)

func TestNewRadix(t *testing.T) {
	node := NewRadix()
	ok := reflect.DeepEqual(node, &Radix{root: &radixNode{}})
	if !ok {
		t.Error("NewRadix() cannot create trie")
	}
}

func TestRadix_Insert(t *testing.T) {
	root := NewRadix()
	root.Insert("/api/user")
	root.Insert("/api/users")
	root.Insert("/api/book")
	if !root.Search("/api/user") {
		t.Errorf("%s cannot be insert into redix", "/api/user")
	}
	if !root.Search("/api/users") {
		t.Errorf("%s cannot be insert into redix", "/api/users")
	}
	if !root.Search("/api/book") {
		t.Errorf("%s cannot be insert into redix", "/api/book")
	}
}

func TestRadix_Del(t *testing.T) {
	root := NewRadix()
	root.Insert("/api/user")
	root.Insert("/api/users")
	root.Insert("/api/userx")
	root.Insert("/api/book")
	root.Insert("/api/")
	if !root.Search("/api/user") {
		t.Errorf("%s cannot be insert into redix", "/api/user")
	}
	if !root.Search("/api/users") {
		t.Errorf("%s cannot be insert into redix", "/api/users")
	}
	if !root.Search("/api/userx") {
		t.Errorf("%s cannot be insert into redix", "/api/userx")
	}
	if !root.Search("/api/book") {
		t.Errorf("%s cannot be insert into redix", "/api/book")
	}
	if !root.Search("/api/") {
		t.Errorf("%s cannot be insert into redix", "/api/")
	}
	if !root.Del("/api/userx") {
		t.Errorf("%s cannot be delete", "/api/userx")
	}
	if !root.Del("/api/") {
		t.Errorf("%s cannot be delete", "/api/userx")
	}
	if root.Search("/api/") {
		t.Errorf("%s can be search into redix", "/api/")
	}
	if root.Search("/api/userx") {
		t.Errorf("%s can be search into redix", "/api/userx")
	}
}

func TestRadix_StartWith(t *testing.T) {
	root := NewRadix()
	root.Insert("/api/user")
	root.Insert("/api/users")
	root.Insert("/api/userx")
	root.Insert("/api/book")
	root.Insert("/api/")
	if !root.StartWith("/api/") {
		t.Errorf("cannot be match by prfix %s", "/api/")
	}
	if !root.StartWith("/api/book") {
		t.Errorf("cannot be match by prfix %s", "/api/book")
	}
}

func TestRadix_PassCnt(t *testing.T) {
	root := NewRadix()
	root.Insert("/api/user")
	root.Insert("/api/users")
	root.Insert("/api/userx")
	root.Insert("/api/book")
	root.Insert("/api/")
	if root.PassCnt("/api/") != 5 {
		t.Errorf("prefix \"/api/\" PassCnt != %d", 5)
	}
	if root.PassCnt("/api/b") != 1 {
		t.Errorf("prefix \"/api/b\" PassCnt != %d", 1)
	}
}
