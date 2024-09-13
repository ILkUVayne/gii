package dialect

import (
	"fmt"
	"github.com/ILkUVayne/utlis-go/v2/ulog"
	"reflect"
	"strings"
	"time"
)

func init() {
	RegisterDialect("mysql", &mysql{})
}

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

func (m *mysql) TagOf(p reflect.StructField) map[string]any {
	res := make(map[string]any)
	tag := ""
	if v, ok := p.Tag.Lookup("orm"); ok {
		tagArr := strings.Split(v, ";")
		for _, v1 := range tagArr {
			if !strings.Contains(v1, ":") {
				// 单字段标签
				if v1 == "primaryKey" {
					res["PrimaryKey"] = p.Name
					continue
				}
				tag += v1 + " "
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

func (m *mysql) TableExistSql(tableName string) (string, []any) {
	return fmt.Sprintf("show TABLES LIKE '%s'", tableName), nil
}

func (m *mysql) AlterSql(tableName string, args ...any) (string, []any) {
	t, args := args[0], args[1:]
	t, ok := t.(AlterType)
	if !ok {
		ulog.Error("alter type not fount")
	}
	args = append([]any{tableName}, args...)
	switch t {
	case Add:
		return fmt.Sprintf("ALTER TABLE `%s` ADD %s %s", args...), nil
	case Modify:
		return fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN %s %s", args...), nil
	case Drop:
		return fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN %s", args...), nil
	}
	return "", nil
}
