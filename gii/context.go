package gii

import (
	"fmt"
	"gii/render"
	"math"
	"net/http"
)

const abortIndex int8 = math.MaxInt8 >> 1

type HandlerFunc func(ctx *Context)

type HandlersChain []HandlerFunc

type Context struct {
	Rw  http.ResponseWriter
	Req *http.Request

	Path   string
	Method string

	index   int8
	Handles HandlersChain
}

func (c *Context) reset() {
	c.Path = c.Req.URL.Path
	c.Method = c.Req.Method
	c.index = -1
	c.Handles = nil
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.Handles)) {
		c.Handles[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) Status(code int) {
	c.Rw.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	c.Rw.Header().Set(key, value)
}

func (c *Context) render(code int, r render.Render) {
	// 设置httpCode,暂时这么处理
	r.WriteContentType(c.Rw)
	c.Status(code)
	err := r.Render(c.Rw)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
}

func (c *Context) String(code int, format string, values ...any) {
	c.render(code, render.String{Format: format, Data: values})
}

func (c *Context) JSON(code int, obj any) {
	c.render(code, render.JSON{Data: obj})
}

func (c *Context) XML(code int, obj any) {
	c.render(code, render.XML{Data: obj})
}

func (c *Context) Fail(code int, err string) {
	c.Abort()
	c.JSON(code, H{"message": err})
}
