package route

import (
	"net/http"
	"quant/backend/bootstrap"
	c "quant/backend/controllers"
	"time"

	"github.com/gin-contrib/cors"
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
		c.NewStock().Router(v1)
		c.NewSignal().Router(v1)
		c.NewPool().Router(v1)
		c.NewClassify().Router(v1)
	}

	if b.Config.App.Mode == "debug" {
		b.App.Use(cors.New(cors.Config{
			AllowHeaders:     []string{"Content-Type", "Access-Control-Allow-Origin", "Authorization"},
			AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
			AllowCredentials: true,
			AllowOrigins:     []string{"http://localhost:8080", "http://localhost:6999"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowOriginFunc: func(origin string) bool {
				return origin == "https://localhost:6999"
			},
			MaxAge: 12 * time.Hour,
		}))
	}
}
