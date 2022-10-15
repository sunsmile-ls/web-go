package web

import "strings"

const (
	// 根节点，只有根用这个
	nodeTypeRoot = iota
	// *
	nodeTypeAny

	// 路径参数
	nodeTypeParam

	// // 正则
	// nodeTypeReg

	// 静态，完全匹配
	nodeTypeStatic
)

const any = "*"

type matchFunc func(path string, c *Context) bool

type node struct {
	path      string
	children  []*node
	handler   handlerFunc // 匹配到叶子节点，需要执行相应的 handler
	matchFunc matchFunc   // 判断节点是否匹配的方法
	pattern   string      // 原始的 pattern。注意，它不是完整的pattern，而是匹配到这个节点的pattern
	nodeType  int
}

// newStaticNode 创建静态节点
func newStaticNode(path string) *node {
	return &node{
		children: make([]*node, 0, 2),
		matchFunc: func(p string, c *Context) bool {
			return path == p && p != "*"
		},
		nodeType: nodeTypeStatic,
		pattern:  path,
	}
}

// newRootNode 创建根节点
func newRootNode(method string) *node {
	return &node{
		children: make([]*node, 0, 2),
		matchFunc: func(p string, c *Context) bool {
			panic("never call me")
		},
		nodeType: nodeTypeRoot,
		pattern:  method,
	}
}

// newParamNode 创建路径参数节点
func newParamNode(path string) *node {
	paramName := path[1:]
	return &node{
		children: make([]*node, 0, 2),
		matchFunc: func(p string, c *Context) bool {
			if c != nil {
				c.PathParams[paramName] = p
			}
			// 如果自身是一个参数路由，
			// 然后又来一个通配符，我们认为是不匹配的
			return p != any
		},
		nodeType: nodeTypeParam,
		pattern:  path,
	}
}

// newAnyNode 创建任何节点
func newAnyNode() *node {
	return &node{
		// any节点没有 children 节点，所以不需要创建 children 属性
		matchFunc: func(p string, c *Context) bool {
			// 没有孩子节点
			return true
		},
		nodeType: nodeTypeAny,
		pattern:  any,
	}
}

// newNode 创建节点
func newNode(path string) *node {
	if path == "" {
		return newAnyNode()
	}
	if strings.HasPrefix(path, ":") {
		return newParamNode(path)
	}
	return newStaticNode(path)
}
