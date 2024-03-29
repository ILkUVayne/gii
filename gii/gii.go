package gii

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

const defaultMultipartMemory = 32 << 20

type Engine struct {
	*RouterGroup

	MaxMultipartMemory int64

	pool sync.Pool

	trees methodTrees
}

func (e *Engine) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, e))
}

func New() *Engine {
	engine := &Engine{
		RouterGroup: &RouterGroup{
			basePath: "/",
			root:     true,
			Handlers: nil,
		},
		trees:              make(methodTrees, 0, 9),
		MaxMultipartMemory: defaultMultipartMemory,
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() any {
		return engine.allocateContext()
	}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func (e *Engine) allocateContext() *Context {
	return &Context{index: -1, engine: e}
}

func (e *Engine) Use(handlers ...HandlerFunc) GRoutes {
	e.RouterGroup.Use(handlers...)
	return e
}

func (e *Engine) addRouter(method string, absolutePath string, handlers HandlersChain) {
	key := method + "-" + absolutePath
	// 获取请求方法对应的radix树
	tree := e.getMethodTree(method)
	// 判断路由是否存在
	if exist, _ := tree.Search(absolutePath, static); exist {
		log.Fatalf("router %s is exisit", key)
	}
	tree.Insert(absolutePath, handlers)
	fmt.Printf("[GIN-debug] %s\n", key)
}

func (e *Engine) getMethodTree(method string) *Radix {
	for _, v := range e.trees {
		if v.method == method {
			return v.root
		}
	}
	newTree := methodTree{
		method: method,
		root:   NewRadix(),
	}
	e.trees = append(e.trees, newTree)
	return newTree.root
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get context from pool
	c := e.pool.Get().(*Context)
	c.Rw, c.Req = w, r
	c.reset()
	// do handle
	e.httpHandle(c)
	// put context to pool
	e.pool.Put(c)
}

func (e *Engine) httpHandle(c *Context) {
	key := strings.ToUpper(c.Method) + "-" + c.Path
	// 获取请求方法对应的radix树
	tree := e.getMethodTree(c.Method)
	// 判断路由是否存在
	exist, params := tree.Search(c.Path, param)
	if !exist {
		c.String(http.StatusNotFound, "method not found path: %s\n", key)
		return
	}
	c.params = params
	// 绑定handles到context
	c.Handles = tree.GetHandles(c.Path)
	// 执行处理函数
	c.Next()
}
