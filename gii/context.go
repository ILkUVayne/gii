package gii

import (
	"fmt"
	"gii/render"
	"net/http"
)

type HandlerFunc func(ctx *Context)

type HandlersChain []HandlerFunc

type Context struct {
	Rw  http.ResponseWriter
	Req *http.Request

	Path   string
	Method string

	Handles HandlersChain
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Rw:     w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
	}
}

func (c *Context) Next() {
	for _, v := range c.Handles {
		v(c)
	}
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
