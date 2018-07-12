package main

import (
	"evolution/backend/time/bootstrap"
	"evolution/backend/time/bootstrap/database"
	"evolution/backend/time/bootstrap/route"
)

func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("time", "wuzhongyang@wzy.com")
	app.Bootstrap()
	app.Configure(database.Configure, route.Configure)
	return app
}

func main() {
	app := newApp()
	app.Listen(":8080")
}
