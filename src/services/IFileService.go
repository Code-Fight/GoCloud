package services

import (
	"gocloud/dao"
	"gocloud/datamodels"
)

type IFileService interface {
	GetFileMeta(filehash string) (file *datamodels.FileModel,err error)
	OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) (succ bool,err error)
	QueryUserFileMetas(username string, limit int) (userfile []datamodels.UserFile, err error)
	OnUserFileUploadFinished(username, filehash, filename string, filesize int64) (succ bool,err error)
}

type fileService struct {
	dao dao.IFileDao
}

func NewFileService(dao dao.IFileDao) IFileService {
	return &fileService{dao}
}

func (this *fileService) GetFileMeta(filehash string) (file *datamodels.FileModel,err error){
	return this.dao.SelectFile(filehash)
}
func (this *fileService) OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) (succ bool,err error){
	return this.dao.InsertFile(filehash,filename,filesize,fileaddr)
}
func (this *fileService) QueryUserFileMetas(username string, limit int) (userfile []datamodels.UserFile, err error) {
	return this.dao.SelectUserFileMetas(username,limit)
}
func (this *fileService) OnUserFileUploadFinished(username, filehash, filename string, filesize int64) (succ bool,err error){
	return this.dao.InsertUserFile(username, filehash, filename ,filesize)
}