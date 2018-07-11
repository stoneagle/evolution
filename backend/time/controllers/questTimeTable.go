package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type QuestTimeTable struct {
	structs.Controller
	QuestTimeTableSvc *services.QuestTimeTable
}

func NewQuestTimeTable() *QuestTimeTable {
	QuestTimeTable := &QuestTimeTable{}
	QuestTimeTable.Init()
	QuestTimeTable.ProjectName = QuestTimeTable.Config.Time.System.Name
	QuestTimeTable.Name = "quest-timetable"
	QuestTimeTable.Prepare()
	QuestTimeTable.QuestTimeTableSvc = services.NewQuestTimeTable(QuestTimeTable.Engine, QuestTimeTable.Cache)
	return QuestTimeTable
}

func (c *QuestTimeTable) Router(router *gin.RouterGroup) {
	questTimeTable := router.Group(c.Name).Use(middles.One(c.QuestTimeTableSvc, c.Name))
	questTimeTable.GET("/get/:id", c.One)
	questTimeTable.GET("/list", c.List)
	questTimeTable.POST("", c.Add)
	questTimeTable.POST("/list", c.ListByCondition)
	questTimeTable.PUT("/:id", c.Update)
	questTimeTable.DELETE("/:id", c.Delete)
}

func (c *QuestTimeTable) One(ctx *gin.Context) {
	questTimeTable := ctx.MustGet(c.Name).(models.QuestTimeTable)
	resp.Success(ctx, questTimeTable)
}

func (c *QuestTimeTable) List(ctx *gin.Context) {
	questTimeTables, err := c.QuestTimeTableSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questTimeTable get error", err)
		return
	}
	resp.Success(ctx, questTimeTables)
}

func (c *QuestTimeTable) ListByCondition(ctx *gin.Context) {
	var questTimeTable models.QuestTimeTable
	if err := ctx.ShouldBindJSON(&questTimeTable); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	questTimeTables, err := c.QuestTimeTableSvc.ListWithCondition(&questTimeTable)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questTimeTable get error", err)
		return
	}
	resp.Success(ctx, questTimeTables)
}

func (c *QuestTimeTable) Add(ctx *gin.Context) {
	var questTimeTable models.QuestTimeTable
	if err := ctx.ShouldBindJSON(&questTimeTable); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestTimeTableSvc.Add(questTimeTable)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questTimeTable insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *QuestTimeTable) Update(ctx *gin.Context) {
	var questTimeTable models.QuestTimeTable
	if err := ctx.ShouldBindJSON(&questTimeTable); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestTimeTableSvc.Update(questTimeTable.Id, questTimeTable)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questTimeTable update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *QuestTimeTable) Delete(ctx *gin.Context) {
	questTimeTable := ctx.MustGet(c.Name).(models.QuestTimeTable)
	err := c.QuestTimeTableSvc.Delete(questTimeTable.Id, questTimeTable)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questTimeTable delete error", err)
		return
	}
	resp.Success(ctx, questTimeTable.Id)
}
