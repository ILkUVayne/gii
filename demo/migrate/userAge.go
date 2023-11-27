package migrate

import (
	"gii/demo/model"
	"gii/orm/dialect"
)

type UserAge struct {
}

func (u *UserAge) GetRecordName() string {
	return "user_add_age"
}

func (u *UserAge) Do() {
	e := model.Engine()
	e.NewSession().Model(&User{}).Comment("年龄").Alter(dialect.Add, "age", "int")
}
