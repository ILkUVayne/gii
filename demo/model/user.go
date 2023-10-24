package model

type User struct {
	Id   int    `orm:"primaryKey;NOT NULL;AUTO_INCREMENT;column:id" json:"id"`
	Name string `orm:"type:varchar(50);column:name" json:"name"`
	Addr string `orm:"type:varchar(70);column:addr" json:"addr"`
}

func (u *User) CheckTableExist() {
	e := Engine()
	e.NewSession().Model(&User{}).CreateTable()
}
