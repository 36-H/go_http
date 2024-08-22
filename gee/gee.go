package gee

import (
	"net/http"
)

// 定义请求处理器
type HandlerFunc func(*Context)

// 处理引擎
type Engine struct {
	//请求路径和对应的处理器
	router *router
}

// 构造函数 构造空的处理引擎
func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method,pattern,handler)
}

// GET 路由
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST 路由
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP implements http.Handler.
// 解析路由 查找映射表
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w,r)
	engine.router.handle(c)
}
