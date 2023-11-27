package main

import (
	"gii/demo/migrate"
	_ "github.com/go-sql-driver/mysql"

	"gii/demo/config"
	"gii/glog"
)

func main() {
	glog.SetLevel(glog.InfoLevel)
	// migrate
	migrate.Migrate()
	config.Router().Run("localhost:8000")
}
