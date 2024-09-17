package services

import (
	"gorm.io/gorm/clause"
	"server/app/models"
	"server/app/types"
	"server/database"
)

type TaskServiceImpl struct{}

func (TaskServiceImpl) save(req *models.Task) error {

	if req.Result != "" {
		req.Status = types.TaskStatusFinished
	}

	return database.GetMySQL().Save(req).Error
}

func (TaskServiceImpl) list(req *models.TaskReq) (models.TaskListResp, error) {

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

func (TaskServiceImpl) getTask() (models.TaskResp, error) {

}
