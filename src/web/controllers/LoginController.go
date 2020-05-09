package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"gocloud/common"
	"gocloud/datamodels"
	"gocloud/services"
	"strings"
)

type LoginController struct {
	Ctx     iris.Context
	Service services.ILoginService

}

func (this *LoginController) Get() mvc.View {
	sess:=sessions.Get(this.Ctx)
	if auth,err:=sess.GetBoolean("authenticated");auth&&err==nil{
		return mvc.View{
			Name: "index/index.html",
		}
	}
	return mvc.View{
		Layout:"shared/layout.fw.html",
		Name:"login/login.html",
	}
}

// PostSignin the user signin
func (this *LoginController) PostSignin()  {
	data :=&datamodels.UserModel{}

	this.Ctx.ReadJSON(data)

	if len(strings.TrimSpace(data.Email))==0{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg: "Invalid Email",
		})
		return
	}

	if len(strings.TrimSpace(data.Userpwd))==0{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg: "Invalid Password",
		})

		return
	}


	user,err :=this.Service.Signin(data)
	if err!=nil{
		this.Ctx.Application().Logger().Error(err)
		this.Ctx.JSON("system error")
		return
	}
	if user.Userpwd != common.Sha1([]byte(data.Userpwd+common.User_Pwd_Sha1_Salt)){
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg: "The Email or Password Error",
		})
		return
	}

	sessions.Get(this.Ctx).Set("authenticated",true)
	sessions.Get(this.Ctx).Set("user",user)
	this.Ctx.JSON(datamodels.RespModel{
		Status: 0,
		Msg: "OK",
	})
}

// PostSignup the user singUp
func (this *LoginController) PostSignup()  {
	data :=&datamodels.UserModel{}

	this.Ctx.ReadJSON(data)

	_,err :=this.Service.Signup(data)
	if err!=nil{
		this.Ctx.Application().Logger().Error(err)
		this.Ctx.JSON("system error")
		return
	}
	this.Ctx.JSON("OK")
}

// GetLogout the user loguout
func (this *LoginController) GetLogout() {
	sessions.Get(this.Ctx).Set("authenticated",false)
	this.Ctx.JSON(datamodels.RespModel{
		Status: 0,
		Msg: "OK",
	})
}
