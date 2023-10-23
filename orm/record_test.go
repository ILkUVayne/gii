package orm

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type UserAddr struct {
	Id   int    `orm:"primaryKey;column:id;NOT NULL;AUTO_INCREMENT" json:"id"`
	Addr string `orm:"column:addr;type:varchar(255)" json:"addr"`
	No   int    `orm:"column:no" json:"no"`
}

var (
	addr1 = &UserAddr{Addr: "xxxx路1号", No: 18}
	addr2 = &UserAddr{Addr: "xxxx路2号", No: 25}
	addr3 = &UserAddr{Addr: "xxxx路3号", No: 25}
)

func TestSession_Insert(t *testing.T) {
	engine := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession().Model(&UserAddr{})
	s.DropTable()
	s.CreateTable()
	i, err := s.Insert(addr1, addr2, addr3)
	if err != nil {
		t.Error(err)
	}
	if i != 3 {
		t.Error("insert num not 3")
	}
}

func TestSession_All(t *testing.T) {
	engine := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession().Model(&UserAddr{})
	var userAddr []UserAddr
	s.Where("id > ?", 1).OrderBy("id desc").Limit(2).All(&userAddr)
	if len(userAddr) != 2 {
		t.Error("failed to query all")
	}
}
