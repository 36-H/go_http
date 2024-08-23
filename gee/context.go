package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct{
	Writer http.ResponseWriter
	Req *http.Request
	//请求参数 已解析
	Method string
	Path string
	Params map[string]string
	//响应状态码
	StatusCode int
	//中间件
	handlers []HandlerFunc
	index int

	// 引擎指针
	engine *Engine
}

func newContext(w http.ResponseWriter, r *http.Request) *Context{
	return &Context{
		Writer: w,
		Req: r,
		Method: r.Method,
		Path: r.URL.Path,
		index: -1,
	}
}

func (context *Context)Next(){
	context.index++;
	s := len(context.handlers)
	//不是所有的handler都会调用 Next()。
	//手工调用 Next()，一般用于在请求前后各实现一些行为。如果中间件只作用于请求前，可以省略调用Next()，此种写法可以兼容 不调用Next的写法
	//并且使用c.index<s;c.index++ 也能保证中间件只执行一次
	//当中间件不调用 next函数时,通过此循环保证中间件执行顺序
	//当中间件调用next 函数时,且当中间件执行完毕,对应c.index也已经到达指定index
	//前面存在的for循环因为是通过c.index<s 做循环判断,则不会重复执行已经执行过的中间件
	for ; context.index < s; context.index++{
		context.handlers[context.index](context)
	}
}

func (context *Context) Fail(code int, err string) {
	context.index = len(context.handlers)
	context.JSON(code, H{"message": err})
}

func (context *Context) Param(key string)string{
	value := context.Params[key]
	return value
}

// 获取GET参数
func (context *Context)Query(key string) string{
	return context.Req.URL.Query().Get(key)
}

//	获取POST参数
func (context *Context)PostForm(key string) string{
	return context.Req.FormValue(key)
}

//	设置状态码
func (context *Context)Status(code int){
	context.StatusCode = code
	context.Writer.WriteHeader(code)
}

// 设置响应头
func (context *Context)SetHeader(key,value string){
	context.Writer.Header().Set(key,value)
}

// text响应
func (context *Context)STRING(code int,format string, values... interface{}){
	context.SetHeader("Content-Type","text/plain")
	context.Status(code)
	context.Writer.Write([]byte(fmt.Sprintf(format,values...)))
}

// JSON响应
func (context *Context)JSON(code int,Object interface{})  {
	context.SetHeader("Content-Type","application/json")
	context.Status(code)
	encoder := json.NewEncoder(context.Writer)
	if err := encoder.Encode(Object); err != nil{
		http.Error(context.Writer,err.Error(),500)
	}
}

// 直接写出
func (context *Context)DATA(code int,data []byte){
	context.Status(code)
	context.Writer.Write(data)
} 

// HTML响应
func (context *Context)HTML(code int, html string){
	context.SetHeader("Content-Type","text/html")
	context.Status(code)
	context.Writer.Write([]byte(html))
}

func (context *Context) HTML_TEMPLATE(code int, name string, data interface{}) {
	context.SetHeader("Content-Type", "text/html")
	context.Status(code)
	if err := context.engine.htmlTemplates.ExecuteTemplate(context.Writer, name, data); err != nil {
		context.Fail(500, err.Error())
	}
}

