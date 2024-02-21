package main

import (
	"net/http"

	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"goblog/pkg/logger"
)

func main() {
	router := bootstrap.SetupRoute()
	bootstrap.SetupDB()

	// 通过命名路由获取 URL 示例
	// homeURL, _ := router.Get("home").URL()
	// fmt.Println("homeURL: ", homeURL)
	// articleURL, _ := router.Get("articles.show").URL("id", "1")
	// fmt.Println("articleURL: ", articleURL)
	err := http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
	logger.LogErr(err)
}
