package model

type Order struct {
	Id      int    `orm:"primaryKey;NOT NULL;AUTO_INCREMENT;column:id" json:"id"`
	UserId  int    `orm:"column:user_id" json:"user_id"`
	OrderNo string `orm:"type:varchar(50);column:order_no" json:"order_no"`
}

func (u *Order) CheckTableExist() {
	e := Engine()
	e.NewSession().Model(&Order{}).CreateTable()
}
