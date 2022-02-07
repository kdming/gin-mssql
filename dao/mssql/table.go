package mssql

import (
	"app/models"
	"fmt"
)

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
