package args

import (
	"strings"
)

// RouteHandle 路由
// authType 是否要验证登录态 0-不验证 1-验证 2-验证且必须为管理员登录
type RouteHandle struct {
	authType int8
	uri      string
	fn       interface{}
}

func NewRouteHandle(uri string, auth int8, fn interface{}) *RouteHandle {
	return &RouteHandle{
		uri:      strings.ToLower(uri),
		fn:       fn,
		authType: auth,
	}
}

// GetURI .
func (r *RouteHandle) GetURI() string {
	return r.uri
}

// GetAuthType .
func (r *RouteHandle) GetAuthType() int8 {
	return r.authType
}

// ExecRouteCBK .
func (r *RouteHandle) ExecRouteCBK(route *HTTPContext) (interface{}, error) {
	return r.fn.(RouteCBK)(route)
}

// ExecFunc .
func (r *RouteHandle) ExecFunc(route *HTTPContext) error {
	return r.fn.(func(*HTTPContext) error)(route)
}
