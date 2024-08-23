package gee

import (
	"net/http"
)

// 定义请求处理器
type HandlerFunc func(*Context)

// 路径分组
type RouterGroup struct {
	// 路由前缀
	prefix string
	// 支持中间件
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 处理引擎
type Engine struct {
	*RouterGroup
	//请求路径和对应的处理器
	router *router
	groups []*RouterGroup
}

// 构造函数 构造空的处理引擎
func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRoute(method, pattern, handler)
}

// GET
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP implements http.Handler.
// 解析路由 查找映射表
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	engine.router.handle(c)
}
