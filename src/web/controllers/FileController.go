package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"gocloud/common"
	"gocloud/datamodels"
	"gocloud/services"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type FileController struct {
	Ctx     iris.Context
	Service services.IFileService
}

func (this *FileController) PostUpload() {
	sess := sessions.Get(this.Ctx)

	user, ok := (sess.Get("user")).(*datamodels.UserModel)
	if !ok {
		this.Ctx.Application().Logger().Error("parse user err by sesssion")
	}

	//recv file info

	qetag := this.Ctx.Request().FormValue("qetag")
	flowIdentifier := this.Ctx.Request().FormValue("flowIdentifier")
	flowTotalChunks := this.Ctx.Request().FormValue("flowTotalChunks")
	flowChunkNumber := this.Ctx.Request().FormValue("flowChunkNumber")

	if len(qetag) == 0 || len(flowIdentifier) == 0 || len(flowTotalChunks) == 0 || len(flowChunkNumber) == 0 {
		this.Ctx.JSON(datamodels.RespModel{Status: 0, Msg: fmt.Sprintf("the upload info error")})
		return
	}

	//recv file stream
	this.Ctx.Request().ParseForm()
	file, fileHeader, err := this.Ctx.FormFile("file")
	if err != nil {
		this.Ctx.JSON(datamodels.RespModel{Status: 0, Msg: fmt.Sprintf("Failed to get data,err:%s\n", err.Error())})
		return
	}
	defer file.Close()

	fileMeta := datamodels.FileModel{
		FileName: fileHeader.Filename,
		Location: common.Local_Storage_Mount + user.Username + "/" + qetag + "-[" + flowIdentifier + "]-" + flowTotalChunks + "-" + flowChunkNumber,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	if !common.Exists(common.Local_Storage_Mount + user.Username + "/") {
		err := os.MkdirAll(common.Local_Storage_Mount+user.Username+"/", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		this.Ctx.JSON(datamodels.RespModel{Status: 0, Msg: fmt.Sprintf("Failed to create file,err:%s\n", err.Error())})
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
	sess := sessions.Get(this.Ctx)

	user, ok := (sess.Get("user")).(*datamodels.UserModel)
	if !ok {
		this.Ctx.Application().Logger().Error("parse user err by sesssion")
	}

	qetag := this.Ctx.FormValue("qetag")
	flowIdentifier := this.Ctx.FormValue("flowIdentifier")
	flowTotalChunks := this.Ctx.FormValue("flowTotalChunks")
	fileExt := this.Ctx.FormValue("fileExt")
	parent_dir, _ := strconv.ParseInt(this.Ctx.FormValue("parentDir"), 10, 0)
	fileSize, _ := strconv.ParseInt(this.Ctx.Request().FormValue("fileSize"), 10, 0)
	fileName := this.Ctx.Request().FormValue("fileName")

	if !common.Exists(common.Local_Storage_Mount + user.Username + "/" + qetag + "-[" + flowIdentifier + "]-" + flowTotalChunks + "-" + "1") {
		this.Ctx.JSON(datamodels.RespModel{Status: 0, Msg: "file not exist"})
		return
	}
	newFileName := common.Local_Storage_Mount + user.Username + "/" + qetag + "-[" + flowIdentifier + "]." + fileExt
	newFile, err := os.Create(newFileName)
	if err != nil {
		this.Ctx.JSON(datamodels.RespModel{Status: 0, Msg: fmt.Sprintf("Failed to create file,err:%s\n", err.Error())})
		return
	}

	filesCount, err := strconv.ParseInt(flowTotalChunks, 10, 0)
	if err != nil {
		panic(err)
	}
	for i := 1; i <= int(filesCount); i++ {
		oldFilePath := common.Local_Storage_Mount + user.Username + "/" + qetag + "-[" + flowIdentifier + "]-" + flowTotalChunks + "-" + strconv.Itoa(i)
		oldfile, err := os.Open(oldFilePath)
		if err != nil {
			_ = newFile.Close()
			_, _ = this.Ctx.JSON(datamodels.RespModel{Status: 0, Msg: fmt.Sprintf("Failed to open file,err:%s\n", err.Error())})
			return
		}

		_, err = io.Copy(newFile, oldfile)
		err = os.Remove(oldFilePath)
		if err != nil {
			_ = oldfile.Close()
			_ = newFile.Close()
			this.Ctx.Application().Logger().Error("delete temp file error:" + err.Error())
			return
		}
		_ = oldfile.Close()
	}
	_ = newFile.Close()

	// calc the file qetag
	qetagBK, err := common.GetEtag(newFileName)
	if err != nil {
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    err.Error(),
		})
		return
	}
	if qetagBK != qetag {
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "qetag inconsistent",
		})
		return
	}

	// save file info to db
	succ, err := this.Service.AddFile(qetag, fileName, fileSize, newFileName)
	if !succ {
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    err.Error(),
		})
		return
	}

	succ, err = this.Service.AddUserFileRelation(user.Username, qetag, fileName, newFileName, fileSize, 0, parent_dir)
	if !succ {
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    err.Error(),
		})
		return
	}

	this.Ctx.JSON(datamodels.RespModel{
		Status: 1,
		Msg:    "upload success",
	})
}

func (this *FileController) GetUserindexfiles() {
	sess := sessions.Get(this.Ctx)
	user, ok := (sess.Get("user")).(*datamodels.UserModel)
	if !ok {
		this.Ctx.Application().Logger().Error("parse user err by sesssion")
		return
	}

	parents, _ := strconv.ParseInt(this.Ctx.Request().FormValue("p"), 10, 0)

	files, err := this.Service.QueryUserFils(user.Username, parents,1)
	if err != nil {
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "get user index files error",
		})
		return
	}
	this.Ctx.JSON(datamodels.RespModel{
		Status: 1,
		Msg:    "OK",
		Data:   files,
	})
}

func (this *FileController) GetDownloadfileBy(filename string) {

	fileqetag := this.Ctx.Request().FormValue("fileqetag")
	if len(strings.TrimSpace(fileqetag)) == 0 {
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "upload error:invalid fileqetag",
		})
		return
	}

	file, err := this.Service.GetFileMeta(fileqetag)
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
	//this.Ctx.SendFile(file.Location,file.FileName)
}

func (this *FileController) GetCreatedirBy(parent_dir int,dir_name string) {

	sess := sessions.Get(this.Ctx)

	user, ok := (sess.Get("user")).(*datamodels.UserModel)
	if !ok {
		this.Ctx.Application().Logger().Error("parse user err by sesssion")
		return
	}

	// check the

	// it is a simple method to generate a 'dir sha1',
	// but I think it's enough
	// because the dir_name is unique
	dir_sha1 :=common.Sha1([]byte(user.Username+dir_name+strconv.Itoa(parent_dir)))
	succ, err := this.Service.AddUserFileRelation(user.Username, dir_sha1, dir_name, "", 0, 1,int64(parent_dir))
	if !succ {
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

func (this *FileController) GetDeleteBy(qetag string,parent_id int64) {
	sess := sessions.Get(this.Ctx)

	user, ok := (sess.Get("user")).(*datamodels.UserModel)
	if !ok {
		this.Ctx.Application().Logger().Error("parse user err by sesssion")
		return
	}
	succ,err:=this.Service.DeleteFile(user.Username,qetag,parent_id)
	if err !=nil||!succ{
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

func (this *FileController) PostFilesecondspass()  {
	sess := sessions.Get(this.Ctx)

	user, ok := (sess.Get("user")).(*datamodels.UserModel)
	if !ok {
		this.Ctx.Application().Logger().Error("parse user err by sesssion")
	}

	qetag := this.Ctx.FormValue("qetag")
	parent_dir, _ := strconv.ParseInt(this.Ctx.FormValue("parentDir"), 10, 0)
	fileName := this.Ctx.Request().FormValue("fileName")

	meta,err :=this.Service.GetFileMeta(qetag)
	if err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    err.Error(),
		})
		return
	}

	if meta==nil{
		//don't hava same file,need upload file
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "don't hava same file",
		})
		return
	}

	succ,err :=this.Service.AddUserFileRelation(user.Username,qetag,fileName,meta.Location,meta.FileSize,0,parent_dir)
	if err!=nil || !succ{
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

func (this *FileController) GetRenamefileBy(id int64, name string) {
	if id<=0{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    "the id error",
		})
		return
	}

	succ,err := this.Service.UpdateUserFileName(id,name)
	if err!=nil ||!succ{
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

func (this *FileController) GetUserdirsBy(id int)  {
	sess := sessions.Get(this.Ctx)

	user, ok := (sess.Get("user")).(*datamodels.UserModel)
	if !ok {
		this.Ctx.Application().Logger().Error("parse user err by sesssion")
		return
	}
	dirs,err:=this.Service.GetUserDirByUser(user.Username,id)
	if err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    err.Error(),
		})
		return
	}
	this.Ctx.JSON(datamodels.RespModel{
		Status: 1,
		Msg:    "OK",
		Data:   dirs,
	})

}

func (this *FileController) PostMovefile() {
	user_file_id,err :=strconv.ParseInt( this.Ctx.Request().FormValue("id"),10,0)
	dir ,err:=strconv.ParseInt(  this.Ctx.Request().FormValue("dir"),10,0)

	if err!=nil{
		this.Ctx.JSON(datamodels.RespModel{
			Status: 0,
			Msg:    err.Error(),
		})
		return
	}

	succ,err := this.Service.MoveFileTo(user_file_id,dir)
	if !succ||err!=nil{
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