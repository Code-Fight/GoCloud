package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"gocloud/datamodels"
	"gocloud/services"
	"gocloud/web/middleware"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ShareController struct {
	Ctx iris.Context
	Service services.IFileService
}

func (this *ShareController)PostCreateshare()  {
	this.Ctx.Proceed (middleware.NewAuth())
	if this.Ctx.IsStopped(){
		return
	}

	user_file_id,_ :=strconv.ParseInt( this.Ctx.Request().FormValue("user_file_id"),10,0)

	share_pwd :=this.Ctx.Request().FormValue("share_pwd")
	share_time ,_:=strconv.ParseInt( this.Ctx.Request().FormValue("share_time"),10,0)

	if user_file_id==0{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "share error:invalid user_file_id",
		})
		return
	}

	sharefile,err:=this.Service.QueryShareFileByUserFileId(user_file_id)
	if err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    err.Error(),
		})
		return
	}

	if sharefile!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "this file has been shared",
		})
		return
	}



	link,succ,err:=this.Service.CreateShareFile(user_file_id,share_time,share_pwd)
	if !succ||err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    err.Error(),
		})
		return
	}


	this.Ctx.JSON(datamodels.RespModel{
		Status: 1,
		Msg:    "",
		Data: iris.Map{
			"link":link,
		},
	})

}

func (this *ShareController) Get() mvc.View  {
	return mvc.View{
		Layout:"shared/layout.fw.html",
		Name:"error/error.html",
	}
}

func (this *ShareController) GetBy(share_id string) mvc.View  {
	sess := sessions.Get(this.Ctx)

	user, ok := (sess.Get("user")).(*datamodels.UserModel)
	username:=""
	if !ok {
		username=""
	}else {
		username = user.Username
	}

	return mvc.View{
		Name:"share/share.html",
		Data: iris.Map{
			"share_id":share_id,
			"username":username,
		},
	}

}

func (this *ShareController) GetFileBy(share_id string)  {
	if len(share_id)==0||share_id=="undefined"{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "invalid share_id",
		})
		return
	}

	share, err:=this.Service.QueryShareFileBy(share_id)
	if err !=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    err.Error(),
		})
		return
	}

	if share == nil{
		//file not exist
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "",
		})
		return
	}

	if len(share.SharePwd)==0{
		//public share
		this.Ctx.JSON(datamodels.RespModel{
			Status: 1,
			Msg:    "OK",
			Data: iris.Map{
				"pwd":0,
			},
		})
		return
	}
	//password share
	this.Ctx.JSON(datamodels.RespModel{
		Status: 1,
		Msg:    "OK",
		Data: iris.Map{
			"pwd":1,
		},
	})
}

func (this *ShareController) PostValid()   {
	share_id:=this.Ctx.Request().FormValue("share_id")
	pwd:=this.Ctx.Request().FormValue("pwd")
	if len(share_id)==0{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "Invalid share id",
		})
		return
	}

	usershare,err :=this.Service.QueryUserShareFileBy(share_id)
	if err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    err.Error(),
		})
		return
	}

	if usershare.SharePwd !=pwd{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "the share passwprd invalid",
		})
		return
	}
	usershare.SharePwd =""
	this.Ctx.JSON(datamodels.RespModel{
		Status: 1,
		Msg:    "OK",
		Data: usershare,
	})
}

func (this *ShareController) GetDownloadfileBy(filename string)  {

	share_id :=this.Ctx.Request().FormValue("share_id")
	share_pwd:=this.Ctx.Request().FormValue("share_pwd")
	if len(strings.TrimSpace(share_id)) == 0 ||len(strings.TrimSpace(share_pwd))==0{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "upload error:invalid share_id/share_pwd",
		})
		return
	}
	_,err:=this.Service.QueryShareFileAndValid(share_id,share_pwd)
	if err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "upload error:" + err.Error(),
		})
		return
	}

	share,err:=this.Service.QueryUserShareFileBy(share_id)
	if err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "upload error:" + err.Error(),
		})
		return
	}



	file, err := this.Service.GetFileMeta(share.FileQetag)
	if err != nil {
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "upload error:" + err.Error(),
		})
		return
	}

	//download
	f, err := os.Open(file.Location)
	if err != nil {
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "open file error:" + err.Error(),
		})
		return
	}
	defer f.Close()
	_t, err := time.Parse("2006-01-02 15:04:05", file.UploadAt)
	if err != nil {
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "parse time error:" + err.Error(),
		})
		return
	}
	this.Ctx.ResponseWriter().Header().Set("Content-Disposition","attachment; filename=\""+filename+"\"")

	http.ServeContent(this.Ctx.ResponseWriter(), this.Ctx.Request(), "", _t, f)
}

func (this *ShareController) PostSavefile() {

	this.Ctx.Proceed (middleware.NewAuth())
	if this.Ctx.IsStopped(){
		return
	}


	sess := sessions.Get(this.Ctx)

	user, ok := (sess.Get("user")).(*datamodels.UserModel)
	if !ok {
		this.Ctx.Application().Logger().Error("parse user err by sesssion")
		return
	}

	share_id :=this.Ctx.Request().FormValue("share_id")
	share_pwd:=this.Ctx.Request().FormValue("share_pwd")
	dir ,_:=strconv.ParseInt( this.Ctx.Request().FormValue("dir"),10,0)

	if len(strings.TrimSpace(share_id)) == 0 ||len(strings.TrimSpace(share_pwd))==0{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "upload error:invalid share_id/share_pwd",
		})
		return
	}
	fileshare,err:=this.Service.QueryShareFileAndValid(share_id,share_pwd)
	if err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "error:" + err.Error(),
		})
		return
	}

	userfile,err:=this.Service.GetUserFileByID(int64(fileshare.UserFileId))
	succ,err :=this.Service.AddUserFileRelation(user.Username,userfile.FileQetag,userfile.FileName,"",userfile.FileSize,userfile.IsDir, dir)
	if !succ||err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    err.Error(),
		})
	}

	this.Ctx.JSON(datamodels.RespModel{
		Status: 1,
		Msg:    "OK",
	})
}


func (this *ShareController) GetCancelshareBy(share_id string)  {
	this.Ctx.Proceed (middleware.NewAuth())
	if this.Ctx.IsStopped(){
		return
	}
	//idInt, _:= strconv.ParseInt(id,10,0)

	succ,err:=this.Service.CancelShareFile(share_id)

	if !succ || err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    err.Error(),
		})
		return
	}

	this.Ctx.JSON(datamodels.RespModel{
		Status: 1,
		Msg:    "OK",
	})

}

