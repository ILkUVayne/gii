package schema

import (
	dialect2 "gii/orm/dialect"
	"testing"
)

type User struct {
	Id   int `orm:"primaryKey"`
	Name string
}

func TestParse(t *testing.T) {
	schema := Parse(&User{}, dialect2.GetDialect("mysql"))
	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Error("failed to parse User struct")
	}
	if schema.GetField("Id").Tag != "PRIMARY KEY" {
		t.Error("failed to parse User Tag")
	}
}
