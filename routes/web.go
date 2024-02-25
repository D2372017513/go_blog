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
	router.HandleFunc("/articles", middlewares.Auth(ac.Store)).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", middlewares.Auth(ac.Create)).Methods("GET").Name("articles.create")

	// 文章更新
	router.HandleFunc("/articles/{id:[0-9]+}/edit", middlewares.Auth(ac.Edit)).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", middlewares.Auth(ac.Update)).Methods("POST").Name("articles.update")

	// 文章删除
	router.HandleFunc("/articles/{id:[0-9]+}/delete", middlewares.Auth(ac.Delete)).Methods("POST").Name("articles.delete")

	auc := new(controllers.AuthController)
	// 登录/注册界面
	router.HandleFunc("/auth/register", middlewares.Guest(auc.Register)).Methods("GET").Name("auth.register")
	router.HandleFunc("/auth/do_register", middlewares.Guest(auc.DoRegister)).Methods("POST").Name("auth.doregister")
	router.HandleFunc("/auth/login", middlewares.Guest(auc.Login)).Methods("GET").Name("auth.login")
	router.HandleFunc("/auth/do_login", middlewares.Guest(auc.DoLogin)).Methods("POST").Name("auth.dologin")
	router.HandleFunc("/auth/logout", middlewares.Auth(auc.Logout)).Methods("POST").Name("auth.logout")

	uc := new(controllers.UserController)
	// 用户相关界面
	router.HandleFunc("/user/{id:[0-9]+}", middlewares.Auth(uc.Show)).Methods("GET").Name("users.show")

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
