package models

import (
	"server/app/types"
	"time"
)

type Task struct {
	ID        int                  `json:"id" gorm:"primaryKey"`
	Url       string               `json:"url"`
	Status    types.TaskStatusType `json:"status"`
	Result    string               `json:"result"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
}

type TaskReq struct {
	ID     int                  `json:"id"`
	Url    string               `json:"url"`
	Status types.TaskStatusType `json:"status"`
	Result string               `json:"result"`
}

type TaskResp struct {
	ID        int                  `json:"id"`
	Url       string               `json:"url"`
	Status    types.TaskStatusType `json:"status"`
	Result    string               `json:"result"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
	ExpiredAt time.Time            `json:"expiredAt" gorm:"-"`
}

type TaskListResp struct {
	Total   int64      `json:"total"`
	Records []TaskResp `json:"records"`
}

func (TaskResp) TableName() string {
	return "task"
}
