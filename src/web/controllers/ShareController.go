package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"gocloud/services"
)

type ShareController struct {
	Ctx iris.Context
	Service services.IFileService
}

func (this *ShareController) GetBy(qetag string) mvc.View  {
	if len(qetag)==0{
		return mvc.View{
			Layout:"shared/layout.fw.html",
			Name:"error/error.html",
		}
	}

	share, err:=this.Service.QueryShareFileBy(qetag)
	if err !=nil{
		return mvc.View{
			Layout:"shared/layout.fw.html",
			Name:"error/error.html",
		}
	}

	if len(share.SharePwd)==0{
		//public share
		return mvc.View{

		}
	}
	//password share
	return mvc.View{}

}
