package dialect

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type mysql struct{}

var _ Dialect = (*mysql)(nil)

func (m *mysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "int"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "double"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

// TagOf 解析tag标签

func (m *mysql) TagOf(p reflect.StructField) map[string]interface{} {
	res := make(map[string]interface{})
	tag := ""
	if v, ok := p.Tag.Lookup("orm"); ok {
		tagArr := strings.Split(v, ";")
		for _, v1 := range tagArr {
			if !strings.Contains(v1, ":") {
				// 单字段标签
				if v1 == "primaryKey" {
					res["PrimaryKey"] = p.Name
				}
				continue
			}
			// 标签对
			// type:varchar(50)
			v2 := strings.Split(v1, ":")
			if v2[0] == "type" {
				res["Type"] = v2[1]
				continue
			}
			// column:name
			if v2[0] == "column" {
				res["Name"] = v2[1]
				continue
			}
			tag += v2[0] + " " + v2[1] + " "
		}
		// 保证主键名正确
		if strings.Contains(v, "primaryKey") && res["Name"] != res["PrimaryKey"] {
			res["PrimaryKey"] = res["Name"]
		}
		res["Tag"] = tag
		return res
	}
	return nil
}

func (m *mysql) TableExistSql(tableName string) (string, []interface{}) {
	return fmt.Sprintf("show TABLES LIKE '%s'", tableName), nil
}
