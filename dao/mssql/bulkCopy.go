package mssql

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
	fmt.Println(len(columns), useGormModel, "columns len")

	// 准备bulkCopy
	txn, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := txn.Prepare(msq.CopyIn(tableName, msq.BulkOptions{}, columns...))
	if err != nil {
		return err
	}

	fmt.Println("插入开始", time.Now())
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
			fmt.Println(err.Error())
			return err
		}
	}

	// 执行
	result, err := stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 关闭
	err = stmt.Close()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// 提交
	err = txn.Commit()
	if err != nil {
		return err
	}
	rowCount, _ := result.RowsAffected()
	fmt.Printf("%d row copied\n", rowCount)
	fmt.Println("插入结束", time.Now())
	return nil
}

func snakeString(s string) string {
	ssr := []byte{}
	for i := 0; i < len(s); i++ {
		a := s[i]
		ii := i + 1
		c := byte('n')
		if ii <= len(s)-1 {
			b := s[ii]
			if (a >= 65 && a <= 90) && (b >= 97 && b <= 122) && b != '_' && i > 0 && ssr[i-1] != '_' {
				c = '_'
			} else {
				c = 'n'
			}
			if c != 'n' {
				ssr = append(ssr, []byte{c, a}...)
			} else {
				ssr = append(ssr, a)
			}
		} else {
			ssr = append(ssr, a)
		}
	}
	return strings.ToLower(string(ssr))
}
