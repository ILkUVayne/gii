package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	//config.Router().Run("localhost:8000")
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/gii")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(db)
	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(rows)

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(id, name)
	}
}
