package gii

import (
	"log"
	"net/http"
	"strings"
)

type Router struct {
	Handles map[string]HandlerFunc
}

func (r *Router) Get(pattern string, handler HandlerFunc) {
	r.addRouter("GET", pattern, handler)
}

func (r *Router) Post(pattern string, handler HandlerFunc) {
	r.addRouter("POST", pattern, handler)
}

func (r *Router) addRouter(method string, pattern string, handler HandlerFunc) {
	pattern = method + "-" + pattern
	log.Printf("Router %s", pattern)
	r.Handles[pattern] = handler
}

func (r *Router) handle(c *Context) {
	key := strings.ToUpper(c.Method) + "-" + c.Path
	handler, ok := r.Handles[key]
	if !ok {
		c.String(http.StatusNotFound, "method not found path: %s\n", key)
		return
	}
	handler(c)
}
