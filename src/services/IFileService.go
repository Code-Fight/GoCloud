package services

import (
	"errors"
	"gocloud/dao"
	"gocloud/datamodels"
	"math/big"
)

type IFileService interface {
	GetFileMeta(fileqetag string) (file *datamodels.FileModel,err error)
	AddFile(fileqetag string, filename string, filesize int64, fileaddr string) (succ bool,err error)
	QueryUserFils(username string,parent_dir ,status int64) (userfile []datamodels.UserFileModel, err error)
	QueryUserFilsByStatus(username string,status int64) (userfile []datamodels.UserFileModel, err error)
	GetUserFileByID(id int64) (userfile *datamodels.UserFileModel, err error)
	AddUserFileRelation(username, fileqetag, filename,fileaddr string, filesize ,is_dir,parent_dir int64) (succ bool,err error)
	DeleteFile(username, fileqetag string,parent_id int64)(succ bool,err error)
	UpdateUserFileName(id int64,name string)(succ bool,err error)
	GetUserDirByUser(user_name string,ignoreNode int)(dirs *[]map[string]interface{},err error)
	MoveFileTo(id, parent_dir int64) ( bool,  error)
	DeleteRecyle() ( bool,  error)


	CreateShareFile(user_file_id,share_id int64, pwd string)  (share_link string,succ bool, err error)
	QueryShareFileBy(qetag string) (share *datamodels.FileShareModel, err error)
	QueryShareFileAndValid(share_id,pwd string) (share *datamodels.FileShareModel, err error)
	QueryUserShareFileBy(share_id string) (share *datamodels.UserFileShareModel, err error)
	QueryUserShareFiles(user_name string) (share []datamodels.UserFileShareModel, err error)
	CancelShareFile(share_id string) (succ bool, err error)
	QueryShareFileByUserFileId(user_file_id int64) (share *datamodels.FileShareModel, err error)
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
func (this *fileService) QueryUserFils(username string,parent_dir ,status int64) (userfile []datamodels.UserFileModel, err error) {
	return this.dao.SelectUserFiles(username,parent_dir,status)
}
func (this *fileService) QueryUserFilsByStatus(username string,status int64) (userfile []datamodels.UserFileModel, err error) {
	return this.dao.SelectUserFilesByStatus(username,status)
}

func (this *fileService) GetUserFileByID(id int64) (userfile *datamodels.UserFileModel, err error){
	 return this.dao.SelectUserFilesByID(id)
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

func (this *fileService)  DeleteRecyle() ( bool,  error) {
	return this.dao.DeleteUserFile()
}


func (this *fileService)  CreateShareFile(user_file_id ,share_time int64, pwd string) (share_link string,succ bool, err error){
	id,err := this.dao.InsertShareFile(user_file_id,share_time,pwd)
	if err!=nil{
		return "",false,err
	}
	share_link =big.NewInt(id).Text(62)
	succ,err = this.dao.UpdateShareFileShareID(id,share_link)
	if !succ||err!=nil{
		return "",false,err
	}
	return

}

func (this *fileService)  QueryShareFileBy(share_id string) (share *datamodels.FileShareModel, err error){
	return this.dao.SelectShareFileBy(share_id)
}

func (this *fileService)  QueryUserShareFileBy(share_id string) (share *datamodels.UserFileShareModel, err error){
	return this.dao.SelectShareFileAndUserFile(share_id)
}

func (this *fileService)  QueryShareFileAndValid(share_id,pwd string) (share *datamodels.FileShareModel, err error){
	share ,err = this.dao.SelectShareFileBy(share_id)

	if err!=nil{
		return nil, err
	}

	if pwd !=share.SharePwd{
		return nil,errors.New("the share password invalid")
	}


	return share,nil
}

func (this *fileService)  QueryUserShareFiles(user_name string) (share []datamodels.UserFileShareModel, err error){

	return this.dao.SelectUserShareFiles(user_name)
}

func (this *fileService)  CancelShareFile(share_id string) (succ bool, err error){

	return this.dao.DeleteShareFileByID(share_id)
}

func (this *fileService)  QueryShareFileByUserFileId(user_file_id int64) (share *datamodels.FileShareModel, err error){

	return this.dao.SelectShareFileByUserFileId(user_file_id)
}



func createTree(tree *[]map[string]interface{},dirs []datamodels.UserFileModel,nodes  []datamodels.UserFileModel,ignoreNode int)  {

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

func getNode(data []datamodels.UserFileModel,parent_id int64) []datamodels.UserFileModel {
	temp :=[]datamodels.UserFileModel{}
	for _,v :=range data{
		if v.ParentDir == parent_id{
			temp = append(temp, v)
		}
	}
	return temp
}


