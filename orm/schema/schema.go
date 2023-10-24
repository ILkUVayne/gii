package schema

import (
	"gii/glog"
	"gii/orm/dialect"
	"gii/tools"
	"go/ast"
	"reflect"
)

// tags
const (
	PK     = "primaryKey"
	TYPE   = "type"
	COLUMN = "column"
)

type Field struct {
	Name   string
	Column string
	IsPk   bool
	Type   string
	Tag    string
}

type Schema struct {
	Model          interface{}
	Name           string
	UnderscoreName string
	Fields         []*Field
	FieldNames     []string
	FieldColumns   []string
	fieldMap       map[string]*Field
	PrimaryKey     string
}

func (s *Schema) GetField(name string) *Field {
	field, ok := s.fieldMap[name]
	if !ok {
		glog.ErrorF("field %s not exist", name)
	}
	return field
}

func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}

	for _, v := range s.Fields {
		if v.IsPk && destValue.FieldByName(v.Name).IsZero() {
			continue
		}
		fieldValues = append(fieldValues, destValue.FieldByName(v.Name).Interface())
	}
	return fieldValues
}

func (s *Schema) SaveFields(dest interface{}) []string {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldNames []string

	for _, v := range s.Fields {
		if v.IsPk && destValue.FieldByName(v.Name).IsZero() {
			continue
		}
		fieldNames = append(fieldNames, v.Column)
	}
	return fieldNames
}

func Parse(dest interface{}, dialect dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:          dest,
		Name:           modelType.Name(),
		UnderscoreName: tools.CamelCaseToUnderscore(modelType.Name()),
		fieldMap:       make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name:   p.Name,
				Column: p.Name,
				Type:   dialect.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			// 根据解析的tag标签，调整field
			if tags := dialect.TagOf(p); tags != nil {
				field.Tag = tags["Tag"].(string)
				if name, ok := tags["Name"]; ok {
					field.Column = name.(string)
				}
				if typ, ok := tags["Type"]; ok {
					field.Type = typ.(string)
				}
				if pk, ok := tags["PrimaryKey"]; ok {
					field.IsPk = true
					schema.PrimaryKey = pk.(string)
				}
			}

			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, field.Name)
			schema.FieldColumns = append(schema.FieldColumns, tools.CamelCaseToUnderscore(field.Column))
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
