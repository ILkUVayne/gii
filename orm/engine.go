package orm

import (
	"database/sql"
	"gii/orm/dialect"
	"gii/orm/session"
	"github.com/ILkUVayne/utlis-go/v2/ulog"
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
		ulog.Error(err)
	}
	// ping db
	if err = db.Ping(); err != nil {
		ulog.Error(err)
	}
	e = &Engine{
		db:      db,
		dialect: dialect.GetDialect(driver),
	}
	ulog.Info("Connect database success")
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		ulog.Error(e)
	}
	ulog.Info("Close database success")
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
