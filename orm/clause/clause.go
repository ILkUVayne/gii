package clause

import (
	"strings"
)

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBy
	UPDATE
	DELETE
	COUNT
	ALTER
	COMMENT
)

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

func (c *Clause) Set(typ Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql, c.sqlVars = make(map[Type]string), make(map[Type][]interface{})
	}
	sql, sqlVal := generators[typ](vars...)
	c.sql[typ] = sql
	c.sqlVars[typ] = sqlVal
}

func (c *Clause) Build(types ...Type) (string, []interface{}) {
	var sqls []string
	var sqlVals []interface{}
	for _, typ := range types {
		if sql, ok := c.sql[typ]; ok {
			sqls = append(sqls, sql)
			sqlVals = append(sqlVals, c.sqlVars[typ]...)
		}
	}
	return strings.Join(sqls, " "), sqlVals
}
