package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"goblog/app/models/article"
	"goblog/bootstrap"
	"goblog/logger"
	"goblog/pkg/route"
)

var router *mux.Router

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
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 读取成功，显示表单
		rs.URL = route.Name2URL("articles.update", "id", id)
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogErr(err)

		err = tmpl.Execute(w, rs)
		logger.LogErr(err)
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
			logger.LogErr(err)
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

		// var showErr string = checkArticleData(title, body)
		var showErr string = ""

		// 校验通过允许更新
		if len(showErr) != 0 {
			query := "UPDATE articles SET title = ?, body = ? WHERE id = ?"
			rs, err := bootstrap.GetDB().Exec(query, title, body, id)
			if err != nil {
				logger.LogErr(err)
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
			updateURL := route.Name2URL("articles.edit", "id", id)
			data := article.ArticlesData{
				Title: title,
				Body:  body,
				URL:   updateURL,
				// ShowErr: showErr,
			}

			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogErr(err)

			err = tmpl.Execute(w, data)
			logger.LogErr(err)
		}
	}

}

func articleDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	article, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 文章未找到")
		} else {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {
		affectRow, err := article.Delete()
		if err != nil {
			// 应该是 SQL 报错了
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			if affectRow > 0 {
				url, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, url.String(), http.StatusFound)
			} else {
				// Edge case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}
	}

}

func getArticleByID(id string) (article.ArticlesData, error) {
	article := article.ArticlesData{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := bootstrap.GetDB().QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func main() {
	router = bootstrap.SetupRoute()
	bootstrap.SetupDB()

	router.HandleFunc("/articles/{id:[0-9]+}/edit", articleEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articleUpdateHandler).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articleDeleteHandler).Methods("POST").Name("articles.delete")

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	// 通过命名路由获取 URL 示例
	// homeURL, _ := router.Get("home").URL()
	// fmt.Println("homeURL: ", homeURL)
	// articleURL, _ := router.Get("articles.show").URL("id", "1")
	// fmt.Println("articleURL: ", articleURL)
	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
