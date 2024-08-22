package main

import (
	"gee"
	"net/http"
)

func main(){
	r := gee.New()
	r.GET("/", func(ctx *gee.Context){
		ctx.HTML(http.StatusOK,"<h1> Hello </h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.STRING(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}