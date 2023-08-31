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
	Router := gii.New().Use(gii.Base())

	r1 := Router.Group("/v1").Use(gii.V1())
	{
		r1.Get("/ping", handle)
		r1.Get("/hello", reJson)
		r1.Post("/hello", reJson)
		r1.Get("/rexml", reXML)
	}

	r2 := Router.Group("/v2").Use(gii.V2())
	{
		r2.Get("/ping", handle)
		r2.Get("/hello", reJson)
		r2.Post("/hello", reJson)
		r2.Get("/rexml", reXML)
	}

	r3 := r2.Group("/v3").Use(gii.V3())
	{
		r3.Get("/ping", handle)
		r3.Get("/hello", reJson)
		r3.Post("/hello", reJson)
		r3.Get("/rexml", reXML)
	}

	Router.Get("/ping", handle)
	Router.Get("/hello", reJson)
	Router.Post("/hello", reJson)
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
	c.JSON(http.StatusMultipleChoices, gii.H{
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
