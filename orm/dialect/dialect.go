package dialect

import (
	"reflect"

	"gii/glog"
)

type AlterType int

const (
	Add AlterType = iota
	Modify
	Drop
)

var dialectMaps = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TagOf(p reflect.StructField) map[string]interface{}
	TableExistSql(tableName string) (string, []interface{})
	AlterSql(tableName string, args ...interface{}) (string, []interface{})
}

func init() {
	RegisterDialect("mysql", &mysql{})
}

func RegisterDialect(name string, dialect Dialect) {
	if _, ok := dialectMaps[name]; ok {
		return
	}
	dialectMaps[name] = dialect
}

func GetDialect(name string) Dialect {
	dialect, ok := dialectMaps[name]
	if !ok {
		glog.ErrorF("%s dialect is not exist", name)
	}
	return dialect
}
