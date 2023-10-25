package main

import (
	"gii/demo/config"
	"gii/glog"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	glog.SetLevel(glog.InfoLevel)
	// check table
	config.CheckTable()
	config.Router().Run("localhost:8000")
}
