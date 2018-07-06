package controllers

import (
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"evolution/backend/system/models"
	"evolution/backend/system/services"

	"github.com/gin-gonic/gin"
)

type User struct {
	structs.Controller
	Name    string
	UserSvc *services.User
}

func NewUser() *User {
	User := &User{}
	User.Init()
	User.Name = "user"
	User.ProjectName = User.Config.System.System.Name
	User.Prepare()
	User.UserSvc = services.NewUser(User.Engine, User.Cache)
	return User
}

func (c *User) Router(router *gin.RouterGroup) {
	user := router.Group(c.Name).Use(middles.One(c.UserSvc, c.Name))
	user.GET("/get/:id", c.One)
	user.GET("/list", c.List)
	user.POST("", c.Add)
	user.POST("/list", c.ListByCondition)
	user.PUT("/:id", c.Update)
	user.DELETE("/:id", c.Delete)
}

func (c *User) One(ctx *gin.Context) {
	user := ctx.MustGet(c.Name).(models.User)
	resp.Success(ctx, user)
}

func (c *User) List(ctx *gin.Context) {
	users, err := c.UserSvc.List()
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "user get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "user get error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "user insert error", err)
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
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "user update error", err)
		return
	}
	resp.Success(ctx, struct{}{})
}

func (c *User) Delete(ctx *gin.Context) {
	user := ctx.MustGet(c.Name).(models.User)
	err := c.UserSvc.Delete(user.Id, user)
	if err != nil {
		resp.ErrorBusiness(ctx, resp.ErrorMysql, "user delete error", err)
		return
	}
	resp.Success(ctx, user.Id)
}
