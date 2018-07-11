package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type QuestEntity struct {
	structs.Controller
	QuestEntitySvc *services.QuestEntity
}

func NewQuestEntity() *QuestEntity {
	QuestEntity := &QuestEntity{}
	QuestEntity.Init()
	QuestEntity.ProjectName = QuestEntity.Config.Time.System.Name
	QuestEntity.Name = "quest-entity"
	QuestEntity.Prepare()
	QuestEntity.QuestEntitySvc = services.NewQuestEntity(QuestEntity.Engine, QuestEntity.Cache)
	return QuestEntity
}

func (c *QuestEntity) Router(router *gin.RouterGroup) {
	questEntity := router.Group(c.Name).Use(middles.One(c.QuestEntitySvc, c.Name))
	questEntity.GET("/get/:id", c.One)
	questEntity.GET("/list", c.List)
	questEntity.POST("", c.Add)
	questEntity.POST("/list", c.ListByCondition)
	questEntity.PUT("/:id", c.Update)
	questEntity.DELETE("/:id", c.Delete)
}

func (c *QuestEntity) One(ctx *gin.Context) {
	questEntity := ctx.MustGet(c.Name).(models.QuestEntity)
	resp.Success(ctx, questEntity)
}

func (c *QuestEntity) List(ctx *gin.Context) {
	questEntitys, err := c.QuestEntitySvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questEntity get error", err)
		return
	}
	resp.Success(ctx, questEntitys)
}

func (c *QuestEntity) ListByCondition(ctx *gin.Context) {
	var questEntity models.QuestEntity
	if err := ctx.ShouldBindJSON(&questEntity); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	questEntitys, err := c.QuestEntitySvc.ListWithCondition(&questEntity)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questEntity get error", err)
		return
	}
	resp.Success(ctx, questEntitys)
}

func (c *QuestEntity) Add(ctx *gin.Context) {
	var questEntity models.QuestEntity
	if err := ctx.ShouldBindJSON(&questEntity); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestEntitySvc.Add(questEntity)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questEntity insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *QuestEntity) Update(ctx *gin.Context) {
	var questEntity models.QuestEntity
	if err := ctx.ShouldBindJSON(&questEntity); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.QuestEntitySvc.Update(questEntity.Id, questEntity)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questEntity update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *QuestEntity) Delete(ctx *gin.Context) {
	questEntity := ctx.MustGet(c.Name).(models.QuestEntity)
	err := c.QuestEntitySvc.Delete(questEntity.Id, questEntity)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "questEntity delete error", err)
		return
	}
	resp.Success(ctx, questEntity.Id)
}
