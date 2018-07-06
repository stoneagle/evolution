package route

import (
	sysApi "evolution/backend/common/api/system"
	"evolution/backend/common/middles"
	"evolution/backend/time/bootstrap"
	"evolution/backend/time/controllers"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Configure(b *bootstrap.Bootstrapper) {
	store := sessions.NewCookieStore([]byte("secret"))
	b.App.Use(sessions.Sessions("session", store))

	prefix := b.Config.Time.System.Prefix + "/" + b.Config.Time.System.Version
	var v1 *gin.RouterGroup
	switch b.Config.Time.System.Auth {
	case "BasicAuth":
		BAConf, err := sysApi.UserBAMap(b.Config.System.System)
		if err != nil {
			panic(err)
		}
		sign := b.App.Group(prefix + "/sign")
		sign.GET("/login", gin.BasicAuth(BAConf), middles.BasicAuthLogin())
		sign.GET("/logout", middles.BasicAuthLogout())
		v1 = b.App.Group(prefix, middles.BasicAuthCheck())
	default:
		v1 = b.App.Group(prefix)
	}
	{
		controllers.NewTask().Router(v1)
		controllers.NewCountry().Router(v1)
		controllers.NewArea().Router(v1)
		controllers.NewField().Router(v1)
		controllers.NewPhase().Router(v1)
		controllers.NewEntity().Router(v1)
	}
}
