package middles

import (
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func One(svc structs.ServiceGeneral, name string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				resp.ErrorBusiness(ctx, resp.ErrorParams, "id params error", err)
				return
			}
			model, err := svc.One(id)
			if err != nil {
				resp.ErrorBusiness(ctx, resp.ErrorMysql, fmt.Sprintf("get model %s error", name), err)
			}
			ctx.Set(name, model)
		}
		return
	}
}
