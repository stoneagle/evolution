package main

import (
	"evolution/backend/quant/bootstrap"
	"evolution/backend/quant/bootstrap/database"
	"evolution/backend/quant/bootstrap/route"
)

func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("quant", "wuzhongyang@wzy.com")
	app.Bootstrap()
	app.Configure(database.Configure, route.Configure)

	return app
}

func main() {
	app := newApp()
	app.Listen(":8080")
}
