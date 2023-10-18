package orm

import (
	"database/sql"
	"gii/glog"
	"strings"
	"sync"
)

type Session struct {
	db      *sql.DB
	sql     strings.Builder
	sqlVars []interface{}
	mux     sync.Mutex
}

func NewSession(db *sql.DB) *Session {
	return &Session{db: db}
}

func (s *Session) Clear() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.sql.Reset()
}

func (s *Session) Db() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, sqlVars ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, sqlVars...)
	return s
}

func (s *Session) Exec() sql.Result {
	defer s.Clear()
	glog.Info(s.sql.String(), s.sqlVars)
	result, err := s.Db().Exec(s.sql.String(), s.sqlVars...)
	if err != nil {
		glog.Error(err)
	}
	return result
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	glog.Info(s.sql.String(), s.sqlVars)
	return s.db.QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) Query() *sql.Rows {
	defer s.Clear()
	glog.Info(s.sql.String(), s.sqlVars)
	rows, err := s.db.Query(s.sql.String(), s.sqlVars...)
	if err != nil {
		glog.Error(err)
	}
	return rows
}
