package main

import (
	"evolution/backend/quant/bootstrap"
	"evolution/backend/quant/bootstrap/database"
	"evolution/backend/quant/bootstrap/route"
	"time"
)

func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("quant", "wuzhongyang@wzy.com")
	app.Bootstrap()
	app.Configure(database.Configure, route.Configure)

	_, err := time.LoadLocation(app.Config.Quant.System.Location)
	if err != nil {
		panic(err)
	}
	return app
}

func main() {
	app := newApp()
	app.Listen(":8080")
}
