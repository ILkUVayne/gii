package migrate

import (
	"gii/demo/model"
	"gii/orm/session"
)

type Order struct {
	Id      int    `orm:"primaryKey;NOT NULL;AUTO_INCREMENT;column:id" json:"id"`
	UserId  int    `orm:"column:user_id" json:"user_id"`
	OrderNo string `orm:"type:varchar(50);column:order_no" json:"order_no"`
}

func init() {
	session.MigrateMappings["order"] = &Order{}
}

func (u *Order) GetRecordName() string {
	return "add_order"
}

func (u *Order) Do() {
	e := model.Engine()
	e.NewSession().Model(&Order{}).CreateTable()
}
