package services

import (
	"errors"
	"gocloud/dao"
	"gocloud/datamodels"
)

type IFileService interface {
	GetFileMeta(fileqetag string) (file *datamodels.FileModel,err error)
	AddFile(fileqetag string, filename string, filesize int64, fileaddr string) (succ bool,err error)
	QueryUserFils(username string,parent_dir ,status int64) (userfile []datamodels.UserFile, err error)
	AddUserFileRelation(username, fileqetag, filename,fileaddr string, filesize ,is_dir,parent_dir int64) (succ bool,err error)
	DeleteFile(username, fileqetag string,parent_id int64)(succ bool,err error)
	UpdateUserFileName(id int64,name string)(succ bool,err error)
	GetUserDirByUser(user_name string,ignoreNode int)(dirs *[]map[string]interface{},err error)
	MoveFileTo(id, parent_dir int64) ( bool,  error)

	CreateShareFile(qetag, pwd string) (succ bool, err error)
	QueryShareFileBy(qetag string) (share *datamodels.FileShareModel, err error)
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
func (this *fileService) QueryUserFils(username string,parent_dir ,status int64) (userfile []datamodels.UserFile, err error) {
	return this.dao.SelectUserFiles(username,parent_dir,status)
}
func (this *fileService) AddUserFileRelation(username, fileqetag, filename,fileaddr string, filesize,is_dir,parent_dir int64) (succ bool,err error){

	// if the file was deleted,we only need to do confirm and update the file status
	file ,err :=this.dao.SelectUserFilesByQetag(username,fileqetag,parent_dir,3)
	if err!=nil{
		return false,nil
	}
	if file!=nil {
		return this.dao.UpdateUserFileStatus(1,int64(file.ID))
	}
	// else we insert a new record
	return this.dao.InsertUserFile(username, fileqetag, filename ,filesize,is_dir,parent_dir)
}

func (this *fileService) DeleteFile(username, fileqetag string,parent_id int64)(succ bool,err error){

	// if is a directory ,need it is empty
	file,err :=this.dao.SelectUserFiles(username,parent_id,1)
	if err!=nil{
		return
	}
	if len(file)>0{
		err=errors.New("the dir not empty!")
		return
	}
	// soft delete by update status
	return this.dao.UpdateUserFileStatus(3,parent_id)
}

func (this *fileService) UpdateUserFileName(id int64,name string)(succ bool,err error){
	return this.dao.UpdateUserFileName(id,name)
}

func (this *fileService) GetUserDirByUser(user_name string,ignoreNode int)(dirs *[]map[string]interface{},err error){
	 user_dir,err:= this.dao.SelectUserDirs(user_name)
	 if err!=nil{
	 	return nil,err
	 }
	 dirs = &[]map[string]interface{}{}
	 root := getNode(user_dir,0)
	 createTree(dirs,user_dir,root,ignoreNode)
	 return dirs,nil
}

func (this *fileService)  MoveFileTo(id, parent_dir int64) ( bool,  error) {
	return this.dao.UpdateUserFileParentDir(id,parent_dir)
}


func (this *fileService)  CreateShareFile(qetag, pwd string) (succ bool, err error){
	return this.dao.InsertShareFile(qetag,pwd)
}
func (this *fileService)  QueryShareFileBy(qetag string) (share *datamodels.FileShareModel, err error){
	return this.dao.SelectShareFileBy(qetag)
}


func createTree(tree *[]map[string]interface{},dirs []datamodels.UserFile,nodes  []datamodels.UserFile,ignoreNode int)  {

	for _,v :=range nodes{
		if v.ID == ignoreNode {
			continue
		}
		_node :=map[string]interface{}{}
		_node["id"] = v.ID
		_node["label"] = v.FileName
		*tree = append(*tree, _node)

		temp := getNode(dirs,int64(v.ID))
		if len(temp)>0{
			children :=&[]map[string]interface{}{}
			_node["children"] = children
			createTree(children,dirs,temp,ignoreNode)
		}
	}

}

func getNode(data []datamodels.UserFile,parent_id int64) []datamodels.UserFile {
	temp :=[]datamodels.UserFile{}
	for _,v :=range data{
		if v.ParentDir == parent_id{
			temp = append(temp, v)
		}
	}
	return temp
}

