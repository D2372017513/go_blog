package controllers

import (
	"goblog/pkg/view"
	"net/http"
)

type AuthController struct {
}

// Register 渲染注册界面
func (auth *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

func (auth *AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
}
