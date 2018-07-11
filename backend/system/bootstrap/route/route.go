package route

import (
	"evolution/backend/common/database"
	"evolution/backend/common/middles"
	"evolution/backend/common/resp"
	"evolution/backend/system/bootstrap"
	"evolution/backend/system/controllers"
	"evolution/backend/system/services"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Configure(b *bootstrap.Bootstrapper) {
	prefix := b.Config.System.System.Prefix + "/" + b.Config.System.System.Version
	var v1 *gin.RouterGroup

	switch b.Config.System.Auth.Type {
	case middles.TypeBasicAuth:
		store, err := database.SessionByRedis(b.Config.System.Redis)
		if err != nil {
			panic(err)
		}
		b.App.Use(sessions.Sessions(b.Config.System.Auth.Session, store))

		BAConf := GetBAList(b.Config.System.System.Name)
		sign := b.App.Group(prefix + "/sign")
		sign.GET("/login", gin.BasicAuth(BAConf), middles.BasicAuthLogin())
		sign.GET("/logout", middles.BasicAuthLogout())
		sign.GET("/current", middles.BasicAuthCurrent())
		v1 = b.App.Group(prefix, middles.BasicAuthCheck())
	case middles.TypeBAJwt:
		BAConf := GetBAList(b.Config.System.System.Name)
		sign := b.App.Group(prefix + "/sign")
		sign.GET("/login", gin.BasicAuth(BAConf), middles.JWTAuthLogin())
		sign.GET("/current", middles.JWTAuthCurrent())
		sign.GET("/logout", func(c *gin.Context) { resp.Success(c, struct{}{}) })
		v1 = b.App.Group(prefix, middles.JWTAuthCheck())
	default:
		v1 = b.App.Group(prefix)
	}
	{
		controllers.NewUser().Router(v1)
	}
}

func GetBAList(projectName string) map[string]string {
	UserSlice, err := services.NewUser(database.GetXorm(projectName), nil).List()
	if err != nil {
		panic(err)
	}
	BAConf := make(map[string]string)
	for _, one := range UserSlice {
		BAConf[one.Name] = one.Password
	}
	return BAConf
}
