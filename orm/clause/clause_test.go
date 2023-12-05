package clause

import (
	"reflect"
	"testing"
)

func TestClause_Set(t *testing.T) {
	var clause Clause
	clause.Set(SELECT, "user", []string{"id", "name"})
	clause.Set(ORDERBy, "id desc")
	clause.Set(LIMIT, 4)
	clause.Set(WHERE, "name = ?", "ly")
	sql, sqlvars := clause.Build(SELECT, WHERE, ORDERBy, LIMIT)
	t.Log(sql, sqlvars)
	if sql != "SELECT id,name FROM user WHERE name = ? ORDER BY id desc LIMIT ?" {
		t.Error("sql build failed")
	}
	if !reflect.DeepEqual(sqlvars, []any{"ly", 4}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("clause_set", func(t *testing.T) {
		TestClause_Set(t)
	})
}
