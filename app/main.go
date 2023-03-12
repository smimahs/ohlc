package main

import (
	"app/database"
	"app/api"
	"github.com/kataras/iris/v12"

	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func newApp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel("debug")

	app.Use(recover.New())
	app.Use(logger.New())

	//routes
	db, _ := database.Connect()
	api.Update_API(app, db)
	api.Query_API(app, db)
	return app
}

func main() {
	app := newApp()
	//run database
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))

}
