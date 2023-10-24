package controller

import (
	"gii/demo/model"
	"gii/gii"
	"gii/glog"
	"net/http"
)

func AddUser(c *gii.Context) {
	engine := model.Engine()
	_, err := engine.NewSession().Insert(&model.User{
		Name: "LY",
		Addr: "xxx市xxx区25号",
	})
	if err != nil {
		glog.Error(err)
	}
	c.JSON(http.StatusOK, gii.H{
		"data": []string{},
		"code": http.StatusOK,
		"msg":  "",
	})
}
