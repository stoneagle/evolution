package route

import (
	"evolution/backend/quant/bootstrap"
	c "evolution/backend/quant/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Configure(b *bootstrap.Bootstrapper) {
	b.App.StaticFS("/static", http.Dir("./static"))
	b.App.StaticFile("/favicon.ico", "./static/favicon.ico")
	b.App.GET("/index", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	prefix := b.Config.Quant.System.Prefix + "/v1"
	v1 := b.App.Group(prefix)
	{
		c.NewPool().Router(v1)
		c.NewClassify().Router(v1)
		c.NewConfig().Router(v1)
		c.NewItem(b.Websocket).Router(v1)
		c.NewSignal().Router(v1)
	}
}
