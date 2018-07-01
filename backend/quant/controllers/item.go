package controllers

import (
	"encoding/json"
	"evolution/backend/common/resp"
	"evolution/backend/quant/models"
	"evolution/backend/quant/rpc/engine"
	"evolution/backend/quant/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Item struct {
	Base
	ItemSvc     *services.Item
	ClassifySvc *services.Classify
	Ws          *resp.Websocket
}

func NewItem(ws *resp.Websocket) *Item {
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
	item.GET("/point/:id", initItem(c.ItemSvc), c.SyncPoint)
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
			resp.ErrorBusiness(ctx, resp.ErrorParams, "id params error", err)
		}

		item, err := svc.One(id)
		if err != nil {
			resp.ErrorBusiness(ctx, resp.ErrorMysql, "get item error", err)
		}

		ctx.Set("item", item)
		return
	}
}

func (c *Item) One(ctx *gin.Context) {
	item := ctx.MustGet("item").(models.Item)
	resp.Success(ctx, item)
}

func (c *Item) List(ctx *gin.Context) {
	var item models.Item
	if err := ctx.ShouldBindJSON(&item); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	if len(item.Classify) > 0 && item.Classify[0].AssetString != "" {
		atype, err := engine.AssetTypeFromString(item.Classify[0].AssetString)
		if err != nil {
			resp.ErrorBusiness(ctx, resp.ErrorEngine, "asset type illegal", err)
			return
		}
		item.Classify[0].Asset = atype
	}

	items, err := c.ItemSvc.List(&item)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "item get error", err)
		return
	}
	resp.Success(ctx, items)
}

func (c *Item) SyncClassify(ctx *gin.Context) {
	var classify models.Classify
	if err := ctx.ShouldBindJSON(&classify); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	items, err := c.Rpc.GetItem(classify)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorEngine, "get item error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorEngine, "classify save error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Item) WsSyncClassify(sourceJson []byte) resp.WebsocketResponse {
	var classify models.Classify
	err := json.Unmarshal(sourceJson, &classify)
	if err != nil {
		return c.Ws.BusinessError(resp.ErrorParams, "source json unmarshal failed", err)
	}

	items, err := c.Rpc.GetItem(classify)
	if err != nil {
		return c.Ws.BusinessError(resp.ErrorEngine, "get item error", err)
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
		return c.Ws.BusinessError(resp.ErrorMysql, "classify save error", err)
	}
	return c.Ws.Message(classify.Name + " sync success")
}

func (c *Item) SyncPoint(ctx *gin.Context) {
	item := ctx.MustGet("item").(models.Item)
	err := c.Rpc.GetItemPoint(item.Classify[0], item.Code)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorEngine, "get item point error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}
