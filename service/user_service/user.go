package user_service

import (
	"app/common/e"
	"app/dao/mssql"
	"app/models"
)

type UserSvc struct {
}

func (*UserSvc) Login(user *models.User) error {
	if err := mssql.FindOne("users", "name = ? and password = ?", e.I{user.Name, user.Password}, user); err != nil {
		return err
	}
	return nil
}

func (*UserSvc) Register(user *models.User) error {
	if err := mssql.FindOne("users", "name = ? and password = ?", e.I{user.Name, user.Password}, user); err != nil {
		return err
	}
	if user.ID != 0 {
		return e.NewError("用户已注册！", nil)
	}
	return mssql.Insert("users", user)
}
