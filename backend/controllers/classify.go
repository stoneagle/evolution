package controllers

import (
	"quant/backend/common"
	"quant/backend/models"
	"quant/backend/rpc/engine"
	"quant/backend/services"
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
			common.ResponseErrorBusiness(ctx, common.ErrorParams, "id params error", err)
		}

		classify, err := svc.One(id)
		if err != nil {
			common.ResponseErrorBusiness(ctx, common.ErrorMysql, "get classify error", err)
		}

		ctx.Set("classify", classify)
		return
	}
}

func (c *Classify) One(ctx *gin.Context) {
	classify := ctx.MustGet("classify").(models.Classify)
	common.ResponseSuccess(ctx, classify)
}

func (c *Classify) List(ctx *gin.Context) {
	classifies, err := c.ClassifySvc.List()
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorMysql, "classify get error", err)
		return
	}
	common.ResponseSuccess(ctx, classifies)
}

func (c *Classify) ListByAssetSource(ctx *gin.Context) {
	var classify models.Classify
	if err := ctx.ShouldBindJSON(&classify); err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorParams, "params error: ", err)
		return
	}

	atype, err := engine.AssetTypeFromString(classify.AssetString)
	classifies, err := c.ClassifySvc.ListByAssetSource(atype, classify.Type, classify.Main, classify.Sub)
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorMysql, "classify get error", err)
		return
	}
	common.ResponseSuccess(ctx, classifies)
}

func (c *Classify) Sync(ctx *gin.Context) {
	var classify models.Classify
	if err := ctx.ShouldBindJSON(&classify); err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorParams, "params error: ", err)
		return
	}

	atype, err := engine.AssetTypeFromString(classify.AssetString)
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorEngine, "asset type illegal", err)
		return
	}

	classifySlice, err := c.Rpc.GetClassify(atype, classify.Type, classify.Main, classify.Sub)
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorEngine, "classify get error", err)
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
		common.ResponseErrorBusiness(ctx, common.ErrorEngine, "classify save error", err)
		return
	}
	common.ResponseSuccess(ctx, struct{}{})
}

func (c *Classify) Delete(ctx *gin.Context) {
	classify := ctx.MustGet("classify").(models.Classify)
	err := c.ClassifySvc.Delete(classify.Id)
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorMysql, "classify delete error", err)
		return
	}
	common.ResponseSuccess(ctx, classify.Id)
}
