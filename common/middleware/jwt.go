package middleware

import (
	"app/service/jwt_service"
	"github.com/gin-gonic/gin"
)

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")

		if token == "" {
			c.JSON(401, gin.H{"code": 1, "msg": "token不能为空"})
			c.Abort()
			return
		}

		user, err := (&jwt_service.JwtSvc{}).ParseToken(token)
		if err != nil {
			c.JSON(401, gin.H{"code": 1, "msg": err.Error()})
			c.Abort()
			return
		}
		c.Set("user", user.ID)
		c.Next()
	}
}
