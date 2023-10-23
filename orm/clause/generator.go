package clause

import (
	"fmt"
	"strings"
)

type generator func(vars ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator, 0)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBy] = _orderBy
}

func genBindStr(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ",")
}

func _insert(vars ...interface{}) (string, []interface{}) {
	// INSERT INTO table_name (column1,column2,column3,...)
	tableName := vars[0]
	fields := strings.Join(vars[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), nil
}

func _values(vars ...interface{}) (string, []interface{}) {
	// VALUES ($v1_1,$v1_2 ...), ($v2_1,$2_2 ...), ...
	// [[1,"ly"],[2,"lk"]]
	var sql strings.Builder
	var sqlVars []interface{}

	sql.WriteString("VALUES ")
	// value是一个二维数组
	for idx, value := range vars {
		v := value.([]interface{})
		sql.WriteString(fmt.Sprintf("(%s)", genBindStr(len(v))))
		if idx+1 < len(vars) {
			sql.WriteString(", ")
		}
		sqlVars = append(sqlVars, v...)
	}
	return sql.String(), sqlVars
}

func _select(vars ...interface{}) (string, []interface{}) {
	// SELECT column1, column2, ... FROM table_name
	tableName := vars[0]
	fields := strings.Join(vars[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), nil
}

func _limit(vars ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", vars
}

func _where(vars ...interface{}) (string, []interface{}) {
	// WHERE condition
	// _where("id=? and name=?", 1, "ly")
	desc, values := vars[0], vars[1:]
	return fmt.Sprintf("WHERE %s", desc), values
}

func _orderBy(vars ...interface{}) (string, []interface{}) {
	// ORDER BY column1, column2, ... ASC|DESC
	return fmt.Sprintf("ORDER BY %s", vars[0]), nil
}
