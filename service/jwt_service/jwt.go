package jwt_service

import (
	"app/common/config"
	"app/models"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtSvc struct {
}

func (*JwtSvc) MakeToken(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["role"] = user.Role
	claims["userId"] = user.ID
	claims["expDate"] = time.Now().Add(time.Hour * 20).Format("2006-01-02 15:04:05")
	tk, err := token.SignedString([]byte(config.GetConfig().JwtKey))
	if err != nil {
		return "", err
	}
	return tk, nil
}

func (*JwtSvc) ParseToken(token string) (*models.User, error) {
	tk, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().JwtKey), nil
	})
	if err != nil {
		return nil, errors.New("token秘钥加载失败" + err.Error())
	}

	if claims, ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
		userId := claims["userId"]
		expDate := claims["expDate"]
		role := claims["role"]
		timeNow := time.Now().Format("2006-01-02 15:04:05")
		if timeNow > expDate.(string) {
			return nil, errors.New("token已过期")
		}
		if userId == "" {
			return nil, errors.New("token解密，id为空")
		}
		user := &models.User{}
		user.Role = role.(int)
		user.ID = uint(userId.(float64))
		return user, nil
	} else {
		return nil, errors.New("token解析失败")
	}
}
