package routes

import (
	"goblog/app/http/controllers"
	"goblog/app/http/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterWebRoutes 注册相关路由
func RegisterWebRoutes(router *mux.Router) {
	pc := new(controllers.PagesController)
	router.HandleFunc("/about", pc.About).Methods("GET").Name("about")

	// 文章相关界面
	ac := new(controllers.ArticleController)
	router.HandleFunc("/", ac.Index).Methods("GET").Name("home")
	router.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")

	// 文章列表
	router.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")

	// 文章创建
	router.HandleFunc("/articles", ac.Store).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", ac.Create).Methods("GET").Name("articles.create")
	

	// 文章更新
	router.HandleFunc("/articles/{id:[0-9]+}/edit", ac.Edit).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", ac.Update).Methods("POST").Name("articles.update")

	// 文章删除
	router.HandleFunc("/articles/{id:[0-9]+}/delete", ac.Delete).Methods("POST").Name("articles.delete")

	auc := new(controllers.AuthController)
	// 登录/注册界面
	router.HandleFunc("/auth/register", auc.Register).Methods("GET").Name("auth.register")
	router.HandleFunc("/auth/do_register", auc.DoRegister).Methods("POST").Name("auth.doregister")
	router.HandleFunc("/auth/login", auc.Login).Methods("GET").Name("auth.login")
	router.HandleFunc("/auth/do_login", auc.DoLogin).Methods("POST").Name("auth.dologin")

	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(pc.NotFound)

	// 静态资源
	// 	PathPrefix() 匹配参数里 /css/ 前缀的 URI ， 链式调用 Handler() 指定处理器为 http.FileServer()。
	// http.FileServer() 是文件目录处理器，参数 http.Dir("./public") 是指定在此目录下寻找文件。
	router.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	router.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	// 中间件：强制内容类型为 HTML
	// router.Use(middlewares.ForceHTML)

	// 启动session
	router.Use(middlewares.StartSession)
}
