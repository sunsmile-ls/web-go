package web

import (
	"net/http"
	"sync"
)

// 一种常用的GO设计模式，
// 用于确保HandlerBasedOnMap肯定实现了这个接口
var _ Handler = &HandlerBasedOnMap{}

type HandlerBasedOnMap struct {
	handlers sync.Map // 此处防止并发操作
}

func (h *HandlerBasedOnMap) Route(method string, pattern string, handler handlerFunc) error {
	// 构建key
	key := h.key(method, pattern)
	// 保存相应的key
	h.handlers.Store(key, handler)
	return nil
}

func (h *HandlerBasedOnMap) ServeHTTP(c *Context) {
	request := c.R
	// 获取相应的key 对应的值
	key := h.key(request.Method, request.URL.Path)
	handler, ok := h.handlers.Load(key)
	if !ok {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("not any router match"))
		return
	}

	handler.(handlerFunc)(c)
}

func (h *HandlerBasedOnMap) key(key string, pattern string) string {
	return key + "#" + pattern
}

func NewHandlerBasedOnMap() *HandlerBasedOnMap {
	return &HandlerBasedOnMap{}
}
