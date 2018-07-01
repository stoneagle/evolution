package controllers

import (
	"evolution/backend/common/resp"
	"evolution/backend/quant/models"
	"evolution/backend/quant/rpc/engine"
	"evolution/backend/quant/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pool struct {
	Base
	PoolSvc *services.Pool
}

func NewPool() *Pool {
	Pool := &Pool{}
	Pool.Prepare()
	Pool.PoolSvc = services.NewPool(Pool.Engine, Pool.Cache)
	return Pool
}

func (c *Pool) Router(router *gin.RouterGroup) {
	pool := router.Group("pool")
	pool.GET("/get/:id", initPool(c.PoolSvc), c.One)
	pool.GET("/list", c.List)
	pool.POST("", c.Add)
	pool.POST("/items/add", c.AddItems)
	pool.POST("/items/delete", c.DeleteBatchItems)
	pool.PUT("/:id", initPool(c.PoolSvc), c.Update)
	pool.DELETE("/:id", initPool(c.PoolSvc), c.Delete)
}

func initPool(svc *services.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			resp.ErrorBusiness(ctx, resp.ErrorParams, "id params error", err)
		}

		pool, err := svc.One(id)
		if err != nil {
			resp.ErrorBusiness(ctx, resp.ErrorMysql, "get pool error", err)
		}

		ctx.Set("pool", pool)
		return
	}
}

func (c *Pool) One(ctx *gin.Context) {
	pool := ctx.MustGet("pool").(models.Pool)
	resp.Success(ctx, pool)
}

func (c *Pool) List(ctx *gin.Context) {
	pools, err := c.PoolSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "pool get error", err)
		return
	}
	resp.Success(ctx, pools)
}

func (c *Pool) Add(ctx *gin.Context) {
	var pool models.Pool
	if err := ctx.ShouldBindJSON(&pool); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	atype, err := engine.AssetTypeFromString(pool.AssetString)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorEngine, "asset type illegal", err)
		return
	}
	pool.Asset = atype

	err = c.PoolSvc.Add(&pool)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "pool insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Pool) Update(ctx *gin.Context) {
	var pool models.Pool
	if err := ctx.ShouldBindJSON(&pool); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.PoolSvc.Update(pool.Id, &pool)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "pool update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Pool) Delete(ctx *gin.Context) {
	pool := ctx.MustGet("pool").(models.Pool)
	err := c.PoolSvc.Delete(pool.Id)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "pool delete error", err)
		return
	}
	resp.Success(ctx, pool.Id)
}

func (c *Pool) AddItems(ctx *gin.Context) {
	var pool models.Pool
	if err := ctx.ShouldBindJSON(&pool); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.PoolSvc.AddItems(&pool)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "pool items join insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *Pool) DeleteBatchItems(ctx *gin.Context) {
	var pool models.Pool
	if err := ctx.ShouldBindJSON(&pool); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.PoolSvc.DeleteItems(&pool)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "pool items join delete error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}
