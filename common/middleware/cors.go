package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

// 跨域
func Cors(g *gin.Engine) {
	g.Use(cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "OPTIONS", "PUT", "POST", "DELETE"},
		AllowedHeaders:     []string{"*"},
		AllowCredentials:   true,
		OptionsPassthrough: false,
		Debug:              false,
	}))
	g.OPTIONS("*options_support", func(context *gin.Context) {
		context.AbortWithStatus(http.StatusNoContent)
		return
	})
}
