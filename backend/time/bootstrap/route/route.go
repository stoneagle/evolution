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

	prefix := b.Config.Time.System.Prefix + "/" + b.Config.Time.System.Version
	var v1 *gin.RouterGroup
	switch b.Config.System.Auth.Type {
	case middles.TypeBasicAuth:
		store, err := database.SessionByRedis(b.Config.System.Redis)
		if err != nil {
			panic(err)
		}
		b.App.Use(sessions.Sessions(b.Config.System.Auth.Session, store))
		v1 = b.App.Group(prefix, middles.BasicAuthCheck())
	case middles.TypeBAJwt:
		v1 = b.App.Group(prefix, middles.JWTAuthCheck())
	default:
		v1 = b.App.Group(prefix)
	}
	{
		controllers.NewCountry().Router(v1)
		controllers.NewArea().Router(v1)
		controllers.NewField().Router(v1)
		controllers.NewPhase().Router(v1)
		controllers.NewResource().Router(v1)
		controllers.NewUserResource().Router(v1)
		controllers.NewQuest().Router(v1)
		controllers.NewQuestTarget().Router(v1)
		controllers.NewQuestTeam().Router(v1)
		controllers.NewProject().Router(v1)
		controllers.NewTask().Router(v1)
	}
}
