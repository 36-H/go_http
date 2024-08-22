package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct{
	// e.g roots["GET"] 支持模糊匹配
	roots	map[string]*node
	// e.g handlers["GET-/p/:lang/doc"]
	handlers map[string]HandlerFunc
}

func newRouter() *router{
	return &router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func (router *router) addRoute(method string, pattern string, handler HandlerFunc){
	log.Printf("Route %4s - %s", method, pattern)
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_,ok := router.roots[method]
	if !ok{
		router.roots[method] = &node{}
	}
	router.roots[method].insert(pattern,parts,0)
	router.handlers[key] = handler
}

func (router *router) getRoute(method,path string) (*node, map[string]string){
	params := make(map[string]string)
	// 解析请求路径
	searchParts := parsePattern(path)
	// 请求方式判定
	root , ok := router.roots[method]
	// 不存在当前请求方式路由
	if !ok {
		return nil,nil
	}
	//找到当前路径对应的路由节点
	node := root.search(searchParts , 0)
	if node != nil {
		//重新解析node节点的模式串
		parts := parsePattern(node.pattern)
		for index,part  := range parts {
			if part[0] == ':' {
				//匹配到 "go/doc" => ":lang/doc" => "lang" = "go"
				params[part[1:]] = searchParts[index]
			}else if part[0] == '*' && len(part) > 1{
				//匹配到 "/static/*filepath" => "static/css/index.css" => "filepath" = css/index.css
				params[part[1:]] = strings.Join(searchParts[index:],"/")
			}
		}
		return node,params
	}
	return nil,nil
}

func (router *router) handle(ctx *Context) {
	node, params := router.getRoute(ctx.Method, ctx.Path)
	if node != nil{
		ctx.Params = params
		//重新组合key "GET-go/doc" => "GET-:lang/doc"
		key := ctx.Method + "-" + node.pattern
		router.handlers[key](ctx)  
	}else{
		ctx.STRING(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Path)
	}
}

// 模式串解析
func parsePattern(pattern string) []string{
	paths := strings.Split(pattern,"/")

	parts := make([]string,0)

	for _, item := range paths {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}