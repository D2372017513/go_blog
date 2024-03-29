package controllers

import (
	"fmt"
	"net/http"

	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"goblog/pkg/logger"
	"goblog/pkg/view"
)

type AuthController struct {
}

// Register 渲染注册界面
func (auth *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

func (a *AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	_user := user.User{
		Name:            r.FormValue("name"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("password_confirm"),
	}

	errs := requests.ValidateRegistrationForm(_user)

	if len(errs) > 0 {
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.register")
	} else {
		err := _user.Create()
		logger.LogErr(err)
		if _user.ID > 0 {
			// 登录用户并跳转到首页
			flash.Success("恭喜您注册成功！")
			auth.Login(_user)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建用户失败，请联系管理员")
		}
	}
}

// Login 渲染登录界面
func (auth *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.login")
}

// DoLogin 处理登录逻辑
func (a *AuthController) DoLogin(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	passwd := r.PostFormValue("password")

	if err := auth.Attempt(email, passwd); err == nil {
		// 登录成功
		flash.Success("欢迎回来！")
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		view.RenderSimple(w, view.D{
			"Error":    err.Error(),
			"Email":    email,
			"password": passwd,
		}, "auth.login")
	}
}

// Logout 退出登录
func (*AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	auth.Logout()
	flash.Success("您已退出登录")
	http.Redirect(w, r, "/", http.StatusFound)
}
