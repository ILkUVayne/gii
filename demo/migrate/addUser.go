package migrate

import (
	"gii/demo/model"
)

type User struct {
	Id   int    `orm:"primaryKey;NOT NULL;AUTO_INCREMENT;column:id" json:"id"`
	Name string `orm:"type:varchar(50);column:name" json:"name"`
	Addr string `orm:"type:varchar(70);column:addr" json:"addr"`
}

func (u *User) GetRecordName() string {
	return "add_user"
}

func (u *User) Do() {
	e := model.Engine()
	e.NewSession().Model(&User{}).CreateTable()
}
