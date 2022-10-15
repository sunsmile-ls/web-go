package main

import (
	"fmt"
	"web_go/src/web"
)

func main() {
	server := web.NewSdkHttpServer("sunsmile")
	server.Route("GET", "/hello", func(c *web.Context) {
		fmt.Fprintf(c.W, "hello %s, %s", c.R.Method, c.R.URL.Path)
	})
	server.Start(":8081")
}
