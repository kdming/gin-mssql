package api

import (
	v1 "app/api/v1"
	"app/common/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(g *gin.Engine) {
	middleware.Cors(g)
	middleware.PPROF(g)
	g.Use(middleware.ErrCatch)

	api := g.Group("/api/")
	api.POST("login", v1.Login)
	api.Use(middleware.TokenAuth())
	api.POST("register", v1.Register)
}
