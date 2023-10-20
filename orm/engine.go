package orm

import (
	"database/sql"
	"gii/glog"
	"gii/orm/dialect"
	"gii/orm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine) {
	// connect db
	db, err := sql.Open(driver, source)
	if err != nil {
		glog.Error(err)
	}
	// ping db
	if err = db.Ping(); err != nil {
		glog.Error(err)
	}
	e = &Engine{
		db:      db,
		dialect: dialect.GetDialect(driver),
	}
	glog.Info("Connect database success")
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		glog.Error(e)
	}
	glog.Info("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.NewSession(e.db, e.dialect)
}
