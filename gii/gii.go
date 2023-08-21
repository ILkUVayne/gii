package gii

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Engine struct {
	Router map[string]http.HandlerFunc
}

func (e *Engine) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, e))
}

func New() *Engine {
	return &Engine{
		Router: make(map[string]http.HandlerFunc),
	}
}

func (e *Engine) addRouter(method string, pattern string, handler http.HandlerFunc) {
	pattern = method + "-" + pattern
	e.Router[pattern] = handler
}

func (e *Engine) Get(pattern string, handler http.HandlerFunc) {
	e.addRouter("get", pattern, handler)
}

func (e *Engine) Post(pattern string, handler http.HandlerFunc) {
	e.addRouter("post", pattern, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.ToLower(r.Method) + "-" + r.URL.Path
	handle, ok := e.Router[key]
	if !ok {
		fmt.Fprintf(w, "method not found path: %s\n", key)
	}
	handle(w, r)
}
