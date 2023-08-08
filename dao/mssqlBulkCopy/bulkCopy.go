package mssqlBulkCopy

import (
	"app/common/config"
	"database/sql"
	"fmt"
	msq "github.com/denisenkom/go-mssqldb"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// BulkCopy 批量插入
func BulkCopy(tableName string, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	// 连接数据库
	conf := config.GetConfig()
	db, err := sql.Open("sqlserver", conf.SqlUrl)
	defer db.Close()
	if err != nil {
		fmt.Println("Open connection failed:", err.Error())
		return err
	}

	// 生成columns
	columns := []string{}
	sourceColumns := []string{}
	useGormModel := false
	tp := reflect.TypeOf(data[0])
	for i := 0; i < tp.NumField(); i++ {
		name := tp.Field(i).Name
		if name == "Model" || name == "gorm.Model" {
			useGormModel = true
			continue
		}
		if name == "ID" {
			continue
		}
		tag := tp.Field(i).Tag.Get("gorm")
		reg := regexp.MustCompile(`(column:.?.*)`)
		arr := reg.FindAllString(tag, 1)
		var column string
		if len(arr) == 0 {
			column = snakeString(name)
		} else {
			column = strings.ReplaceAll(arr[0], "column:", "")
		}
		columns = append(columns, column)
		sourceColumns = append(sourceColumns, name)
	}
	if useGormModel {
		columns = append(columns, "created_at")
		columns = append(columns, "updated_at")
		columns = append(columns, "deleted_at")
	}

	// 准备bulkCopy
	txn, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := txn.Prepare(msq.CopyIn(tableName, msq.BulkOptions{}, columns...))
	if err != nil {
		return err
	}
	defer stmt.Close()

	fmt.Println("bulkCopy开始：", time.Now())
	cLen := len(columns)
	createdAt := time.Now()
	for i := 0; i < len(data); i++ {
		v := reflect.ValueOf(data[i])
		values := make([]interface{}, cLen)
		sLen := len(sourceColumns)
		for i := 0; i < sLen; i++ {
			column := sourceColumns[i]
			val := v.FieldByName(column).Interface()
			values[i] = val
		}
		if useGormModel {
			values[sLen] = createdAt
			values[sLen+1] = createdAt
			values[sLen+2] = nil
		}
		_, err = stmt.Exec(values...)
		if err != nil {
			txn.Rollback()
			return err
		}
	}

	// 提交
	if err := txn.Commit(); err != nil {
		return err
	}

	// 获取插入条数信息
	result, err := stmt.Exec()
	if err != nil {
		return err
	}
	rowCount, _ := result.RowsAffected()
	fmt.Println("bulkCopy结束，插入条数:", rowCount, time.Now())
	return nil
}

// 驼峰命名转换
func snakeString(s string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snakeCase := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snakeCase)
}
