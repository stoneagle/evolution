package middles

import (
	systemApi "evolution/backend/common/api/system"
	"evolution/backend/common/resp"
	"fmt"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	UserKey = "user"
)

func UserFromSession(sessionKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		username := session.Get(sessionKey)
		if username == nil {
			ctx.AbortWithStatusJSON(http.StatusOK, resp.Response{
				Code: resp.ErrorSign,
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
