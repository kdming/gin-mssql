package api

import (
	"app/api/user_api"
	"app/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(g *gin.Engine) {
	middleware.Cors(g)
	middleware.PPROF(g)
	g.Use(middleware.ErrCatch)

	api := g.Group("/api/")
	api.POST("login", user_api.Login)
	api.Use(middleware.TokenAuth())
	api.POST("register", user_api.Register)
}
