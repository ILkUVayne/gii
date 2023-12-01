package session

import (
	"gii/glog"
	"time"
)

type Migrater interface {
	GetRecordName() string
	Do()
}

type MigrateLog struct {
	Id       int    `orm:"primaryKey;NOT NULL;AUTO_INCREMENT;column:id" json:"id"`
	Name     string `orm:"type:varchar(100);column:name;UNIQUE" json:"name"`
	CreateAt int64  `orm:"column:create_at" json:"create_at"`
}

var MigrateMappings = map[string]any{}

func (s *Session) Migrate() {
	for _, m := range MigrateMappings {
		if m1, ok := m.(Migrater); ok {
			if needMigrate := s.beforeMigrate(m1); needMigrate {
				m1.Do()
				s.afterMigrateAfter(m1)
			}
		}
	}
}

func (s *Session) beforeMigrate(m Migrater) bool {
	s = s.Model(&MigrateLog{})
	if !s.HasTable() {
		return true
	}
	i := s.Where("name  = ?", m.GetRecordName()).Count()
	if i >= 1 {
		return false
	}
	return true
}

func (s *Session) afterMigrateAfter(m Migrater) {
	s = s.Model(&MigrateLog{})
	if !s.HasTable() {
		s.CreateTable()
	}
	_, err := s.Insert(&MigrateLog{
		Name:     m.GetRecordName(),
		CreateAt: time.Now().Unix(),
	})
	if err != nil {
		glog.Error(err)
	}
}
