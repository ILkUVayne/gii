package gii

import (
	"net/http"
	"reflect"
	"testing"
)

func reXML(c *Context) {
	c.XML(http.StatusOK, H{
		"code":    200,
		"message": "操作成功",
	})
}
func reJson(c *Context) {
	c.JSON(http.StatusMultipleChoices, H{
		"code":    200,
		"message": "操作成功",
	})
}

func getHandles() HandlersChain {
	return HandlersChain{
		reXML,
		reJson,
	}
}

func TestNewRadix(t *testing.T) {
	node := NewRadix()
	ok := reflect.DeepEqual(node, &Radix{root: &radixNode{}})
	if !ok {
		t.Error("NewRadix() cannot create trie")
	}
}

func TestRadix_Insert(t *testing.T) {
	root := NewRadix()
	root.Insert("/api/user", getHandles())
	root.Insert("/api/users", getHandles())
	root.Insert("/api/book", getHandles())
	if !root.Search("/api/user", static) {
		t.Errorf("%s cannot be insert into redix", "/api/user")
	}
	if !root.Search("/api/users", static) {
		t.Errorf("%s cannot be insert into redix", "/api/users")
	}
	if !root.Search("/api/book", static) {
		t.Errorf("%s cannot be insert into redix", "/api/book")
	}
}

func TestRadix_Del(t *testing.T) {
	root := NewRadix()
	root.Insert("/api/user", getHandles())
	root.Insert("/api/users", getHandles())
	root.Insert("/api/userx", getHandles())
	root.Insert("/api/book", getHandles())
	root.Insert("/api/", getHandles())
	root.Insert("/api/:id", getHandles())
	root.Insert("/:api", getHandles())
	if !root.Search("/api/user", static) {
		t.Errorf("%s cannot be insert into redix", "/api/user")
	}
	if !root.Search("/api/users", static) {
		t.Errorf("%s cannot be insert into redix", "/api/users")
	}
	if !root.Search("/api/userx", static) {
		t.Errorf("%s cannot be insert into redix", "/api/userx")
	}
	if !root.Search("/api/book", static) {
		t.Errorf("%s cannot be insert into redix", "/api/book")
	}
	if !root.Search("/api/", static) {
		t.Errorf("%s cannot be insert into redix", "/api/")
	}
	if !root.Search("/api/:id", static) {
		t.Errorf("%s cannot be insert into redix", "/api/:id")
	}
	if !root.Search("/:api", static) {
		t.Errorf("%s cannot be insert into redix", "/:api")
	}
	if !root.Del("/api/userx") {
		t.Errorf("%s cannot be delete", "/api/userx")
	}
	if !root.Del("/api/") {
		t.Errorf("%s cannot be delete", "/api/userx")
	}
	if !root.Del("/api/:id") {
		t.Errorf("%s cannot be delete", "/api/:id")
	}
	if !root.Del("/:api") {
		t.Errorf("%s cannot be delete", "/:api")
	}
	if root.Search("/api/", static) {
		t.Errorf("%s can be search into redix", "/api/")
	}
	if root.Search("/api/userx", static) {
		t.Errorf("%s can be search into redix", "/api/userx")
	}
	if root.Search("/api/:id", static) {
		t.Errorf("%s can be search into redix", "/api/:id")
	}
	if root.Search("/:api", static) {
		t.Errorf("%s can be search into redix", "/:api")
	}
}

func TestRadix_StartWith(t *testing.T) {
	root := NewRadix()
	root.Insert("/api/user", getHandles())
	root.Insert("/api/users", getHandles())
	root.Insert("/api/userx", getHandles())
	root.Insert("/api/book", getHandles())
	root.Insert("/api/", getHandles())
	if !root.StartWith("/api/") {
		t.Errorf("cannot be match by prfix %s", "/api/")
	}
	if !root.StartWith("/api/book") {
		t.Errorf("cannot be match by prfix %s", "/api/book")
	}
}

func TestRadix_PassCnt(t *testing.T) {
	root := NewRadix()
	root.Insert("/api/user", getHandles())
	root.Insert("/api/users", getHandles())
	root.Insert("/api/userx", getHandles())
	root.Insert("/api/book", getHandles())
	root.Insert("/api/", getHandles())
	if root.PassCnt("/api/") != 5 {
		t.Errorf("prefix \"/api/\" PassCnt != %d", 5)
	}
	if root.PassCnt("/api/book") != 1 {
		t.Errorf("prefix \"/api/b\" PassCnt != %d", 1)
	}
}
