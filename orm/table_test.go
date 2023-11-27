package orm

import (
	"testing"

	"gii/orm/dialect"
)

type UserOrder struct {
	Id   int    `orm:"primaryKey;column:id"`
	Name string `orm:"column:name;type:varchar(50)"`
}

func TestSession_CreateTable(t *testing.T) {
	engine := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession().Model(&UserOrder{})
	s.DropTable()
	s.CreateTable()
	if !s.HasTable() {
		t.Error("failed to Create Order table")
	}
}

func TestSession_Alter(t *testing.T) {
	engine := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession().Model(&UserOrder{})
	s.DropTable()
	s.CreateTable()
	s.Comment("年龄").Alter(dialect.Add, "age", "integer")
	s.Comment("修改年龄").Alter(dialect.Modify, "age", "bigint")
	s.Alter(dialect.Drop, "age")
}
