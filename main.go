package main

import (
	_ "github.com/go-sql-driver/mysql"

	"gii/demo/config"
	"gii/glog"
)

func main() {
	glog.SetLevel(glog.InfoLevel)
	// check table
	config.CheckTable()
	config.Router().Run("localhost:8000")
}
