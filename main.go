package main

import (
	"gii/gii"
	"net/http"
)

type User struct {
	Name string
	Age  int
}

func main() {
	Router := gii.New()

	Router.Get("/ping", handle)
	Router.Get("/hello", reJson)
	Router.Get("/rexml", reXML)

	Router.Run("localhost:8000")
}

func handle(c *gii.Context) {
	c.String(http.StatusOK, "s: %s", "asdaksda")
}

func reJson(c *gii.Context) {
	user := User{
		Name: "sdja",
		Age:  1,
	}
	c.JSON(http.StatusOK, gii.H{
		"code":    200,
		"message": "操作成功",
		"data":    user,
	})
}

func reXML(c *gii.Context) {
	user := User{
		Name: "sdja",
		Age:  1,
	}

	c.XML(http.StatusOK, gii.H{
		"code":    200,
		"message": "操作成功",
		"data":    user,
	})
}
