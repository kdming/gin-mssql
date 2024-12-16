package v1

import (
	"app/common/app"
	"app/models"
	"app/service/jwt_service"
	"app/service/user_service"
	"github.com/gin-gonic/gin"
)

// Login 用户登录
func Login(c *gin.Context) {
	user := &models.User{}
	if err := c.ShouldBind(user); err != nil {
		app.Err("参数绑定失败", err)
	}

	userSvc := &user_service.UserSvc{}
	if err := userSvc.Login(user); err != nil {
		app.Err("登录失败", err)
	}
	if user.ID == 0 {
		app.Err("用户不存在", nil)
	}

	token, err := (&jwt_service.JwtSvc{}).MakeToken(user)
	if err != nil {
		app.Err("token生成失败", err)
	}
	app.Ok(c, "登录成功", app.M{"token": token, "name": user.Name})
}

// Register 用户注册
func Register(c *gin.Context) {
	user := &models.User{}
	if err := c.ShouldBind(user); err != nil {
		app.Err("参数绑定失败", err)
	}
	if user.Name == "" || user.Password == "" {
		app.Err("请输入完整的账号与密码", nil)
	}

	if err := (&user_service.UserSvc{}).Register(user); err != nil {
		app.Err("注册失败", err)
	}

	token, err := (&jwt_service.JwtSvc{}).MakeToken(user)
	if err != nil {
		app.Err("token生成失败", err)
	}
	app.Ok(c, "注册成功", app.M{"token": token, "name": user.Name})
}
