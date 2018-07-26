package controllers

import (
	"evolution/backend/common/middles"

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
	user.POST("", c.Create)
	user.POST("/list", c.List)
	user.POST("/count", c.Count)
	user.PUT("/:id", c.Update)
	user.DELETE("/:id", c.Delete)
}
