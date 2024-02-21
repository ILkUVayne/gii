package main

import (
	_ "gii/demo/migrate"
	_ "github.com/go-sql-driver/mysql"

	"gii/demo/config"
	"gii/glog"
)

func main() {
	glog.SetLevel(glog.InfoLevel)
	// migrate
	// Need to connect to database,see database.yaml
	//model.Engine().NewSession().Migrate()
	config.Router().Run("localhost:8000")
}
