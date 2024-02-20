package routes

import (
	"goblog/app/http/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterWebRoutes 注册相关路由
func RegisterWebRoutes(router *mux.Router) {
	pc := new(controllers.PagesController)
	router.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	router.HandleFunc("/about", pc.About).Methods("GET").Name("about")

	// 文章相关界面
	ac := new(controllers.ArticleController)
	router.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("article.show")
	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(pc.NotFound)
}
