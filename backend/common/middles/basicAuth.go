package middles

import (
	"evolution/backend/common/resp"
	systemModel "evolution/backend/system/models"
	"evolution/backend/system/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	SessionBasicAuth = "basicAuth"
	TypeBasicAuth    = "BasicAuth"
	TypeBAJwt        = "BAJwt"
)

func BasicAuthLogin(userSvc *services.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.MustGet(gin.AuthUserKey).(string)
		session := sessions.Default(c)
		session.Set(SessionBasicAuth, username)
		err := session.Save()
		if err != nil {
			resp.ErrorBusiness(c, resp.ErrorSign, "session save error", err)
		} else {
			condition := systemModel.User{
				Name: username,
			}
			user, err := userSvc.OneByCondition(&condition)
			if err != nil {
				resp.ErrorBusiness(c, resp.ErrorApi, fmt.Sprintf("get user %s error", username), err)
				return
			}
			user.Password = ""
			resp.Success(c, user)
		}
	}
}

func BasicAuthLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(SessionBasicAuth)
		if user == nil {
			resp.ErrorBusiness(c, resp.ErrorSign, "invalid session token", nil)
		} else {
			session.Delete(SessionBasicAuth)
			session.Save()
			resp.Success(c, struct{}{})
		}
	}
}

func BasicAuthCurrent(userSvc *services.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get(SessionBasicAuth)
		if username == nil {
			resp.ErrorBusiness(c, resp.ErrorSign, "invalid session token", nil)
		} else {
			name := username.(string)
			condition := systemModel.User{
				Name: name,
			}
			user, err := userSvc.OneByCondition(&condition)
			if err != nil {
				resp.ErrorBusiness(c, resp.ErrorApi, fmt.Sprintf("get user %s error", name), err)
				return
			}
			user.Password = ""
			resp.Success(c, user)
		}
	}
}

func BasicAuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(SessionBasicAuth)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp.Response{
				Code: resp.ErrorSign,
				Data: struct{}{},
				Desc: "invalid session token",
			})
			return
		} else {
			c.Next()
		}
	}
}
