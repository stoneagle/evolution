package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/system/models"

	"github.com/gin-gonic/gin"
)

type User struct {
	BaseController
}

func NewUser() *User {
	User := &User{}
	User.Resource = ResourceUser
	return User
}

func (c *User) Router(router *gin.RouterGroup) {
	user := router.Group(c.Resource).Use(middles.OnInit(c))
	user.GET("/get/:id", c.One)
	user.GET("/list", c.List)
	user.POST("", c.Add)
	user.POST("/list", c.ListByCondition)
	user.PUT("/:id", c.Update)
	user.DELETE("/:id", c.Delete)
}

func (c *User) One(ctx *gin.Context) {
	user := ctx.MustGet(c.Resource).(*models.User)
	resp.Success(ctx, user)
}

func (c *User) List(ctx *gin.Context) {
	users, err := c.UserSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "user get error", err)
		return
	}
	resp.Success(ctx, users)
}

func (c *User) ListByCondition(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}
	users, err := c.UserSvc.ListWithCondition(&user)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "user get error", err)
		return
	}
	resp.Success(ctx, users)
}

func (c *User) Add(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.UserSvc.Add(user)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "user insert error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *User) Update(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorParams, "params error: ", err)
		return
	}

	err := c.UserSvc.Update(user.Id, user)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "user update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *User) Delete(ctx *gin.Context) {
	user := ctx.MustGet(c.Resource).(models.User)
	err := c.UserSvc.Delete(user.Id, user)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorDatabase, "user delete error", err)
		return
	}
	resp.Success(ctx, user.Id)
}
