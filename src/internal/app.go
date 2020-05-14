package internal

import (
	"database/sql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"gocloud/common"
	"gocloud/dao"
	"gocloud/services"
	"gocloud/web/controllers"
	"gocloud/web/middleware"
	"time"
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

	app.HandleDir("/assets", "./web/assets")


	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewLayout("shared/layout.fw.html")
		ctx.View("error/error.html")
	})

	sess := sessions.New(sessions.Config{
		// Cookie string, the session's client cookie name, for example: "_session_id"
		//
		// Defaults to "irissessionid"
		Cookie: "_session_id",
		// it's time.Duration, from the time cookie is created, how long it can be alive?
		// 0 means no expire, unlimited life.
		// -1 means expire when browser closes
		// or set a value, like 2 hours:
		Expires: time.Hour * 2,
		// if you want to invalid cookies on different subdomains
		// of the same host, then enable it.
		// Defaults to false.
		DisableSubdomainPersistence: false,
		// Allow getting the session value stored by the request from the same request.
		AllowReclaim: true,
	})

	app.Use(sess.Handler())

	db := getDb()
	registerHandler(app,db)


	// run web app
	app.Listen(":8080",iris.WithOptimizations)


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

	//Index
	indexService := services.NewIndexService()
	indexParty := app.Party("/")
	indexParty.Use(middleware.NewAuth())
	index := mvc.New(indexParty)
	index.Register(indexService)
	index.Handle(new(controllers.IndexController))

	//Login
	loginDao :=dao.NewUserDao("tbl_user",db)
	loginService := services.NewLoginService(loginDao)
	loginParty := app.Party("/login")
	login := mvc.New(loginParty)
	login.Register(loginService)
	login.Handle(new(controllers.LoginController))

	//File
	fileDao :=dao.NewFileDao(db)
	fileService :=services.NewFileService(fileDao)
	file:=mvc.New(app.Party("/file",middleware.NewAuth()))
	file.Register(fileService)
	file.Handle(new(controllers.FileController))


}
