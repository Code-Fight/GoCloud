package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"gocloud/datamodels"
	"gocloud/services"
)

type IndexController struct {
	Ctx iris.Context
	IndexService services.IIndexService
}

func (this *IndexController) Get() mvc.View{
	sess:=sessions.Get(this.Ctx)
	//if auth,err:=sess.GetBoolean("authenticated");!auth||err!=nil{
	//	return mvc.View{
	//		Layout: "shared/layout.fw.html",
	//		Name: "login/login.html",
	//	}
	//}

	user ,ok:=(sess.Get("user")).(*datamodels.UserModel)
	if !ok{
		this.Ctx.Application().Logger().Error("get user err by sesssion")
	}

	return mvc.View{
		Name: "index/index.html",
		Data: iris.Map{
			"username":user.Username,
		},
	}
}


