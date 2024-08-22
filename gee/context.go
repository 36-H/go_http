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

	Method string
	Path string
	Params map[string]string

	StatusCode int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context{
	return &Context{
		Writer: w,
		Req: r,
		Method: r.Method,
		Path: r.URL.Path,
	}
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

