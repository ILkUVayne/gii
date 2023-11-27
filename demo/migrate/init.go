package migrate

import (
	"gii/demo/model"
)

func Migrate() {
	e := model.Engine()
	e.NewSession().Migrate()
}
