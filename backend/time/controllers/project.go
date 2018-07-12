package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"

	"github.com/gin-gonic/gin"
)

type Project struct {
	structs.Controller
	ProjectSvc   *services.Project
	QuestSvc     *services.Quest
	QuestTeamSvc *services.QuestTeam
}

func NewProject() *Project {
	Project := &Project{}
	Project.Init()
	Project.ProjectName = Project.Config.Time.System.Name
	Project.Name = "project"
	Project.Prepare()
	Project.ProjectSvc = services.NewProject(Project.Engine, Project.Cache)
	Project.QuestSvc = services.NewQuest(Project.Engine, Project.Cache)
	Project.QuestTeamSvc = services.NewQuestTeam(Project.Engine, Project.Cache)
	return Project
}

func (c *Project) Router(router *gin.RouterGroup) {
	project := router.Group(c.Name).Use(middles.One(c.ProjectSvc, c.Name))
	project.GET("/get/:id", c.One)
	project.GET("/list", c.List)
	project.GET("/list/syncfusion/", c.ListSyncfusionFormat)
	project.POST("", c.Add)
	project.POST("/list", c.ListByCondition)
	project.PUT("/:id", c.Update)
	project.DELETE("/:id", c.Delete)
}

func (c *Project) One(ctx *gin.Context) {
	project := ctx.MustGet(c.Name).(models.Project)
	resp.Success(ctx, project)
}

func (c *Project) List(ctx *gin.Context) {
	projects, err := c.ProjectSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "project get error", err)
		return
	}
	resp.Success(ctx, projects)
}

func (c *Project) ListSyncfusionFormat(ctx *gin.Context) {
	user := ctx.MustGet(middles.UserKey).(middles.UserInfo)
	questTeam := models.QuestTeam{}
	questTeam.UserId = user.Id
	questTeams, err := c.QuestTeamSvc.ListWithCondition(&questTeam)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "quest team get error", err)
		return
	}

	Ids := make([]int, 0)
	for _, one := range questTeams {
		Ids = append(Ids, one.QuestId)
	}
	quest := models.Quest{}
	quest.Ids = Ids
	quests, err := c.QuestSvc.ListWithCondition(&quest)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "quest get error", err)
		return
	}

	res := make([]models.SyncfusionGantt, 0)
	for _, one := range quests {
		gantt := models.SyncfusionGantt{}
		gantt.Id = one.Id
		gantt.Parent = 0
		gantt.Name = one.Name
		gantt.StartDate = one.EndDate
		gantt.EndDate = one.EndDate
		gantt.Progress = 0
		gantt.Duration = 0
		gantt.Expanded = false
		child := make([]models.SyncfusionGantt, 0)
		// child = append(child, gantt)
		gantt.Children = child
		res = append(res, gantt)
	}
	resp.CustomSuccess(ctx, res)
}

func (c *Project) ListByCondition(ctx *gin.Context) {
	var project models.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	projects, err := c.ProjectSvc.ListWithCondition(&project)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "project get error", err)
		return
	}
	resp.Success(ctx, projects)
}

func (c *Project) Add(ctx *gin.Context) {
	var project models.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.ProjectSvc.Add(project)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "project insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Project) Update(ctx *gin.Context) {
	var project models.Project
	if err := ctx.ShouldBindJSON(&project); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.ProjectSvc.Update(project.Id, project)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "project update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Project) Delete(ctx *gin.Context) {
	project := ctx.MustGet(c.Name).(models.Project)
	err := c.ProjectSvc.Delete(project.Id, project)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "project delete error", err)
		return
	}
	resp.Success(ctx, project.Id)
}
