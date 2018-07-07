package middles

import (
	"evolution/backend/common/resp"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	sessionBA = "basicAuth"
)

func BasicAuthLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		session := sessions.Default(c)
		session.Set(sessionBA, user)
		err := session.Save()
		if err != nil {
			resp.ErrorBusiness(c, resp.ErrorLogin, "session save error", err)
		} else {
			resp.Success(c, user)
		}
	}
}

func BasicAuthLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(sessionBA)
		if user == nil {
			resp.ErrorBusiness(c, resp.ErrorLogin, "invalid session token", nil)
		} else {
			session.Delete(sessionBA)
			session.Save()
			resp.Success(c, struct{}{})
		}
	}
}

func BasicAuthCurrent() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(sessionBA)
		if user == nil {
			resp.ErrorBusiness(c, resp.ErrorLogin, "invalid session token", nil)
		} else {
			resp.Success(c, user)
		}
	}
}

func BasicAuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(sessionBA)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp.Response{
				Code: resp.ErrorLogin,
				Data: struct{}{},
				Desc: "invalid session token",
			})
			return
		} else {
			c.Next()
		}
	}
}
