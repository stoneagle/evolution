package controllers

import (
	"evolution/backend/common/middles"

	"github.com/gin-gonic/gin"
)

type QuestTimeTable struct {
	BaseController
}

func NewQuestTimeTable() *QuestTimeTable {
	QuestTimeTable := &QuestTimeTable{}
	return QuestTimeTable
}

func (c *QuestTimeTable) Router(router *gin.RouterGroup) {
	questTimeTable := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	questTimeTable.GET("/get/:id", c.One)
	questTimeTable.POST("", c.Create)
	questTimeTable.POST("/list", c.List)
	questTimeTable.PUT("/:id", c.Update)
	questTimeTable.DELETE("/:id", c.Delete)
}
