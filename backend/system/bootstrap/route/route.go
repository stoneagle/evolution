package route

import (
	"evolution/backend/system/bootstrap"
	"evolution/backend/system/controllers"
)

func Configure(b *bootstrap.Bootstrapper) {
	prefix := b.Config.System.System.Prefix + "/v1"
	v1 := b.App.Group(prefix)
	{
		controllers.NewUser().Router(v1)
	}
}
