package services

import (
	"context"
	"gorm.io/gorm/clause"
	"server/app/models"
	"server/app/types"
	"server/database"
	"strconv"
	"time"
)

var (
	defaultExpireTime = 10 * time.Minute
)

type TaskServiceImpl struct{}

func (TaskServiceImpl) Save(req *models.Task) error {

	if req.Result != "" {
		req.Status = types.TaskStatusFinished
	}

	return database.GetMySQL().Save(req).Error
}

func (TaskServiceImpl) List(req *models.TaskReq) (models.TaskListResp, error) {

	var resp models.TaskListResp

	db := database.GetMySQL().Model(&models.Task{})

	filters := []QueryOption{
		WithID32(req.ID),
		WithTaskStatus(req.Status),
	}

	ApplyFilters(db, filters...)

	if err := db.Count(&resp.Total).Error; err != nil {
		return resp, err
	}

	if err := db.Preload(clause.Associations).Find(&resp.Records).Error; err != nil {
		return resp, err
	}

	return resp, nil
}

func (TaskServiceImpl) GetTask() (models.TaskResp, error) {

	var tasks []models.TaskResp

	if err := database.GetMySQL().Model(&models.Task{}).Where("status = ?", types.TaskStatusUnfinished).Limit(100).Find(&tasks).Error; err != nil {
		return models.TaskResp{}, err
	}

	for _, task := range tasks {

		taskID := strconv.Itoa(task.ID)

		res, _ := database.GetRedis().Incr(context.Background(), taskID).Result()
		if res > 1 {
			continue
		}

		if _, err := database.GetRedis().Expire(context.Background(), taskID, defaultExpireTime).Result(); err != nil {
			return models.TaskResp{}, err
		}

		return task, nil
	}

	return models.TaskResp{}, nil
}

func releaseTaskLock(taskID string) error {
	_, err := database.GetRedis().Del(context.Background(), taskID).Result()
	return err
}
