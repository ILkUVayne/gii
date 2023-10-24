package session

import (
	"fmt"
	"gii/glog"
	"gii/orm/schema"
	"reflect"
	"strings"
)

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		glog.Error("cannot get refTable")
	}
	return s.refTable
}

func (s *Session) Model(m interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(m) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(m, s.dialect)
	}
	return s
}

func (s *Session) CreateTable() {
	table := s.RefTable()
	if s.HasTable() {
		glog.InfoF("table %s is exist", table.UnderscoreName)
		return
	}

	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("`%s` %s %s", field.Column, field.Type, field.Tag))
	}
	if table.PrimaryKey != "" {
		columns = append(columns, fmt.Sprintf("PRIMARY KEY (`%s`)", table.PrimaryKey))
	}
	columnDesc := strings.Join(columns, ",")
	s.Raw(fmt.Sprintf("CREATE TABLE `%s` (%s)", table.UnderscoreName, columnDesc)).Exec()
}

func (s *Session) DropTable() {
	if !s.HasTable() {
		glog.InfoF("table %s not exist", s.RefTable().UnderscoreName)
		return
	}
	s.Raw(fmt.Sprintf("DROP TABLE %s", s.RefTable().UnderscoreName)).Exec()
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSql(s.RefTable().UnderscoreName)
	res := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = res.Scan(&tmp)
	return tmp == s.RefTable().UnderscoreName
}
