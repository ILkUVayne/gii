package orm

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type UserOrder struct {
	Id   int    `orm:"primaryKey;column:id"`
	Name string `orm:"column:name;type:varchar(50)"`
}

func TestSession_CreateTable(t *testing.T) {
	engine := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession().Model(&UserOrder{})
	s.CreateTable()
	if !s.HasTable() {
		t.Error("failed to Create Order table")
	}
}
