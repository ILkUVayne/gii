package clause

import (
	"fmt"
	"strings"
)

type generator func(vars ...any) (string, []any)

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator, 0)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBy] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
	generators[ALTER] = _alter
	generators[COMMENT] = _comment
}

func genBindStr(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ",")
}

func _insert(vars ...any) (string, []any) {
	// INSERT INTO table_name (column1,column2,column3,...)
	tableName := vars[0]
	fields := strings.Join(vars[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), nil
}

func _values(vars ...any) (string, []any) {
	// VALUES ($v1_1,$v1_2 ...), ($v2_1,$2_2 ...), ...
	// [[1,"ly"],[2,"lk"]]
	var sql strings.Builder
	var sqlVars []any

	sql.WriteString("VALUES ")
	// value是一个二维数组
	for idx, value := range vars {
		v := value.([]any)
		sql.WriteString(fmt.Sprintf("(%s)", genBindStr(len(v))))
		if idx+1 < len(vars) {
			sql.WriteString(", ")
		}
		sqlVars = append(sqlVars, v...)
	}
	return sql.String(), sqlVars
}

func _select(vars ...any) (string, []any) {
	// SELECT column1, column2, ... FROM table_name
	tableName := vars[0]
	fields := strings.Join(vars[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), nil
}

func _limit(vars ...any) (string, []any) {
	// LIMIT $num
	return "LIMIT ?", vars
}

func _where(vars ...any) (string, []any) {
	// WHERE condition
	// _where("id=? and name=?", 1, "ly")
	desc, values := vars[0], vars[1:]
	return fmt.Sprintf("WHERE %s", desc), values
}

func _orderBy(vars ...any) (string, []any) {
	// ORDER BY column1, column2, ... ASC|DESC
	return fmt.Sprintf("ORDER BY %s", vars[0]), nil
}

func _update(vars ...any) (string, []any) {
	// UPDATE table_name SET column1 = value1, column2 = value2, ...
	var desc []string
	var values []any
	tableName := vars[0]
	kv := vars[1].(map[string]any)
	for k, v := range kv {
		desc = append(desc, k+" = ?")
		values = append(values, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(desc, ",")), values
}

func _delete(vars ...any) (string, []any) {
	// DELETE FROM table_name
	return fmt.Sprintf("DELETE FROM %s", vars[0]), nil
}

func _count(vars ...any) (string, []any) {
	// SELECT COUNT(*) FROM table_name
	return _select(vars[0], []string{"count(*)"})
}

func _alter(vars ...any) (string, []any) {
	return vars[0].(string), nil
}

func _comment(vars ...any) (string, []any) {
	// comment "ss"
	return fmt.Sprintf("COMMENT '%s'", vars[0]), nil
}
