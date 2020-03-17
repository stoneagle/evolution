package main

import (
	"evolution/backend/ashare/bootstrap"
	"evolution/backend/ashare/bootstrap/database"
	"evolution/backend/ashare/bootstrap/route"
	"time"
)

func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("ashare", "wuzhongyang@wzy.com")
	app.Bootstrap()

	app.Configure(database.Configure, route.Configure)
	_, err := time.LoadLocation(app.Config.Ashare.System.Location)
	if err != nil {
		panic(err)
	}
	return app
}

func main() {
	app := newApp()
	app.Listen(":" + app.Config.Ashare.System.Port)
}
