package user_service

import (
	"app/common/app"
	"app/dao/mssql"
	"app/models"
)

type UserSvc struct {
}

// Login 用户登录
func (*UserSvc) Login(user *models.User) error {
	return mssql.NewCrud(false).FindOne(&models.User{}, "name = ? and password = ?", app.I{user.Name, user.Password}, user)
}

// Register 用户注册
func (*UserSvc) Register(user *models.User) error {
	crud := mssql.NewCrud(false)
	if err := crud.FindOne("users", "name = ? and password = ?", app.I{user.Name, user.Password}, user); err != nil {
		return err
	}
	if user.ID != 0 {
		return app.NewError("用户已注册！", nil)
	}
	return crud.Insert("users", user)
}
