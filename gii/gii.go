package gii

import (
	"log"
	"net/http"
	"strings"
)

type Engine struct {
	*Router
}

func (e *Engine) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, e))
}

func New() *Engine {
	return &Engine{
		Router: NewRouter(),
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := NewContext(w, r)
	key := strings.ToUpper(context.Method) + "-" + context.Path
	node, ok := e.methodTree[context.Method]
	// 路由不存在
	if !ok || !node.Search(context.Path) {
		context.String(http.StatusNotFound, "method not found path: %s\n", key)
		return
	}
	// 绑定handles到context
	context.Handles = e.getHandles(context.Method, context.Path)
	// 执行处理函数
	context.Next()
}
