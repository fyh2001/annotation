package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm/clause"
	"log"
	"os"
	"server/app/models"
	"server/app/types"
	"server/database"
	"strconv"
	"time"
)

var (
	defaultExpireTime = 10 * time.Second
)

type TaskServiceImpl struct{}

func (TaskServiceImpl) Save(req *models.Task) error {

	if req.Result != "" {
		req.Status = types.TaskStatusFinished
	}

	if err := database.GetMySQL().Save(req).Error; err != nil {
		return err
	}

	if req.Status == types.TaskStatusFinished {

		type result struct {
			Classification int     `json:"classification"`
			X              float64 `json:"x"`
			Y              float64 `json:"y"`
			Width          float64 `json:"width"`
			Height         float64 `json:"height"`
		}

		var resultData []result

		json.Unmarshal([]byte(req.Result), &resultData)

		idStr := strconv.Itoa(req.ID)

		file, err := os.Create("./static/file/" + idStr + ".txt")
		if err != nil {
			return err
		}
		defer file.Close()

		for _, v := range resultData {
			_, err = file.WriteString(fmt.Sprintf("%d %.8f %.8f %.8f %.8f\n", v.Classification, v.X, v.Y, v.Width, v.Height))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
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

		expiredAt := time.Now().Add(defaultExpireTime)
		task.ExpiredAt = expiredAt

		if _, err := database.GetRedis().ExpireAt(context.Background(), taskID, expiredAt).Result(); err != nil {
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
