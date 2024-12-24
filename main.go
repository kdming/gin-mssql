package main

import (
	"app/api"
	"app/common/config"
	"app/dao/mssql"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	mssql.Connect()
	g := gin.New()
	api.InitRouter(g)
	appPort := config.GetConfig().AppPort
	if appPort == "" {
		appPort = ":8080"
	}
	if err := g.Run(fmt.Sprintf(":%v", appPort)); err != nil {
		panic(err)
	}
}
