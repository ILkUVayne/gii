package main

import (
	"gii/glog"
	"gii/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Order struct {
	Id   int `orm:"PRIMARY KEY"`
	Name string
}

func main() {
	//config.Router().Run("localhost:8000")
	glog.SetLevel(glog.InfoLevel)
	engine := orm.NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession().Model(&Order{})
	println(s.HasTable())
	//ss := []interface{}{"user"}
	//res := s.Db().QueryRow("show TABLES LIKE ?", ss...)
	//var tmp string
	//_ = res.Scan(&tmp)
	//println(tmp)
	//res := s.Raw("show TABLES LIKE 'user'").QueryRow()
	//res := s.Db().QueryRow("show TABLES LIKE ? ", "user")
	//res := s.Raw("show TABLES LIKE ?", "user").QueryRow()

	//s.Raw("CREATE TABLE book(Name text);").Exec()
	//s.Raw("CREATE TABLE book(Name text);").Exec()
	//rows := s.Raw("SELECT * FROM user").Query()
	//defer tools.Close(rows)
	//
	//for rows.Next() {
	//	var id int
	//	var name string
	//	err := rows.Scan(&id, &name)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(id, name)
	//}
}
