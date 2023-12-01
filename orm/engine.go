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

type TxFunc func(*session.Session) (any, error)

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

func (e *Engine) Transaction(f TxFunc) (result any, err error) {
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		return nil, err
	}

	defer func() {
		defer func() {
			if p := recover(); p != nil {
				_ = s.RollBack()
				panic(p) // re-throw panic after Rollback
			} else if err != nil {
				_ = s.RollBack() // err is non-nil; don't change it
			} else {
				defer func() {
					if err != nil {
						_ = s.RollBack()
					}
				}()
				err = s.Commit() // err is nil; if Commit returns error update err
			}
		}()
	}()

	return f(s)
}
