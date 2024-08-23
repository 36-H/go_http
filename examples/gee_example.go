package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.STRING(http.StatusOK, "Hello Helios\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"Helios"}
		c.STRING(http.StatusOK, names[100])
	})

	r.Run(":9999")
}