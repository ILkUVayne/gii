package controller

import (
	"gii/demo/model"
	"gii/gii"
	"net/http"
)

func AddUser(c *gii.Context) {
	engine := model.Engine()
	name, addr := c.PostForm("name"), c.PostForm("addr")
	if name == "" {
		Xml(c, []string{}, "name is empty", http.StatusBadRequest)
		return
	}
	if addr == "" {
		Json(c, []string{}, "addr is empty", http.StatusBadRequest)
		return
	}
	_, err := engine.NewSession().Insert(&model.User{
		Name: name,
		Addr: addr,
	})
	if err != nil {
		Json(c, []string{}, "server error", http.StatusInternalServerError)
		return
	}
	Json(c, []string{}, "", http.StatusOK)
}
