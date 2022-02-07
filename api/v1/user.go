package v1

import (
	"app/common/e"
	"app/models"
	"app/service/jwt_service"
	"app/service/user_service"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	user := &models.User{}
	if err := c.ShouldBind(user); err != nil {
		e.Err("参数绑定失败", err)
	}

	userSvc := &user_service.UserSvc{}
	if err := userSvc.Login(user); err != nil {
		e.Err("登录失败", err)
	}
	if user.ID == 0 {
		e.Err("用户不存在", nil)
	}

	token, err := (&jwt_service.JwtSvc{}).MakeToken(user)
	if err != nil {
		e.Err("token生成失败", err)
	}

	e.Ok(c, "登录成功", e.M{"token": token, "name": user.Name})
}

func Register(c *gin.Context) {
	user := &models.User{}
	if err := c.ShouldBind(user); err != nil {
		e.Err("参数绑定失败", err)
	}
	if user.Name == "" || user.Password == "" {
		e.Err("请输入完整的账号与密码", nil)
	}

	if err := (&user_service.UserSvc{}).Register(user); err != nil {
		e.Err("注册失败", err)
	}

	token, err := (&jwt_service.JwtSvc{}).MakeToken(user)
	if err != nil {
		e.Err("token生成失败", err)
	}

	e.Ok(c, "注册成功", e.M{"token": token, "name": user.Name})
}
