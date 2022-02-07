package mssql

import (
	"gorm.io/gorm"
)

func Insert(table string, data interface{}) error {
	return db.Table(table).Create(data).Error
}

func FindOne(table, sql string, params []interface{}, out interface{}) error {
	err := db.Table(table).Where(sql, params...).First(out).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return err
}

func Find(table, sql string, params []interface{}, out interface{}) error {
	return db.Table(table).Where(sql, params...).Find(out).Error
}

func Count(table, sql string, params []interface{}) (int, error) {
	var total int64
	err := db.Table(table).Where(sql, params...).Count(&total).Error
	return int(total), err
}

func Update(table, sql string, params []interface{}, update interface{}) error {
	return db.Table(table).Where(sql, params...).Updates(update).Error
}

func Delete(model interface{}, sql string, params []interface{}) error {
	return db.Where(sql, params...).Unscoped().Delete(model).Error
}

func Exec(table, sql string, params []interface{}, out interface{}) error {
	return db.Table(table).Raw(sql, params...).Scan(out).Error
}

func ExecNTable(sql string, params []interface{}, out interface{}) error {
	return db.Raw(sql, params...).Scan(out).Error
}
