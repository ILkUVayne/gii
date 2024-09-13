package main

import (
	_ "gii/demo/migrate"
	"github.com/ILkUVayne/utlis-go/v2/ulog"
	_ "github.com/go-sql-driver/mysql"

	"gii/demo/config"
)

func main() {
	ulog.SetLevel(ulog.InfoLevel)
	// migrate
	// Need to connect to database,see database.yaml
	//model.Engine().NewSession().Migrate()
	config.Router().Run("localhost:8000")
}
