package main

import (
	"evolution/backend/system/bootstrap"
	"evolution/backend/system/bootstrap/database"
	"evolution/backend/system/bootstrap/route"
	"time"
)

func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("system", "wuzhongyang@wzy.com")
	app.Bootstrap()
	app.Configure(database.Configure, route.Configure)
	_, err := time.LoadLocation(app.Config.System.System.Location)
	if err != nil {
		panic(err)
	}
	return app
}

func main() {
	app := newApp()
	app.Listen(":8080")
}
