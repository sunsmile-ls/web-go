package web

type handlerFunc func(c *Context)

type Handler interface {
	Routable
	ServeHTTP(c *Context)
}
