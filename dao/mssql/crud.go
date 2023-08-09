package mssql

import (
	"gorm.io/gorm"
)

type crudService struct {
	db *gorm.DB
}

// NewCrud new crudSvc
func NewCrud(taxModel bool) *crudService {
	if taxModel {
		return &crudService{
			db: db,
		}
	}
	return &crudService{db: db.Begin()}
}

// Insert 插入数据
func (s *crudService) Insert(model, data interface{}) error {
	return s.db.Model(model).Create(data).Error
}

// FindOne 查询单条数据
func (s *crudService) FindOne(model interface{}, sql string, params []interface{}, out interface{}) error {
	err := s.db.Model(model).Where(sql, params...).First(out).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return err
}

// Find 查询数据
func (s *crudService) Find(model interface{}, sql string, params []interface{}, out interface{}) error {
	return s.db.Model(model).Where(sql, params...).Find(out).Error
}

// Count 统计数据条数
func (s *crudService) Count(model interface{}, sql string, params []interface{}) (int, error) {
	var total int64
	err := s.db.Model(model).Where(sql, params...).Count(&total).Error
	return int(total), err
}

// Update 更新数据
func (s *crudService) Update(model interface{}, sql string, params []interface{}, update interface{}) error {
	return s.db.Model(model).Where(sql, params...).Updates(update).Error
}

// Delete 删除数据（非软删除）
func (s *crudService) Delete(model interface{}, sql string, params []interface{}) error {
	return s.db.Where(sql, params...).Unscoped().Delete(model).Error
}

// ExecNTable 执行原生sql语句
func (s *crudService) ExecNTable(sql string, params []interface{}, out interface{}) error {
	return s.db.Raw(sql, params...).Scan(out).Error
}

// RollBack 事物回滚
func (s *crudService) RollBack() error {
	return s.db.Rollback().Error
}

// Commit 事物提交
func (s *crudService) Commit() error {
	return s.db.Commit().Error
}

func (s *crudService) GetDB() *gorm.DB {
	return s.db
}
