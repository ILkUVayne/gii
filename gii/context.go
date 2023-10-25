package gii

import (
	"errors"
	"fmt"
	"gii/glog"
	"gii/render"
	"math"
	"net/http"
	"net/url"
)

const abortIndex int8 = math.MaxInt8 >> 1

type HandlerFunc func(ctx *Context)

type HandlersChain []HandlerFunc

type Context struct {
	Rw  http.ResponseWriter
	Req *http.Request

	Path   string
	Method string

	engine *Engine

	index   int8
	Handles HandlersChain

	// url参数
	params map[string]string
	// query 参数
	queryCache url.Values
	// form参数
	formCache url.Values
}

func (c *Context) reset() {
	c.Path = c.Req.URL.Path
	c.Method = c.Req.Method
	c.index = -1
	c.Handles = nil
	c.params = nil
	c.queryCache = nil
	c.formCache = nil
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

func (c *Context) Params(s string) string {
	if c.params == nil {
		return ""
	}
	p, ok := c.params[s]
	if !ok {
		return ""
	}
	return p
}

// query

func (c *Context) initQueryCache() {
	if c.queryCache != nil {
		return
	}
	if c.Req == nil {
		c.queryCache = url.Values{}
		return
	}
	c.queryCache = c.Req.URL.Query()
}

func (c *Context) GetQueryArray(key string) (value []string, ok bool) {
	c.initQueryCache()
	value, ok = c.queryCache[key]
	return
}

func (c *Context) GetQuery(key string) (string, bool) {
	value, ok := c.GetQueryArray(key)
	if ok {
		return value[0], ok
	}
	return "", ok
}

func (c *Context) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}

func (c *Context) DefaultQuery(key string, def string) string {
	value, ok := c.GetQuery(key)
	if ok {
		return value
	}
	return def
}

// form

func (c *Context) initPostFormCache() {
	if c.formCache != nil {
		return
	}
	req := c.Req
	c.formCache = make(url.Values)
	if err := req.ParseMultipartForm(c.engine.MaxMultipartMemory); err != nil {
		if !errors.Is(err, http.ErrNotMultipart) {
			glog.ErrorF("error on call ParseMultipartForm : %v", err)
		}
	}
	c.formCache = req.PostForm
}

func (c *Context) GetPostFormArray(key string) (value []string, ok bool) {
	c.initPostFormCache()
	value, ok = c.formCache[key]
	return
}

func (c *Context) GetPostForm(key string) (string, bool) {
	value, ok := c.GetPostFormArray(key)
	if ok {
		return value[0], ok
	}
	return "", ok
}

func (c *Context) PostForm(key string) string {
	value, _ := c.GetPostForm(key)
	return value
}

func (c *Context) DefaultPostForm(key string, def string) string {
	value, ok := c.GetPostForm(key)
	if ok {
		return value
	}
	return def
}
