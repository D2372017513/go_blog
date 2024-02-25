package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"goblog/app/models"
	"goblog/app/models/article"
	"goblog/app/policies"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"goblog/types"
)

type ArticleController struct {
	BaseController
}

func (ac *ArticleController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	rs, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		view.Render(w, view.D{
			"Article":          rs,
			"CanModifyArticle": policies.CanModifyArticle(rs),
		}, "articles.show", "articles._article_meta")
	}
}

func (ac *ArticleController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := article.GetAll()
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		view.Render(w, view.D{
			"Articles": articles,
		}, "articles.index", "articles._article_meta")
	}
}

func (ac *ArticleController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "articles.create", "articles._form_field")
}

func (ac *ArticleController) Store(w http.ResponseWriter, r *http.Request) {
	// Form：存储了 post、put 和 get 参数，在使用之前需要调用 ParseForm 方法。PostForm：存储了 post、put 参数，在使用之前需要调用 ParseForm 方法。
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("解析出错！")
		return
	}

	title, body := r.PostFormValue("title"), r.PostFormValue("body")
	data := article.ArticlesData{
		Title:  title,
		Body:   body,
		UserID: auth.User().ID,
	}
	errors := requests.ValidateArticleForm(data)

	if len(errors) != 0 {
		view.Render(w, view.D{
			"Article": data,
			"Errors":  errors,
		}, "articles.create", "articles._form_field")
		return
	}

	err = data.Create()
	if data.ID > 0 {
		fmt.Fprint(w, "插入成功，ID 为"+strconv.FormatInt(data.ID, 10))
	} else {
		logger.LogErr(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	}
}

func (ac *ArticleController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	rs, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		if !policies.CanModifyArticle(rs) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			view.Render(w, view.D{
				"Article": rs,
			}, "articles.edit", "articles._form_field")
		}
	}
}

func (ac *ArticleController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		err = r.ParseForm()
		if err != nil {
			fmt.Printf("解析出错！")
			return
		}

		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {

			title, body := r.PostFormValue("title"), r.PostFormValue("body")
			updateURL := route.Name2URL("articles.edit", "id", id)
			data := article.ArticlesData{
				BaseModel: models.BaseModel{ID: int64(types.StringToInt64(id))},
				Title:     title,
				Body:      body,
				URL:       updateURL,
			}
			errors := requests.ValidateArticleForm(data)

			// 校验通过允许更新
			if len(errors) == 0 {
				rowsAffected, err := data.Update()
				if err != nil {
					logger.LogErr(err)
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "500 服务器内部错误")
				}

				if rowsAffected > 0 {
					// 跳转到文章详情页
					showURL := route.Name2URL("articles.show", "id", id)
					http.Redirect(w, r, showURL, http.StatusFound)
				} else {
					fmt.Fprintf(w, "您未做任何更改")
				}
			} else {
				// 验证不通过，显示理由
				view.Render(w, view.D{
					"Article": data,
					"Errors":  errors,
				}, "articles.edit", "articles._form_field")
			}
		}
	}

}

// Delete 删除文章
func (ac *ArticleController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	article, err := article.Get(id)
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		if !policies.CanModifyArticle(article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			affectRow, err := article.Delete()
			if err != nil {
				// 应该是 SQL 报错了
				logger.LogErr(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			} else {
				if affectRow > 0 {
					url := route.Name2URL("articles.index")
					http.Redirect(w, r, url, http.StatusFound)
				} else {
					// Edge case
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprint(w, "404 文章未找到")
				}
			}
		}
	}

}

// func saveArticleToDB(title string, body string) (int64, error) {
// 	// 变量初始化
// 	var (
// 		id   int64
// 		err  error
// 		rs   sql.Result
// 		stmt *sql.Stmt
// 	)

// 	// 1. 获取一个 prepare 声明语句
// 	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
// 	// 例行的错误检测
// 	if err != nil {
// 		return 0, err
// 	}

// 	// 2. 在此函数运行结束后关闭此语句，防止占用 SQL 连接
// 	defer stmt.Close()

// 	// 3. 执行请求，传参进入绑定的内容
// 	rs, err = stmt.Exec(title, body)
// 	if err != nil {
// 		return 0, err
// 	}

// 	// 4. 插入成功的话，会返回自增 ID
// 	if id, err = rs.LastInsertId(); id > 0 {
// 		return id, nil
// 	}

// 	return 0, nil
// }
