package controllers

import (
	"encoding/json"
	"quant/backend/common"
	"quant/backend/models"
	"quant/backend/rpc/engine"
	"quant/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Item struct {
	Base
	ItemSvc     *services.Item
	ClassifySvc *services.Classify
	Ws          *common.Websocket
}

func NewItem(ws *common.Websocket) *Item {
	Item := &Item{}
	Item.Prepare()
	Item.ItemSvc = services.NewItem(Item.Engine, Item.Cache)
	Item.ClassifySvc = services.NewClassify(Item.Engine, Item.Cache)
	Item.Ws = ws
	return Item
}

func (c *Item) Router(router *gin.RouterGroup) {
	item := router.Group("item")
	item.POST("/list", c.List)
	item.GET("/get/:id", initItem(c.ItemSvc), c.One)
	item.GET("/sync/classify/ws", func(ctx *gin.Context) {
		wsCtx := c.Ws.BuildContext(c.WsSyncClassify)
		c.Ws.Intence.HandleRequestWithKeys(ctx.Writer, ctx.Request, wsCtx)
	})
	item.POST("/sync/classify", c.SyncClassify)
}

func initItem(svc *services.Item) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			common.ResponseErrorBusiness(ctx, common.ErrorParams, "id params error", err)
		}

		item, err := svc.One(id)
		if err != nil {
			common.ResponseErrorBusiness(ctx, common.ErrorMysql, "get item error", err)
		}

		ctx.Set("item", item)
		return
	}
}

func (c *Item) One(ctx *gin.Context) {
	item := ctx.MustGet("item").(models.Item)
	common.ResponseSuccess(ctx, item)
}

func (c *Item) List(ctx *gin.Context) {
	var item models.Item
	if err := ctx.ShouldBindJSON(&item); err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorParams, "params error: ", err)
		return
	}

	if len(item.Classify) > 0 && item.Classify[0].AssetString != "" {
		atype, err := engine.AssetTypeFromString(item.Classify[0].AssetString)
		if err != nil {
			common.ResponseErrorBusiness(ctx, common.ErrorEngine, "asset type illegal", err)
			return
		}
		item.Classify[0].Asset = atype
	}

	items, err := c.ItemSvc.ListWithClassify(&item)
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorMysql, "item get error", err)
		return
	}
	common.ResponseSuccess(ctx, items)
}

func (c *Item) SyncClassify(ctx *gin.Context) {
	var classify models.Classify
	if err := ctx.ShouldBindJSON(&classify); err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorParams, "params error: ", err)
		return
	}
	items, err := c.Rpc.GetItem(classify)
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorEngine, "get item error", err)
		return
	}

	itemBatch := []models.Item{}
	for _, i := range items {
		one := models.Item{}
		one.Name = i.Name
		one.Code = i.Code
		itemBatch = append(itemBatch, one)
	}

	err = c.ItemSvc.BatchSave(classify, itemBatch)
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorEngine, "classify save error", err)
		return
	}
	common.ResponseSuccess(ctx, struct{}{})
}

func (c *Item) WsSyncClassify(sourceJson []byte) common.WebsocketResponse {
	var classify models.Classify
	err := json.Unmarshal(sourceJson, &classify)
	if err != nil {
		return c.Ws.ResponseBusinessError(common.ErrorParams, "source json unmarshal failed", err)
	}

	items, err := c.Rpc.GetItem(classify)
	if err != nil {
		return c.Ws.ResponseBusinessError(common.ErrorEngine, "get item error", err)
	}

	itemBatch := []models.Item{}
	for _, i := range items {
		one := models.Item{}
		one.Name = i.Name
		one.Code = i.Code
		itemBatch = append(itemBatch, one)
	}

	err = c.ItemSvc.BatchSave(classify, itemBatch)
	if err != nil {
		return c.Ws.ResponseBusinessError(common.ErrorMysql, "classify save error", err)
	}
	return c.Ws.ResponseMessage(classify.Name + " sync success")
}
