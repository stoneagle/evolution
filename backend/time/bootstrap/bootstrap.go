package bootstrap

import (
	"evolution/backend/common/config"
	"evolution/backend/common/middles"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Configurator func(*Bootstrapper)

type Bootstrapper struct {
	App          *gin.Engine
	AppName      string
	AppOwner     string
	AppSpawnDate time.Time
	Config       *config.Conf
}

func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
		App:          gin.New(),
		Config:       config.Get(),
	}
	return b
}

func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	gin.SetMode(b.Config.App.Mode)
	b.App.Use(gin.Logger())
	b.App.Use(gin.Recovery())
	b.App.Use(middles.Recovery(middles.RecoveryHandler))
	// cors must set in bootstrap
	if b.Config.App.Mode == "debug" {
		b.App.Use(cors.New(cors.Config{
			AllowHeaders:     []string{"Content-Type", "Access-Control-Allow-Origin", "Authorization"},
			AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
			AllowCredentials: true,
			AllowOrigins:     b.Config.Time.System.Cors,
			ExposeHeaders:    []string{"Content-Length"},
			AllowOriginFunc: func(origin string) bool {
				return origin == "http://localhost:8080"
			},
			MaxAge: 12 * time.Hour,
		}))
	}
	return b
}

func (b *Bootstrapper) Listen(addr string) {
	b.App.Run(addr)
}

func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}
