package gii

import (
	"net/http"
)

var anyMethods = []string{
	http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
	http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
	http.MethodTrace,
}

type RouterGroup struct {
	Handlers HandlersChain
	basePath string
	engine   *Engine
	root     bool
}

type GRoutes interface {
	Use(...HandlerFunc) GRoutes
	Group(string, ...HandlerFunc) GRoutes

	Any(string, ...HandlerFunc)
	Get(string, ...HandlerFunc)
	Post(string, ...HandlerFunc)
	Delete(string, ...HandlerFunc)
	Patch(string, ...HandlerFunc)
	Put(string, ...HandlerFunc)
	Options(string, ...HandlerFunc)
	Head(string, ...HandlerFunc)
}

func (r *RouterGroup) Group(basePath string, handlers ...HandlerFunc) GRoutes {
	return &RouterGroup{
		Handlers: r.combineHandlers(handlers),
		basePath: r.buildAbsolutePath(basePath),
		engine:   r.engine,
	}
}

func (r *RouterGroup) Use(Handlers ...HandlerFunc) GRoutes {
	r.Handlers = r.combineHandlers(Handlers)
	return r
}

func (r *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) {
	for _, method := range anyMethods {
		r.handle(method, relativePath, handlers)
	}
}

func (r *RouterGroup) Get(relativePath string, handlers ...HandlerFunc) {
	r.handle(http.MethodGet, relativePath, handlers)
}

func (r *RouterGroup) Post(relativePath string, handlers ...HandlerFunc) {
	r.handle(http.MethodPost, relativePath, handlers)
}

func (r *RouterGroup) Delete(relativePath string, handlers ...HandlerFunc) {
	r.handle(http.MethodDelete, relativePath, handlers)
}

func (r *RouterGroup) Patch(relativePath string, handlers ...HandlerFunc) {
	r.handle(http.MethodPatch, relativePath, handlers)
}

func (r *RouterGroup) Put(relativePath string, handlers ...HandlerFunc) {
	r.handle(http.MethodPut, relativePath, handlers)
}

func (r *RouterGroup) Options(relativePath string, handlers ...HandlerFunc) {
	r.handle(http.MethodOptions, relativePath, handlers)
}

func (r *RouterGroup) Head(relativePath string, handlers ...HandlerFunc) {
	r.handle(http.MethodHead, relativePath, handlers)
}

func (r *RouterGroup) handle(method, relativePath string, handlers HandlersChain) {
	// 构建绝对地址
	absolutePath := r.buildAbsolutePath(relativePath)
	// 合并中间件和处理控制器
	handlers = r.combineHandlers(handlers)
	// 添加路由（插入radix）
	r.engine.addRouter(method, absolutePath, handlers)
}

func (r *RouterGroup) buildAbsolutePath(relativePath string) string {
	return joinPaths(r.basePath, relativePath)
}

func (r *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(r.Handlers) + len(handlers)
	finalHandlers := make(HandlersChain, finalSize)
	copy(finalHandlers, r.Handlers)
	copy(finalHandlers[len(r.Handlers):], handlers)
	return finalHandlers
}
