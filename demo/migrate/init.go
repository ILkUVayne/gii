package migrate

import (
	"gii/demo/model"
	"gii/orm/session"
)

func init() {
	session.MigrateMappings = map[string]interface{}{
		"user":    &User{},
		"order":   &Order{},
		"userAge": &UserAge{},
	}
}

func Migrate() {
	e := model.Engine()
	e.NewSession().Migrate()
}
