package gii

import (
	"log"
)

type Router struct {
	methodTree map[string]*Radix
}

func NewRouter() *Router {
	return &Router{make(map[string]*Radix)}
}

func (r *Router) Get(pattern string, handler HandlerFunc) {
	r.addRouter("GET", pattern, handler)
}

func (r *Router) Post(pattern string, handler HandlerFunc) {
	r.addRouter("POST", pattern, handler)
}

func (r *Router) addRouter(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	log.Printf("Router %s", key)
	node, ok := r.methodTree[method]
	// 判断路由是否存在
	if ok && node.Search(pattern) {
		log.Fatalf("router %s is exisit", key)
	}
	// 添加路由
	if !ok {
		node = NewRadix()
	}
	node.Insert(pattern, HandlersChain{handler})
	r.methodTree[method] = node
}

func (r *Router) getHandles(method string, pattern string) HandlersChain {
	return r.methodTree[method].GetHandles(pattern)
}
