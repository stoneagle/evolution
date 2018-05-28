package bootstrap

import (
	"quant/backend/common"
	"time"

	"github.com/gin-gonic/gin"
)

type Configurator func(*Bootstrapper)

type Bootstrapper struct {
	App          *gin.Engine
	AppName      string
	AppOwner     string
	AppSpawnDate time.Time
	Config       *common.Conf
}

func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
		App:          gin.New(),
		Config:       common.GetConfig(),
	}

	return b
}

func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	gin.SetMode(b.Config.App.Mode)
	b.App.Use(gin.Logger())
	b.App.Use(gin.Recovery())

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
