package gii

import (
	"log"
	"net/http"
)

type Engine struct {
	Router
}

func (e *Engine) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, e))
}

func New() *Engine {
	return &Engine{
		Router: Router{
			Handles: map[string]HandlerFunc{},
		},
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := NewContext(w, r)
	e.handle(context)
}
