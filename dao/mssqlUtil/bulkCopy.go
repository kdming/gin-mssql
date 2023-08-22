package mssqlUtil

import (
	"app/common/app"
	"app/common/config"
	"database/sql"
	"errors"
	"fmt"
	mssql "github.com/microsoft/go-mssqldb"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// BulkCopy 批量插入数据
func BulkCopy(model interface{}, data interface{}) error {
	// 检查数据是否是slice类型
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		return app.NewError("data not slice!")
	}

	// 判断model是否为指针类型
	if reflect.ValueOf(model).Kind() != reflect.Ptr {
		return fmt.Errorf("model type not ptr; is %T", model)
	}

	// 判断model是否是struct类型
	vt := reflect.ValueOf(model).Elem() // dereference the pointer
	if vt.Kind() != reflect.Struct {
		return fmt.Errorf("not struct; is %T", model)
	}

	// 调用TableName方法获取表名
	method := reflect.ValueOf(model).MethodByName("TableName")
	methodRes := method.Call(nil)
	if len(methodRes) == 0 {
		return errors.New("not found tableName")
	}
	tableName := methodRes[0].String()

	// 获取要插入的列头
	useGormModel := false       // 是否使用了gorm.Model
	modelTp := vt.Type()        // 源数据类型
	columns := []string{}       // gorm column 标签的值
	sourceColumns := []string{} // 原始字段名称
	for i := 0; i < modelTp.NumField(); i++ {
		name := modelTp.Field(i).Name
		if name == "Model" || name == "gorm.Model" {
			useGormModel = true
			continue
		}
		if name == "ID" {
			continue
		}
		tag := modelTp.Field(i).Tag.Get("gorm")
		if tag == "-" {
			continue
		}
		reg := regexp.MustCompile(`(column:.?.*)`)
		arr := reg.FindAllString(tag, 1)
		// 获取列名称，如果未发现column标签则使用驼峰命名
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
	fmt.Println(tableName, "列数量:", len(columns), useGormModel)

	// 连接数据库
	conf := config.GetConfig()
	db, err := sql.Open("mssql", conf.SqlUrl)
	defer db.Close()
	if err != nil {
		fmt.Println("Open connection failed:", err.Error())
		return err
	}

	// 开启事务
	txn, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := txn.Prepare(mssql.CopyIn(tableName, mssql.BulkOptions{}, columns...))
	if err != nil {
		return err
	}
	defer stmt.Close()

	start := time.Now()
	cLen := len(columns)
	createdAt := time.Now()
	arr := reflect.ValueOf(data)
	for i := 0; i < arr.Len(); i++ {
		item := arr.Index(i)
		values := make([]interface{}, cLen)
		sLen := len(sourceColumns)
		for i := 0; i < sLen; i++ {
			column := sourceColumns[i]
			val := item.FieldByName(column).Interface()
			values[i] = val
		}
		if useGormModel {
			values[sLen] = createdAt
			values[sLen+1] = createdAt
			values[sLen+2] = nil
		}
		_, err = stmt.Exec(values...)
		if err != nil {
			return txn.Rollback()
		}
	}

	// 获取插入条数
	result, err := stmt.Exec()
	if err != nil {
		return err
	}
	rowCount, _ := result.RowsAffected()

	// 提交事务
	if err := txn.Commit(); err != nil {
		return err
	}
	fmt.Println(tableName, "表插入结束, 耗时:", time.Now().Sub(start).Seconds(), "秒, 共插入条数:", rowCount)
	return nil
}

// 驼峰命名转换
func snakeString(s string) string {
	re := regexp.MustCompile("([a-z0-9A-Z])([A-Z])")
	snakeCase := re.ReplaceAllString(s, "${1}_$2")
	return strings.ToLower(snakeCase)
}
