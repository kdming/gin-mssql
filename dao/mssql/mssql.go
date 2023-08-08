package mssql

import (
	"app/common/config"
	"app/models"
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Connect 数据库连接
func Connect() {
	var err error
	dsn := config.GetConfig().SqlUrl
	// gorm db 配置，禁用默认事务，设置日志等级
	dbConfig := &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	}
	db, err = gorm.Open(sqlserver.Open(dsn), dbConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println("数据库连接成功，准备更新表格......")
	migrate()
}

// migrate 表结构生成或者初始化
func migrate() {
	tables := []interface{}{
		&models.User{},
	}
	for i := 0; i < len(tables); i++ {
		if err := db.AutoMigrate(tables[i]); err != nil {
			panic("表格更新失败" + err.Error())
		}
	}
	fmt.Println("表格更新完毕......")
}
