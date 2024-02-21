package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"

	"goblog/app/models"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"goblog/types"

	"gorm.io/gorm"
)

type ArticleController struct {
}

func (ac *ArticleController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	rs, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
		view.Render(w, rs, "articles.show")
	}
}

func (ac *ArticleController) Index(w http.ResponseWriter, r *http.Request) {
	articles, err := article.GetAll()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Fprintf(w, "当前没有任何文章可供浏览")
		} else {
			logger.LogErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "500 服务器内部错误")
		}
	} else {
		view.Render(w, articles, "articles.index")
	}
}

func (ac *ArticleController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, article.ArticlesData{}, "articles.create", "articles._form_field")
}

func (ac *ArticleController) Store(w http.ResponseWriter, r *http.Request) {
	// Form：存储了 post、put 和 get 参数，在使用之前需要调用 ParseForm 方法。PostForm：存储了 post、put 参数，在使用之前需要调用 ParseForm 方法。
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("解析出错！")
		return
	}

	title, body := r.PostFormValue("title"), r.PostFormValue("body")
	errors := validateArticleFormData(title, body)
	data := article.ArticlesData{
		Title:  title,
		Body:   body,
		Errors: errors,
	}

	if len(errors) != 0 {
		view.Render(w, data, "articles.create", "articles._form_field")
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
		if err == gorm.ErrRecordNotFound {
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
		view.Render(w, rs, "articles.edit", "articles._form_field")
	}
}

func (ac *ArticleController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
		updateURL := route.Name2URL("articles.edit", "id", id)
		errors := validateArticleFormData(title, body)
		data := article.ArticlesData{
			BaseModel: models.BaseModel{ID: int64(types.StringToInt64(id))},
			Title:     title,
			Body:      body,
			URL:       updateURL,
			Errors:    errors,
		}

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
			view.Render(w, data, "articles.edit", "articles._form_field")
		}
	}

}

// Delete 删除文章
func (ac *ArticleController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	article, err := article.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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

// 检查提交的内容是否有效
func validateArticleFormData(title, body string) map[string]string {
	errors := make(map[string]string)
	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
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