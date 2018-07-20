package controllers

import (
	"evolution/backend/common/middles"

	"github.com/gin-gonic/gin"
)

type QuestTeam struct {
	BaseController
}

func NewQuestTeam() *QuestTeam {
	QuestTeam := &QuestTeam{}
	QuestTeam.Resource = ResourceQuestTeam
	return QuestTeam
}

func (c *QuestTeam) Router(router *gin.RouterGroup) {
	questTeam := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	questTeam.GET("/get/:id", c.One)
	questTeam.GET("/list", c.List)
	questTeam.POST("", c.Create)
	questTeam.POST("/list", c.ListWithCondition)
	questTeam.PUT("/:id", c.Update)
	questTeam.DELETE("/:id", c.Delete)
}
