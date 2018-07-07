package middles

import (
	systemApi "evolution/backend/common/api/system"
	"evolution/backend/common/resp"
	"evolution/backend/common/structs"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	UserKey = "user"
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
				return
			}
			ctx.Set(name, model)
		}
		return
	}
}

func UserFromSession(sessionKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		username := session.Get(sessionKey)
		if username == nil {
			ctx.AbortWithStatusJSON(http.StatusOK, resp.Response{
				Code: resp.ErrorLogin,
				Data: struct{}{},
				Desc: "invalid user session",
			})
			return
		} else {
			name := username.(string)
			user, err := systemApi.UserByName(name)
			if err != nil {
				resp.ErrorBusiness(ctx, resp.ErrorApi, fmt.Sprintf("get user %s error", name), err)
				return
			}
			ctx.Set(UserKey, user)
		}
		return
	}
}
