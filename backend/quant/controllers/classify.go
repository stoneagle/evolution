package controllers

import (
	"evolution/backend/common/resp"
	"evolution/backend/quant/models"
	"evolution/backend/quant/rpc/engine"
	"evolution/backend/quant/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Classify struct {
	Base
	ClassifySvc *services.Classify
}

func NewClassify() *Classify {
	Classify := &Classify{}
	Classify.Prepare()
	Classify.ClassifySvc = services.NewClassify(Classify.Engine, Classify.Cache)
	return Classify
}

func (c *Classify) Router(router *gin.RouterGroup) {
	classify := router.Group("classify")
	classify.GET("/get/:id", initClassify(c.ClassifySvc), c.One)
	classify.GET("/list", c.List)
	classify.POST("/list", c.ListByAssetSource)
	classify.POST("/sync", c.Sync)
	classify.DELETE("/:id", initClassify(c.ClassifySvc), c.Delete)
}

func initClassify(svc *services.Classify) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			resp.ErrorBusiness(ctx, resp.ErrorParams, "id params error", err)
		}

		classify, err := svc.One(id)
		if err != nil {
			resp.ErrorBusiness(ctx, resp.ErrorMysql, "get classify error", err)
		}

		ctx.Set("classify", classify)
		return
	}
}

func (c *Classify) One(ctx *gin.Context) {
	classify := ctx.MustGet("classify").(models.Classify)
	resp.Success(ctx, classify)
}

func (c *Classify) List(ctx *gin.Context) {
	classifies, err := c.ClassifySvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "classify get error", err)
		return
	}

	resp.Success(ctx, classifies)
}

func (c *Classify) ListByAssetSource(ctx *gin.Context) {
	var classify models.Classify
	if err := ctx.ShouldBindJSON(&classify); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	atype, err := engine.AssetTypeFromString(classify.AssetString)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorEngine, "asset type illegal", err)
		return
	}

	classifies, err := c.ClassifySvc.ListByAssetSource(atype, classify.Type, classify.Main, classify.Sub)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "classify get error", err)
		return
	}
	resp.Success(ctx, classifies)
}

func (c *Classify) Sync(ctx *gin.Context) {
	var classify models.Classify
	if err := ctx.ShouldBindJSON(&classify); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	atype, err := engine.AssetTypeFromString(classify.AssetString)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorEngine, "asset type illegal", err)
		return
	}

	classifySlice, err := c.Rpc.GetClassify(atype, classify.Type, classify.Main, classify.Sub)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorEngine, "classify get error", err)
		return
	}

	classifyBatch := []models.Classify{}
	for _, c := range classifySlice {
		one := models.Classify{}
		one.Name = c.Name
		one.Tag = c.Tag
		one.AssetType.Asset = atype
		one.AssetType.Type = classify.Type
		one.Source.Main = classify.Main
		one.Source.Sub = classify.Sub
		classifyBatch = append(classifyBatch, one)
	}
	err = c.ClassifySvc.BatchSave(classifyBatch)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorEngine, "classify save error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Classify) Delete(ctx *gin.Context) {
	classify := ctx.MustGet("classify").(models.Classify)
	err := c.ClassifySvc.Delete(classify.Id)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "classify delete error", err)
		return
	}
	resp.Success(ctx, classify.Id)
}
