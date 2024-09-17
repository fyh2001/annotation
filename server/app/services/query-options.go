package services

import (
	"gorm.io/gorm"
	"server/app/types"
)

type QueryOption func(*gorm.DB)

func ApplyFilters(db *gorm.DB, opts ...QueryOption) {
	for _, opt := range opts {
		opt(db) // 调用每个过滤器，应用到db上
	}
}

func WithID32(id int) QueryOption {
	return func(db *gorm.DB) {
		if id != 0 {
			db.Where("id = ?", id)
		}
	}
}

func WithTaskStatus(status types.TaskStatusType) QueryOption {
	return func(db *gorm.DB) {
		if status != 0 {
			db.Where("status = ?", status)
		}
	}
}
