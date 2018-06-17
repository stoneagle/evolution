package controllers

import (
	"quant/backend/common"
	"quant/backend/models"
	"quant/backend/services"
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
	pool.PUT("/:id", initPool(c.PoolSvc), c.Update)
	pool.DELETE("/:id", initPool(c.PoolSvc), c.Delete)
}

func initPool(svc *services.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			common.ResponseErrorBusiness(ctx, common.ErrorParams, "id params error", err)
		}

		pool, err := svc.One(id)
		if err != nil {
			common.ResponseErrorBusiness(ctx, common.ErrorMysql, "get pool error", err)
		}

		ctx.Set("pool", pool)
		return
	}
}

func (c *Pool) One(ctx *gin.Context) {
	pool := ctx.MustGet("pool").(models.Pool)
	common.ResponseSuccess(ctx, pool)
}

func (c *Pool) List(ctx *gin.Context) {
	pools, err := c.PoolSvc.List()
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorMysql, "pool get error", err)
		return
	}
	common.ResponseSuccess(ctx, pools)
}

func (c *Pool) Add(ctx *gin.Context) {
	var pool models.Pool
	if err := ctx.ShouldBindJSON(&pool); err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorParams, "params error: ", err)
		return
	}

	err := c.PoolSvc.Add(&pool)
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorMysql, "pool insert error", err)
		return
	}
	common.ResponseSuccess(ctx, struct{}{})
}

func (c *Pool) Update(ctx *gin.Context) {
	var pool models.Pool
	if err := ctx.ShouldBindJSON(&pool); err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorParams, "params error: ", err)
		return
	}

	err := c.PoolSvc.Update(pool.Id, &pool)
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorMysql, "pool update error", err)
		return
	}
	common.ResponseSuccess(ctx, struct{}{})
}

func (c *Pool) Delete(ctx *gin.Context) {
	pool := ctx.MustGet("pool").(models.Pool)
	err := c.PoolSvc.Delete(pool.Id)
	if err != nil {
		common.ResponseErrorBusiness(ctx, common.ErrorMysql, "pool delete error", err)
		return
	}
	common.ResponseSuccess(ctx, pool.Id)
}
