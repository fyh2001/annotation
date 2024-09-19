package controllers

import (
	"github.com/gin-gonic/gin"
	"server/app/models"
	"server/app/services"
	"server/app/types"
	"server/database"
	"strconv"
)

type TaskController struct{}

func (TaskController) GetTask(c *gin.Context) {

	resp, err := services.Task.GetTask()
	if err != nil {
		c.JSON(200, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": nil, "data": resp})
}

func (TaskController) Submit(c *gin.Context) {

	var req models.Task

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	if err := services.Task.Save(&req); err != nil {
		c.JSON(200, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": nil, "data": "保存成功"})
}

func (TaskController) Insert(c *gin.Context) {

	type insertReq struct {
		Start  int    `json:"start"`
		End    int    `json:"end"`
		Suffix string `json:"suffix"`
	}

	var req insertReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	var tasks []*models.Task

	for i := req.Start; i <= req.End; i++ {
		task := &models.Task{
			ID:     i,
			Url:    "images/" + strconv.Itoa(i) + "." + req.Suffix,
			Status: types.TaskStatusUnfinished,
		}

		tasks = append(tasks, task)
	}

	if err := database.GetMySQL().Create(&tasks).Error; err != nil {
		c.JSON(200, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": nil, "data": "操作成功"})
}
