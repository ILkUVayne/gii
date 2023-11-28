package main

import (
	_ "gii/demo/migrate"
	_ "github.com/go-sql-driver/mysql"

	"gii/demo/config"
	"gii/demo/model"
	"gii/glog"
)

func main() {
	glog.SetLevel(glog.InfoLevel)
	// migrate
	model.Engine().NewSession().Migrate()
	config.Router().Run("localhost:8000")
}
