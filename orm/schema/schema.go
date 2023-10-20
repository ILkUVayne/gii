package schema

import (
	"gii/glog"
	"gii/orm/dialect"
	"go/ast"
	"reflect"
	"strings"
)

// tags
const (
	PK     = "primaryKey"
	TYPE   = "type"
	COLUMN = "column"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
	PrimaryKey string
}

func (s *Schema) GetField(name string) *Field {
	field, ok := s.fieldMap[name]
	if !ok {
		glog.ErrorF("field %s not exist", name)
	}
	return field
}

func Parse(dest interface{}, dialect dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: dialect.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			// 根据解析的tag标签，调整field
			if tags := dialect.TagOf(p); tags != nil {
				field.Tag = tags["Tag"].(string)
				if name, ok := tags["Name"]; ok {
					field.Name = name.(string)
				}
				if typ, ok := tags["Type"]; ok {
					field.Type = typ.(string)
				}
				if pk, ok := tags["PrimaryKey"]; ok {
					schema.PrimaryKey = pk.(string)
				}
			}

			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

func GetTag(tags *string) map[string]interface{} {
	tagArr := strings.Split(*tags, ";")
	ts := make(map[string]interface{})
	for _, v := range tagArr {
		if strings.IndexAny("v", ":") == -1 {
			ts[v] = nil
			continue
		}
		v1 := strings.Split(v, ":")
		ts[v1[0]] = v1[1]
	}
	return ts
}
