package route

import (
	"evolution/backend/time/bootstrap"
	"evolution/backend/time/controllers"
)

func Configure(b *bootstrap.Bootstrapper) {
	prefix := b.Config.Time.System.Prefix + "/v1"
	v1 := b.App.Group(prefix)
	{
		controllers.NewTask().Router(v1)
		controllers.NewCountry().Router(v1)
		controllers.NewArea().Router(v1)
		entity := v1.Group("entity")
		{
			controllers.NewEntityLife().Router(entity)
		}
	}
}
