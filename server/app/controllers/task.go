package controllers

import (
	"github.com/gin-gonic/gin"
	"server/app/models"
	"server/app/services"
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

	if err := services.Task.Save(&req); err != nil {
		c.JSON(200, gin.H{"code": 500, "msg": err.Error(), "data": nil})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": nil, "data": "保存成功"})
}
