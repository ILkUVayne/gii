package gii

type RouterGroup struct {
	Handlers HandlersChain
	basePath string
	engine   *Engine
	root     bool
}

func (r *RouterGroup) Group(basePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: r.combineHandlers(handlers),
		basePath: r.buildAbsolutePath(basePath),
		engine:   r.engine,
	}
}

func (r *RouterGroup) Use(Handlers ...HandlerFunc) *RouterGroup {
	r.Handlers = r.combineHandlers(Handlers)
	return r
}

func (r *RouterGroup) Get(relativePath string, handlers ...HandlerFunc) {
	r.handle("GET", relativePath, handlers)
}

func (r *RouterGroup) Post(relativePath string, handlers ...HandlerFunc) {
	r.handle("POST", relativePath, handlers)
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
