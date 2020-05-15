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

func (this *fileService) GetUserDirByUser(user_name string)(dirs *[]map[string]interface{},err error){
	 user_dir,err:= this.dao.SelectUserDirs(user_name)
	 if err!=nil{
	 	return nil,err
	 }
	 dirs = &[]map[string]interface{}{}
	 root := getNode(user_dir,0)
	 CreateTree(dirs,user_dir,root)
	 return dirs,nil
}

func CreateTree(tree *[]map[string]interface{},dirs []datamodels.UserFile,nodes  []datamodels.UserFile,)  {

	for _,v :=range nodes{
		_node :=map[string]interface{}{}
		_node["id"] = v.ID
		_node["label"] = v.FileName
		*tree = append(*tree, _node)

		temp := getNode(dirs,int64(v.ID))
		if len(temp)>0{
			children :=&[]map[string]interface{}{}
			_node["children"] = children
			CreateTree(children,dirs,nodes)
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