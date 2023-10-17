package controller

import (
	"gii/gii"
	"net/http"
)

type User struct {
	Name string
	Age  int
}

func Handle(c *gii.Context) {
	c.String(http.StatusOK, "s: %s", "asdaksda")
}

func ReJson(c *gii.Context) {
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

func ReXML(c *gii.Context) {
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

func PanicC(c *gii.Context) {
	s := []int{1, 2}
	println(s[3])
	c.JSON(http.StatusMultipleChoices, gii.H{
		"code":    200,
		"message": "操作成功",
	})
}
