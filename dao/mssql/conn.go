package mssql

import (
	"app/common/config"
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	var err error
	dsn := config.GetConfig().SqlUrl
	db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("mysql连接成功，准备更新表格......")
	migrate()
}
