package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/time/models"
	"evolution/backend/time/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Area struct {
	structs.Controller
	AreaSvc  *services.Area
	FieldSvc *services.Field
}

func NewArea() *Area {
	Area := &Area{}
	Area.Init()
	Area.Name = "area"
	Area.ProjectName = Area.Config.Time.System.Name
	Area.Prepare()
	Area.AreaSvc = services.NewArea(Area.Engine, Area.Cache)
	Area.FieldSvc = services.NewField(Area.Engine, Area.Cache)
	return Area
}

func (c *Area) Router(router *gin.RouterGroup) {
	area := router.Group(c.Name).Use(middles.One(c.AreaSvc, c.Name))
	area.GET("/get/:id", c.One)
	area.GET("/list/all", c.List)
	area.GET("/list/parent/:fieldId", c.ListParent)
	area.GET("/list/children/:id", c.ListChildren)
	area.GET("/list/tree/one/:fieldId", c.ListOneTree)
	area.GET("/list/tree/all", c.ListAllTree)
	area.POST("/list", c.ListWithCondition)
	area.POST("", c.Add)
	area.PUT("/:id", c.Update)
	area.DELETE("/:id", c.Delete)
}

func (c *Area) One(ctx *gin.Context) {
	area := ctx.MustGet(c.Name).(models.Area)
	resp.Success(ctx, area)
}

func (c *Area) List(ctx *gin.Context) {
	areas, err := c.AreaSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "area get error", err)
		return
	}
	resp.Success(ctx, areas)
}

func (c *Area) ListWithCondition(ctx *gin.Context) {
	var area models.Area
	if err := ctx.ShouldBindJSON(&area); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	areas, err := c.AreaSvc.ListWithCondition(&area)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "area get error", err)
		return
	}
	resp.Success(ctx, areas)
}

func (c *Area) ListAllTree(ctx *gin.Context) {
	areas, err := c.AreaSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "area get error", err)
		return
	}

	fieldMap, err := c.FieldSvc.Map()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "field Id Map get faield", err)
		return
	}

	areaTrees, err := c.AreaSvc.TransferListToTree(areas, fieldMap)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDataService, "area tree transfer error", err)
		return
	}
	resp.Success(ctx, areaTrees)
}

func (c *Area) ListParent(ctx *gin.Context) {
	fieldIdStr := ctx.Param("fieldId")
	fieldId, err := strconv.Atoi(fieldIdStr)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "fieldId params error", err)
		return
	}

	area := models.Area{
		FieldId: fieldId,
		Type:    models.AreaTypeRoot,
	}
	areas, err := c.AreaSvc.ListWithCondition(&area)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "area get error", err)
		return
	}
	resp.Success(ctx, areas)
}

func (c *Area) ListChildren(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "id params error", err)
		return
	}

	area := models.Area{
		Parent: id,
	}
	areas, err := c.AreaSvc.ListWithCondition(&area)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "area get error", err)
		return
	}
	resp.Success(ctx, areas)
}

func (c *Area) ListOneTree(ctx *gin.Context) {
	fieldIdStr := ctx.Param("fieldId")
	fieldId, err := strconv.Atoi(fieldIdStr)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "fieldId params error", err)
		return
	}

	area := models.Area{
		FieldId: fieldId,
	}
	areas, err := c.AreaSvc.ListWithCondition(&area)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "area get error", err)
		return
	}

	fieldMap, err := c.FieldSvc.Map()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "field Id Map get faield", err)
		return
	}
	areaTrees, err := c.AreaSvc.TransferListToTree(areas, fieldMap)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDataService, "area tree transfer error", err)
		return
	}

	_, ok := areaTrees[fieldId]
	if !ok {
		fieldName, exist := fieldMap[fieldId]
		if !exist {
			resp.ErrorBusiness(ctx, resp.ErrorDataService, "area tree transfer error:field id not exist", nil)
			return
		}
		areaTrees[fieldId] = models.AreaTree{
			Value:    fieldName,
			Children: make([]models.AreaNode, 0),
		}
	}
	resp.Success(ctx, areaTrees[fieldId])
}

func (c *Area) Add(ctx *gin.Context) {
	var area models.Area
	if err := ctx.ShouldBindJSON(&area); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.AreaSvc.Add(area)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "area insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Area) Update(ctx *gin.Context) {
	var area models.Area
	if err := ctx.ShouldBindJSON(&area); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.AreaSvc.Update(area.Id, area)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "area update error", err)
		return
	}

	resp.Success(ctx, struct{}{})
}

func (c *Area) Delete(ctx *gin.Context) {
	area := ctx.MustGet(c.Name).(models.Area)
	err := c.AreaSvc.Delete(area.Id, area)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "area delete error", err)
		return
	}
	resp.Success(ctx, area.Id)
}
