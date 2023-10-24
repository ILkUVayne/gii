package session

import "gii/glog"

func (s *Session) Begin() (err error) {
	glog.Info("transaction begin")
	s.tx, err = s.db.Begin()
	return
}

func (s *Session) Commit() error {
	glog.Info("transaction commit")
	return s.tx.Commit()
}

func (s *Session) RollBack() error {
	glog.Info("transaction rollback")
	return s.tx.Rollback()
}
