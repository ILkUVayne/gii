package session

import (
	"gii/glog"
	"gii/orm/clause"
	"gii/orm/schema"
	"gii/tools"
	"reflect"
)

func (s *Session) Set(typ clause.Type, vars ...interface{}) {
	s.clause.Set(typ, vars...)
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	s.clause.Set(clause.WHERE, append([]interface{}{desc}, args...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBy, desc)
	return s
}

func (s *Session) Insert(dest ...interface{}) (int64, error) {
	var table *schema.Schema
	var fieldNames []string
	recordValues := make([]interface{}, 0)
	for _, value := range dest {
		if table == nil {
			table = s.Model(value).RefTable()
			fieldNames = table.SaveFields(value)
		}
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.clause.Set(clause.INSERT, table.UnderscoreName, fieldNames)
	s.clause.Set(clause.VALUES, recordValues...)
	sql, sqlVars := s.clause.Build(clause.INSERT, clause.VALUES)
	res := s.Raw(sql, sqlVars...).Exec()
	return res.RowsAffected()
}

func (s *Session) All(values interface{}) {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.UnderscoreName, table.FieldColumns)
	sql, sqlVar := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBy, clause.LIMIT)
	rows := s.Raw(sql, sqlVar...).Query()

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}

		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			glog.Error(err)
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	tools.Close(rows)
}

func (s *Session) First(value interface{}) {
	s.Limit(1).All(value)
}

func (s *Session) Update(kv ...interface{}) int64 {
	UpdateMap, ok := kv[0].(map[string]interface{})
	if !ok {
		UpdateMap = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			UpdateMap[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().UnderscoreName, UpdateMap)
	sql, sqlVar := s.clause.Build(clause.UPDATE, clause.WHERE)
	res := s.Raw(sql, sqlVar...).Exec()
	affected, err := res.RowsAffected()
	if err != nil {
		glog.Error(err)
	}
	return affected
}

func (s *Session) Delete() int64 {
	s.clause.Set(clause.DELETE, s.RefTable().UnderscoreName)
	sql, sqlVar := s.clause.Build(clause.DELETE, clause.WHERE)
	res := s.Raw(sql, sqlVar).Exec()
	affected, err := res.RowsAffected()
	if err != nil {
		glog.Error(err)
	}
	return affected
}

func (s *Session) Count() int64 {
	s.clause.Set(clause.COUNT, s.RefTable().UnderscoreName)
	sql, sqlVar := s.clause.Build(clause.COUNT, clause.WHERE)
	rows := s.Raw(sql, sqlVar...).QueryRow()
	var count int64
	if err := rows.Scan(&count); err != nil {
		glog.Error(err)
	}
	return count
}
