package model

import (
	"fmt"
	"gii/orm"
	"github.com/ILkUVayne/utlis-go/v2/ulog"
	"gopkg.in/yaml.v3"
	"os"
)

type config struct {
	Mysql struct {
		Driver   string `yaml:"driver"`
		UserName string `yaml:"user_name"`
		Password string `yaml:"password"`
		Protocol string `yaml:"protocol"`
		Ip       string `yaml:"ip"`
		Port     string `yaml:"port"`
		DbName   string `yaml:"db_name"`
	}
}

var engine *orm.Engine

func Engine() *orm.Engine {
	if engine != nil {
		return engine
	}
	// new
	var conf config
	f, err := os.ReadFile("./demo/config/database.yaml")
	if err != nil {
		ulog.Error(err)
	}
	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		ulog.Error(err)
	}

	engine = orm.NewEngine(conf.Mysql.Driver, fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s",
		conf.Mysql.UserName,
		conf.Mysql.Password,
		conf.Mysql.Protocol,
		conf.Mysql.Ip,
		conf.Mysql.Port,
		conf.Mysql.DbName,
	))
	return engine
}
