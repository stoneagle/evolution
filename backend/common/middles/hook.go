package middles

import (
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func OnInit(c structs.ControllerGeneral) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c.Init()
		idStr := ctx.Param("id")
		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				resp.ErrorBusiness(ctx, resp.ErrorParams, "id params error", err)
				return
			}
			svcMap := c.GetResourceServiceMap()
			modelMap := c.GetResourceModelMap()
			resource := c.GetResource()
			model := modelMap[resource]
			svc := svcMap[resource]
			err = svc.One(id, model)
			if err != nil {
				resp.ErrorBusiness(ctx, resp.ErrorDatabase, fmt.Sprintf("get model %s error", resource), err)
				return
			}
			ctx.Set(resource, model)
		}
		ctx.Next()
	}
}
