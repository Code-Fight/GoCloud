package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"gocloud/common"
	"gocloud/datamodels"
	"gocloud/services"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type FileController struct {
	Ctx iris.Context
	Service services.IFileService
}

func(this * FileController) PostUpload()  {
	sess:=sessions.Get(this.Ctx)

	user ,ok:=(sess.Get("user")).(*datamodels.UserModel)
	if !ok{
		this.Ctx.Application().Logger().Error("parse user err by sesssion")
	}


	//recv file info


	qetag :=this.Ctx.Request().FormValue("qetag")
	flowIdentifier :=this.Ctx.Request().FormValue("flowIdentifier")
	flowTotalChunks:=this.Ctx.Request().FormValue("flowTotalChunks")
	flowChunkNumber:=this.Ctx.Request().FormValue("flowChunkNumber")

	if len(qetag)==0 || len(flowIdentifier)==0|| len(flowTotalChunks)==0||len(flowChunkNumber)==0{
		this.Ctx.JSON(datamodels.RespModel{Status: 0,Msg:fmt.Sprintf("the upload info error")})
		return
	}


	//recv file stream
	this.Ctx.Request().ParseForm()
	file, fileHeader, err := this.Ctx.FormFile("file")
	if err != nil {
		this.Ctx.JSON(datamodels.RespModel{Status: 0,Msg:fmt.Sprintf("Failed to get data,err:%s\n", err.Error())})
		return
	}
	defer file.Close()

	fileMeta := datamodels.FileModel{
		FileName: fileHeader.Filename,
		Location: common.Local_Storage_Mount +user.Username+"/"+qetag+"-["+flowIdentifier+"]-"+ flowTotalChunks+"-"+flowChunkNumber,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	if  !common.Exists(common.Local_Storage_Mount +user.Username+"/"){
		err :=os.MkdirAll(common.Local_Storage_Mount +user.Username+"/",os.ModePerm)
		if err!=nil{
			panic(err)
		}
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		this.Ctx.JSON(datamodels.RespModel{Status: 0,Msg:		fmt.Sprintf("Failed to create file,err:%s\n", err.Error())})
		return
	}
	defer newFile.Close()
	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		fmt.Printf("Failed to save data into file,err:%s\n", err.Error())
		return
	}

}

func (this *FileController) PostUploadfinshed() {
	sess:=sessions.Get(this.Ctx)

	user ,ok:=(sess.Get("user")).(*datamodels.UserModel)
	if !ok{
		this.Ctx.Application().Logger().Error("parse user err by sesssion")
	}

	qetag :=this.Ctx.FormValue("qetag")
	flowIdentifier :=this.Ctx.FormValue("flowIdentifier")
	flowTotalChunks:=this.Ctx.FormValue("flowTotalChunks")
	fileExt:=this.Ctx.FormValue("fileExt")
	fileSize,_:= strconv.ParseInt(this.Ctx.Request().FormValue("fileSize"),10,0)
	fileName:=this.Ctx.Request().FormValue("fileName")

	if  !common.Exists(common.Local_Storage_Mount +user.Username+"/"+qetag+"-["+flowIdentifier+"]-"+ flowTotalChunks+"-"+"1"){
		this.Ctx.JSON(datamodels.RespModel{Status: 0,Msg: "file not exist"})
		return
	}
	newFileName := common.Local_Storage_Mount +user.Username+"/"+qetag+"-["+flowIdentifier+"]."+fileExt
	newFile, err := os.Create(newFileName)
	if err != nil {
		this.Ctx.JSON(datamodels.RespModel{Status: 0,Msg:		fmt.Sprintf("Failed to create file,err:%s\n", err.Error())})
		return
	}

	filesCount,err := strconv.ParseInt(flowTotalChunks,10,0)
	if err !=nil{
		panic(err)
	}
	for i :=1;i <= int(filesCount); i++ {
		oldFilePath := common.Local_Storage_Mount +user.Username+"/"+qetag+"-["+flowIdentifier+"]-"+ flowTotalChunks+"-"+strconv.Itoa(i)
		oldfile, err := os.Open(oldFilePath)
		if err != nil {
			oldfile.Close()
			newFile.Close()
			this.Ctx.JSON(datamodels.RespModel{Status: 0,Msg:fmt.Sprintf("Failed to open file,err:%s\n", err.Error())})
			return
		}

		_, err = io.Copy (newFile, oldfile)
		err = os.Remove(oldFilePath)
		if err!=nil{
			oldfile.Close()
			newFile.Close()
			this.Ctx.Application().Logger().Error("delete temp file error:"+err.Error())
			return
		}
		oldfile.Close()
	}
	newFile.Close()

	// calc the file qetag
	qetagBK,err :=common.GetEtag(newFileName)
	if err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg: err.Error(),
		})
		return
	}
	if qetagBK !=qetag{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg: "qetag inconsistent",
		})
		return
	}

	// save file info to db
	succ,err:=this.Service.AddFile(qetag,fileName,fileSize,newFileName)
	if !succ{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg: err.Error(),
		})
		return
	}

	succ,err =this.Service.AddUserFileRelation(user.Username, qetag,fileName,newFileName,fileSize)
	if !succ{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg: err.Error(),
		})
		return
	}


	this.Ctx.JSON(datamodels.RespModel{
		Status: 1,
		Msg: "upload success",
	})
}


func (this *FileController) GetUserindexfiles() {
	sess:=sessions.Get(this.Ctx)
	user ,ok:=(sess.Get("user")).(*datamodels.UserModel)
	if !ok{
		this.Ctx.Application().Logger().Error("parse user err by sesssion")
		return
	}
	files,err :=this.Service.QueryUserFils(user.Username)
	if err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg: "get user index files error",
		})
		return
	}
	this.Ctx.JSON(datamodels.RespModel{
		Status: 1,
		Msg: "OK",
		Data: files,
	})
}

func (this *FileController) GetDownloadfile() {

	fileqetag := this.Ctx.Request().FormValue("fileqetag")
	if len(strings.TrimSpace(fileqetag))==0{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg: "upload error:invalid fileqetag",
		})
		return
	}

	file,err:=this.Service.GetFileMeta(fileqetag)
	if err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg: "upload error:"+err.Error(),
		})
		return
	}

	//download
	this.Ctx.SendFile(file.Location,file.FileName)
}
