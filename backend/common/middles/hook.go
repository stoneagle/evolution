package middles

import (
	"evolution/backend/common/structs"

	"github.com/gin-gonic/gin"
)

func OnInit(controller structs.ControllerGeneral) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		controller.Init()
		ctx.Next()
	}
}
