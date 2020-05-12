package services

import (
	"gocloud/dao"
	"gocloud/datamodels"
)

type IFileService interface {
	GetFileMeta(fileqetag string) (file *datamodels.FileModel,err error)
	AddFile(fileqetag string, filename string, filesize int64, fileaddr string) (succ bool,err error)
	QueryUserFils(username string) (userfile []datamodels.UserFile, err error)
	AddUserFileRelation(username, fileqetag, filename,fileaddr string, filesize int64) (succ bool,err error)
}

type fileService struct {
	dao dao.IFileDao
}

func NewFileService(dao dao.IFileDao) IFileService {
	return &fileService{dao}
}

func (this *fileService) GetFileMeta(fileqetag string) (file *datamodels.FileModel,err error){
	return this.dao.SelectFile(fileqetag)
}
func (this *fileService) AddFile(fileqetag string, filename string, filesize int64, fileaddr string) (succ bool,err error){
	return this.dao.InsertFile(fileqetag,filename,filesize,fileaddr)
}
func (this *fileService) QueryUserFils(username string) (userfile []datamodels.UserFile, err error) {
	return this.dao.SelectUserFiles(username)
}
func (this *fileService) AddUserFileRelation(username, fileqetag, filename,fileaddr string, filesize int64) (succ bool,err error){

	return this.dao.InsertUserFile(username, fileqetag, filename ,filesize)
}