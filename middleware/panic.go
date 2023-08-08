package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// ErrCatch 异常捕获
func ErrCatch(g *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			g.JSON(200, gin.H{"code": 1, "data": nil, "msg": fmt.Sprintf("%v", err)})
			g.Abort()
		}
	}()
	g.Next()
}
