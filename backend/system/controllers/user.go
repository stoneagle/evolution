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
	user := router.Group(c.Resource.String()).Use(middles.OnInit(c))
	user.GET("/get/:id", c.One)
	user.GET("/list", c.List)
	user.POST("", c.Create)
	user.POST("/list", c.ListByCondition)
	user.PUT("/:id", c.Update)
	user.DELETE("/:id", c.Delete)
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
