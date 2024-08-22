package gee

import (
	"log"
	"net/http"
)

type router struct{
	handlers map[string]HandlerFunc
}

func newRouter() *router{
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

func (router *router) addRoute(method string, pattern string, handler HandlerFunc){
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	router.handlers[key] = handler
}

func (router *router) handle(ctx *Context) {
	key := ctx.Method + "-" + ctx.Path
	if handler, ok := router.handlers[key]; ok {
		handler(ctx)
	} else {
		ctx.STRING(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Path)
	}
}