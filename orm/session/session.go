package session

import (
	"database/sql"
	"gii/orm/clause"
	"gii/orm/dialect"
	"gii/orm/schema"
	"github.com/ILkUVayne/utlis-go/v2/ulog"
	"strings"
	"sync"
)

type Session struct {
	db       *sql.DB
	tx       *sql.Tx
	dialect  dialect.Dialect
	refTable *schema.Schema
	clause   clause.Clause
	sql      strings.Builder
	sqlVars  []any
	mux      sync.Mutex
}

type CommonDB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func NewSession(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

func (s *Session) Db() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *Session) Raw(sql string, sqlVars ...any) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, sqlVars...)
	return s
}

func (s *Session) Exec() sql.Result {
	defer s.Clear()
	ulog.Info(s.sql.String(), s.sqlVars)
	result, err := s.Db().Exec(s.sql.String(), s.sqlVars...)
	if err != nil {
		ulog.Error(err)
	}
	return result
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	ulog.Info(s.sql.String(), s.sqlVars)
	return s.Db().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) Query() *sql.Rows {
	defer s.Clear()
	ulog.Info(s.sql.String(), s.sqlVars)
	rows, err := s.Db().Query(s.sql.String(), s.sqlVars...)
	if err != nil {
		ulog.Error(err)
	}
	return rows
}
