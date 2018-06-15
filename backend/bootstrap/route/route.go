package route

import (
	"net/http"
	"quant/backend/bootstrap"
	c "quant/backend/controllers"

	"github.com/gin-gonic/gin"
)

func Configure(b *bootstrap.Bootstrapper) {
	b.App.StaticFS("/static", http.Dir("./static"))
	b.App.StaticFile("/favicon.ico", "./static/favicon.ico")
	b.App.GET("/index", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	v1 := b.App.Group("/v1")
	{
		c.NewPool().Router(v1)
		c.NewClassify().Router(v1)
		c.NewConfig().Router(v1)
		c.NewItem().Router(v1)
		c.NewSignal().Router(v1)
	}
}
