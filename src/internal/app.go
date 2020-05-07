package internal

import (
	"database/sql"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"gocloud/common"
	"gocloud/services"
	"gocloud/web/controllers"
)

func AppRun() {
	//1. create iris instance
	app := iris.New()
	//2. set logger level
	app.Logger().SetLevel("debug")
	//3. register view
	tmplate := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(true)
	tmplate_framework := iris.HTML("./web/views", ".html").Layout("shared/layout.fw.html").Reload(true)
	app.RegisterView(tmplate)
	app.RegisterView(tmplate_framework)

	app.StaticWeb("/assets", "./web/assets")


	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewLayout("shared/layout.fw.html")
		ctx.View("error/error.html")
	})
	db := getDb()
	registerHandler(app,db)

	// run web app
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)


}

func getDb() *sql.DB{
	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		panic(err)
	}
	return db
}

func registerHandler(app *iris.Application,db *sql.DB) {

	//productRepository := repositories.NewProductManager("product", db)
	indexService := services.NewIndexService()
	productParty := app.Party("/")
	product := mvc.New(productParty)
	product.Register( indexService)
	product.Handle(new(controllers.IndexController))
}
