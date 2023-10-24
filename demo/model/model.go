package model

import (
	"gii/glog"
	"gii/orm"
	"gopkg.in/yaml.v3"
	"os"
)

type ICheckTable interface {
	CheckTableExist()
}

type config struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`
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
		glog.Error(err)
	}
	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		glog.Error(err)
	}
	engine = orm.NewEngine(conf.Driver, conf.Source)
	return engine
}
