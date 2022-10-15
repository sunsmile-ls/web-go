package web

import (
	"net/http"
	"sync"
)

type Routable interface {
	// Route 设定一个路由，命中该路由的会执行handlerFunc的代码
	Route(method string, pattern string, handler handlerFunc) error
}
type Server interface {
	Routable
	// Start 启动我们的服务器
	Start(address string) error
}

type sdkHttpServer struct {
	// Name 服务的名字
	Name    string
	handler Handler   // handler 只要实现了 Handler 接口即可
	root    Filter    // Filter 作为对整个应用添加的处理，即为全局处理
	ctxPool sync.Pool // 创建池减少上下文的创建
}

func (s *sdkHttpServer) Route(method string, pattern string,
	handler handlerFunc) error {
	return s.handler.Route(method, pattern, handler)
}

func (s *sdkHttpServer) Start(address string) {
	// 服务启动
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 从池中获取上下文， 重置上下文， 减少创建上下文
		c := s.ctxPool.Get().(*Context)
		defer func() {
			s.ctxPool.Put(c)
		}()
		c.Reset(r, w)
		s.root(c)
	})
	http.ListenAndServe(address, nil)
}

func NewSdkHttpServer(name string, builders ...FilterBuilder) *sdkHttpServer {
	// 构建路由的handler
	handler := NewHandlerBasedOnTree()

	root := handler.ServeHTTP
	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}

	return &sdkHttpServer{
		Name:    name,
		root:    root,
		handler: handler,
		ctxPool: sync.Pool{New: func() interface{} { // 创建池
			return newContext()
		}},
	}
}
