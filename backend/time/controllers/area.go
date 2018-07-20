package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"

	"evolution/backend/time/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Area struct {
	BaseController
}

func NewArea() *Area {
	Area := &Area{}
	Area.Resource = ResourceArea
	return Area
}

func (c *Area) Router(router *gin.RouterGroup) {
	area := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	area.GET("/get/:id", c.One)
	area.GET("/list/all", c.List)
	area.GET("/list/parent/:fieldId", c.ListParent)
	area.GET("/list/children/:id", c.ListChildren)
	area.GET("/list/tree/one/:fieldId", c.ListOneTree)
	area.GET("/list/tree/all", c.ListAllTree)
	area.POST("/list", c.ListWithCondition)
	area.POST("", c.Create)
	area.PUT("/:id", c.Update)
	area.DELETE("/:id", c.Delete)
}

func (c *Area) ListAllTree(ctx *gin.Context) {
	areasPtr, err := c.AreaSvc.List(&models.Area{})
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "area get error", err)
		return
	}
	areas := *(areasPtr.(*[]models.Area))

	fieldMap, err := c.FieldSvc.Map()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "field Id Map get faield", err)
		return
	}

	areaTrees, err := c.AreaSvc.TransferListToTree(areas, fieldMap)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDataTransfer, "area tree transfer error", err)
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
	areas, err := c.AreaSvc.List(&area)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "area get error", err)
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
	areas, err := c.AreaSvc.List(&area)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "area get error", err)
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
	areasPtr, err := c.AreaSvc.List(&area)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "area get error", err)
		return
	}
	areas := *(areasPtr.(*[]models.Area))

	fieldMap, err := c.FieldSvc.Map()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "field Id Map get faield", err)
		return
	}
	areaTrees, err := c.AreaSvc.TransferListToTree(areas, fieldMap)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDataTransfer, "area tree transfer error", err)
		return
	}

	_, ok := areaTrees[fieldId]
	if !ok {
		fieldName, exist := fieldMap[fieldId]
		if !exist {
			resp.ErrorBusiness(ctx, resp.ErrorDataTransfer, "area tree transfer error:field id not exist", nil)
			return
		}
		areaTrees[fieldId] = models.AreaTree{
			Value:    fieldName,
			Children: make([]models.AreaNode, 0),
		}
	}
	resp.Success(ctx, areaTrees[fieldId])
}
