package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/mux"

	blogSql "goblog/sql"
)

var router = mux.NewRouter()

type ArticlesData struct {
	ID          int
	Title, Body string
	URL         *url.URL
	// Errors      map[string]string
	ShowErr string
}

func (a ArticlesData) Link() string {
	URL, err := router.Get("articles.show").URL("id", strconv.Itoa(a.ID))
	checkErr(err)
	return URL.String()
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
	id := getRouteVariable("id", r)
	rs, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			checkErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功
		tem, err := template.ParseFiles("resources/views/articles/show.gohtml")
		checkErr(err)

		tem.Execute(w, rs)
		fmt.Fprint(w, "读取成功，文章标题 —— "+rs.Title)
	}
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, title FROM articles"
	rows, err := blogSql.GetDB().Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(w, "当前没有任何文章可供浏览")
		} else {
			checkErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {
		defer rows.Close()
		var articles []ArticlesData
		for rows.Next() {
			var a ArticlesData
			err := rows.Scan(&a.ID, &a.Title)
			checkErr(err)
			articles = append(articles, a)
		}
		// 2.3 检测遍历时是否发生错误
		err = rows.Err()
		checkErr(err)

		// 3. 加载模板
		tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
		checkErr(err)

		// 4. 渲染模板，将所有文章的数据传输进去
		err = tmpl.Execute(w, articles)
		checkErr(err)
	}
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	// Form：存储了 post、put 和 get 参数，在使用之前需要调用 ParseForm 方法。PostForm：存储了 post、put 参数，在使用之前需要调用 ParseForm 方法。
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("解析出错！")
		return
	}

	title, body := r.PostFormValue("title"), r.PostFormValue("body")

	var showErr string = checkArticleData(title, body)

	if len(showErr) != 0 {
		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesData{
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

	lastInsertID, err := saveArticleToDB(title, body)
	if lastInsertID > 0 {
		fmt.Fprint(w, "插入成功，ID 为"+strconv.FormatInt(lastInsertID, 10))
	} else {
		checkErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	}

	fmt.Fprint(w, "验证通过！</br>")
	fmt.Fprintf(w, "title 长度 ：%d， 内容 ： %s</br>", utf8.RuneCountInString(title), title)
	fmt.Fprintf(w, "body 长度 ：%d， 内容 ： %s</br>", utf8.RuneCountInString(body), body)
}

func saveArticleToDB(title string, body string) (int64, error) {
	db := blogSql.GetDB()

	// 变量初始化
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	// 1. 获取一个 prepare 声明语句
	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
	// 例行的错误检测
	if err != nil {
		return 0, err
	}

	// 2. 在此函数运行结束后关闭此语句，防止占用 SQL 连接
	defer stmt.Close()

	// 3. 执行请求，传参进入绑定的内容
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}

	// 4. 插入成功的话，会返回自增 ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}

	return 0, err
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
	data := ArticlesData{URL: storeURL}
	tmpl, err := template.ParseFiles("resources/views/articles/create.tmpl")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func articleEditHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	rs, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			checkErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功，显示表单
		rs.URL, _ = router.Get("articles.update").URL("id", id)
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		checkErr(err)

		err = tmpl.Execute(w, rs)
		checkErr(err)
	}
}

func articleUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	_, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 文章未找到")
		} else {
			checkErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {
		err = r.ParseForm()
		if err != nil {
			fmt.Printf("解析出错！")
			return
		}

		title, body := r.PostFormValue("title"), r.PostFormValue("body")

		var showErr string = checkArticleData(title, body)

		// 校验通过允许更新
		if len(showErr) != 0 {
			query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"
			rs, err := blogSql.GetDB().Exec(query, title, body, id)
			if err != nil {
				checkErr(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "500 服务器内部错误")
			}

			if n, _ := rs.RowsAffected(); n > 0 {
				// 跳转到文章详情页
				showURL, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprintf(w, "您未做任何更改")
			}
		} else {
			// 验证不通过，显示理由
			updateURL, _ := router.Get("articles.edit").URL("id", id)
			data := ArticlesData{
				Title:   title,
				Body:    body,
				URL:     updateURL,
				ShowErr: showErr,
			}

			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			checkErr(err)

			err = tmpl.Execute(w, data)
			checkErr(err)
		}
	}

}

func getRouteVariable(paramName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[paramName]
}

func getArticleByID(id string) (ArticlesData, error) {
	article := ArticlesData{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := blogSql.GetDB().QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// 检查提交的内容是否有效
func checkArticleData(title, body string) (showErr string) {
	// 验证标题
	if title == "" {
		showErr = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		showErr = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		showErr = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		showErr = "内容长度需大于或等于 10 个字节"
	}

	return
}

func main() {
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articleCreateHandler).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articleEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articleUpdateHandler).Methods("POST").Name("articles.update")

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
