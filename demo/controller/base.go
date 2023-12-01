package controller

import (
	"gii/gii"
	"net/http"
)

func getMsg(msg string, code int) string {
	if code != http.StatusOK && msg == "" {
		msg = http.StatusText(code)
	}
	return msg
}

func Json(c *gii.Context, data any, msg string, code int) {
	c.JSON(code, gii.H{
		"data": data,
		"msg":  getMsg(msg, code),
		"code": code,
	})
}

func Xml(c *gii.Context, data any, msg string, code int) {
	c.XML(code, gii.H{
		"data": data,
		"msg":  getMsg(msg, code),
		"code": code,
	})
}
