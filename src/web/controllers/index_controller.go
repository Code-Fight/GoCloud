package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type IndexController struct {
	Ctx iris.Context
}

func (this *IndexController) Get() mvc.View{

	this.Ctx.Application().Logger().Debug("index get")
	return mvc.View{
		Name: "index/index.html",
	}
}

func (this *IndexController) GetIndex() {
	//test
	a :=[]int{1,2,3}
	this.Ctx.JSON(a)
}

