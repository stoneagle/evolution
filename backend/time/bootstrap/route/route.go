package route

import (
	"evolution/backend/common/database"
	"evolution/backend/common/middles"
	"evolution/backend/time/bootstrap"
	"evolution/backend/time/controllers"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Configure(b *bootstrap.Bootstrapper) {
	store, err := database.SessionByRedis(b.Config.System.Redis)
	if err != nil {
		panic(err)
	}
	b.App.Use(sessions.Sessions(b.Config.System.Auth.Session, store))

	prefix := b.Config.Time.System.Prefix + "/" + b.Config.Time.System.Version
	var v1 *gin.RouterGroup
	switch b.Config.System.Auth.Type {
	case "BasicAuth":
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
