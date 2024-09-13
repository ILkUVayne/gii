package session

import (
	"github.com/ILkUVayne/utlis-go/v2/ulog"
)

func (s *Session) Begin() (err error) {
	ulog.Info("transaction begin")
	s.tx, err = s.db.Begin()
	return
}

func (s *Session) Commit() error {
	ulog.Info("transaction commit")
	return s.tx.Commit()
}

func (s *Session) RollBack() error {
	ulog.Info("transaction rollback")
	return s.tx.Rollback()
}
