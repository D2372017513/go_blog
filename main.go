package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/mux"

	_ "goblog/sql"
)

var router = mux.NewRouter()

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	// Errors      map[string]string
	ShowErr string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章 ID："+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	// Form：存储了 post、put 和 get 参数，在使用之前需要调用 ParseForm 方法。PostForm：存储了 post、put 参数，在使用之前需要调用 ParseForm 方法。
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("解析出错！")
		return
	}

	title, body := r.PostFormValue("title"), r.PostFormValue("body")

	var showErr string = ""

	if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(body) < 3 {
		showErr = "标题或内容长度不能小于 3 字节"
	} else if utf8.RuneCountInString(title) > 40 || utf8.RuneCountInString(body) > 40 {
		showErr = "标题或内容长度不能大于 40 字节"
	}

	if len(showErr) != 0 {
		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			Title:   title,
			Body:    body,
			URL:     storeURL,
			ShowErr: showErr,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.tmpl")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
		return
	}

	fmt.Fprint(w, "验证通过！</br>")
	fmt.Fprintf(w, "title 长度 ：%d， 内容 ： %s</br>", utf8.RuneCountInString(title), title)
	fmt.Fprintf(w, "body 长度 ：%d， 内容 ： %s</br>", utf8.RuneCountInString(body), body)
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 2. 继续处理请求
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func articleCreateHandler(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{URL: storeURL}
	tmpl, err := template.ParseFiles("resources/views/articles/create.tmpl")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func main() {
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articleCreateHandler).Methods("GET").Name("article.create")

	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	// 通过命名路由获取 URL 示例
	// homeURL, _ := router.Get("home").URL()
	// fmt.Println("homeURL: ", homeURL)
	// articleURL, _ := router.Get("articles.show").URL("id", "1")
	// fmt.Println("articleURL: ", articleURL)
	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
