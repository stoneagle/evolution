package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

	"evolution/backend/time/models"

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
	questTimeTable.GET("/list", c.List)
	questTimeTable.POST("", c.Create)
	questTimeTable.POST("/list", c.ListByCondition)
	questTimeTable.PUT("/:id", c.Update)
	questTimeTable.DELETE("/:id", c.Delete)
}

func (c *QuestTimeTable) ListByCondition(ctx *gin.Context) {
	var questTimeTable models.QuestTimeTable
	if err := ctx.ShouldBindJSON(&questTimeTable); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	questTimeTables, err := c.QuestTimeTableSvc.ListWithCondition(&questTimeTable)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "questTimeTable get error", err)
		return
	}
	resp.Success(ctx, questTimeTables)
}
