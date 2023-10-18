package main

import (
	"fmt"
	"gii/glog"
	"gii/orm"
	"gii/tools"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//config.Router().Run("localhost:8000")
	glog.SetLevel(glog.InfoLevel)
	engine := orm.NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	defer engine.Close()
	s := engine.NewSession()

	s.Raw("CREATE TABLE book(Name text);").Exec()
	s.Raw("CREATE TABLE book(Name text);").Exec()
	rows := s.Raw("SELECT * FROM user").Query()
	defer tools.Close(rows)

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		fmt.Println(id, name)
	}
}
