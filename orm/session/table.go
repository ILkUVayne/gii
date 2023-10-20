package session

import (
	"fmt"
	"gii/glog"
	"gii/orm/schema"
	"gii/tools"
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
		glog.ErrorF("table %s is exist", table.Name)
	}

	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("`%s` %s %s", field.Name, field.Type, field.Tag))
	}
	if table.PrimaryKey != "" {
		columns = append(columns, fmt.Sprintf("PRIMARY KEY (`%s`)", table.PrimaryKey))
	}
	columnDesc := strings.Join(columns, ",")
	s.Raw(fmt.Sprintf("CREATE TABLE `%s` (%s)", tools.CamelCaseToUnderscore(table.Name), columnDesc)).Exec()
}

func (s *Session) DropTable() {
	s.Raw(fmt.Sprintf("DROP TABLE %s", s.RefTable().Name)).Exec()
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSql(s.RefTable().Name)
	res := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = res.Scan(&tmp)
	return tmp == s.RefTable().Name
}
