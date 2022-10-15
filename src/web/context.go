package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Context struct {
	R          *http.Request
	W          http.ResponseWriter
	PathParams map[string]string
}

// ReadJson 获取请求参数，转变为 json
func (c *Context) ReadJson(data interface{}) error {
	body, err := io.ReadAll(c.R.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, data)
}

// WriteJson 序列化并且写入 响应中
func (c *Context) WriteJson(status int, data interface{}) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.W.Write(bs)
	if err != nil {
		return err
	}
	c.W.WriteHeader(status)
	return nil
}

// OkJson 成功的响应
func (c *Context) OkJson(data interface{}) error {
	return c.WriteJson(http.StatusOK, data)
}

// SystemErrJson 系统错误
func (c *Context) SystemErrJson(data interface{}) error {
	// http 库里面提前定义好了各种响应码
	return c.WriteJson(http.StatusInternalServerError, data)
}

// BadRequestJson 请求参数错误
func (c *Context) BadRequestJson(data interface{}) error {
	// http 库里面提前定义好了各种响应码
	return c.WriteJson(http.StatusBadRequest, data)
}

// Reset 重置 Context
func (c *Context) Reset(r *http.Request, w http.ResponseWriter) {
	c.W = w
	c.R = r
	c.PathParams = make(map[string]string, 1)
}

// NewContext 创建上下文对象
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W:          w,
		R:          r,
		PathParams: make(map[string]string, 1),
	}
}

// newContext 池对象使用
func newContext() *Context {
	fmt.Println("create new context")
	return &Context{}
}
