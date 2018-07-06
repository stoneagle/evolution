package middles

import (
	"evolution/backend/common/resp"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	sessionKey = "user"
)

func BasicAuthLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)
		session := sessions.Default(c)
		session.Set(sessionKey, user)
		err := session.Save()
		if err != nil {
			resp.ErrorBusiness(c, resp.ErrorLogin, "session save error", err)
		} else {
			resp.Success(c, struct{}{})
		}
	}
}

func BasicAuthLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(sessionKey)
		if user == nil {
			resp.ErrorBusiness(c, resp.ErrorLogin, "invalid session token", nil)
		} else {
			session.Delete(sessionKey)
			session.Save()
			resp.Success(c, struct{}{})
		}
	}
}

func BasicAuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(sessionKey)
		if user == nil {
			c.JSON(http.StatusUnauthorized, resp.Response{
				Code: resp.ErrorLogin,
				Data: struct{}{},
				Desc: "invalid session token",
			})
		} else {
			c.Next()
		}
	}
}
