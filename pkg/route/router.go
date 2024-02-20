package route

import (
	"goblog/logger"
	"net/http"

	"github.com/gorilla/mux"
)

var route *mux.Router

// SetRoute 设置路由实例，以供 Name2URL 等函数使用
func SetRoute(r *mux.Router) {
	route = r
}

func GetRoute() *mux.Router {
	return route
}

// Name2URL 通过路由名称来获取 URL
func Name2URL(routeName string, pairs ...string) string {
	url, err := route.Get(routeName).URL(pairs...)
	if err != nil {
		logger.LogErr(err)
		return ""
	}
	return url.String()
}

// GetRouteVariable 根据字段名获取参数
func GetRouteVariable(paramName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[paramName]
}
