package controllers

import (
	"fmt"
	"net/http"

	"goblog/app/models/user"
	"goblog/pkg/logger"
	"goblog/pkg/view"
)

type AuthController struct {
}

// Register 渲染注册界面
func (auth *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

func (auth *AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	_user := user.User{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	err := _user.Create()
	logger.LogErr(err)
	if _user.ID > 0 {
		fmt.Fprint(w, "插入成功，ID 为"+_user.GetStringID())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "创建用户失败，请联系管理员")
	}
}
