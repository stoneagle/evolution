package route

import (
	"quant/backend/bootstrap"
	c "quant/backend/controllers"
)

func Configure(b *bootstrap.Bootstrapper) {
	v1 := b.App.Group("/v1")
	{
		c.NewStock().Router(v1)
		c.NewSignal().Router(v1)
		c.NewPool().Router(v1)
		c.NewClassify().Router(v1)
	}
}
