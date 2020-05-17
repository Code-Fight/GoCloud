package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"gocloud/common"
	"gocloud/datamodels"
	"time"
)

type IFileDao interface {
	Conn() error
	//tbale_file
	SelectFile(fileqetag string) (file *datamodels.FileModel,err error)
	InsertFile(fileqetag string, filename string, filesize int64, fileaddr string) (succ bool,err error)
	//table_user_file
	SelectUserFiles(username string,parent_dir,status int64) (userfile []datamodels.UserFile, err error)
	SelectUserDirs(username string) (userfile []datamodels.UserFile, err error)
	SelectUserFilesByQetag(username,fileqetag string,parent_dir,status int64) (userfile *datamodels.UserFile, err error)
	InsertUserFile(username, fileqetag, filename string, filesize,is_dir,parent_dir int64) (succ bool,err error)
	DeleteUserFile(username, fileqetag string)(succ bool,err error)
	UpdateUserFileStatus(status,id int64)(succ bool,err error)
	UpdateUserFileName(id int64,name string)(succ bool,err error)
	UpdateUserFileParentDir(id ,parent_dir int64)(succ bool,err error)
	//table_share_file
	InsertShareFile(qetag, pwd string) (succ bool, err error)
	SelectShareFileBy(qetag string) (share *datamodels.FileShareModel, err error)
}

type fileDao struct {
	mysqlConn *sql.DB
}

func NewFileDao( conn *sql.DB) IFileDao {
	return &fileDao{mysqlConn: conn}
}

func (this *fileDao) Conn() error {
	if this.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		this.mysqlConn = mysql
	}

	return nil
}

func (this *fileDao) SelectFile(fileqetag string) (file *datamodels.FileModel,err error) {
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err := this.mysqlConn.Prepare(
		"select file_qetag,file_addr,file_name,file_size,update_at from tbl_file " +
			"where file_qetag=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()


	row, errRow := stmt.Query(fileqetag)

	if errRow != nil {
		return nil, errRow
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return nil, err
	}
	file = &datamodels.FileModel{}
	common.DataToStructByTagSql(result, file)

	return file,nil
}

// InsertUserFile : Add the file info to tbl_user_file
func (this *fileDao) InsertUserFile(username, fileqetag, filename string, filesize,is_dir,parent_dir int64) (succ bool,err error) {
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err := this.mysqlConn.Prepare(
		"insert ignore into tbl_user_file (`user_name`,`file_qetag`,`file_name`," +
			"`file_size`,`upload_at`,`is_dir`,`parent_dir`) values (?,?,?,?,?,?,?)")
	if err != nil {
		return false,err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, fileqetag, filename, filesize, time.Now(),is_dir,parent_dir)
	if err != nil {
		return false,err
	}
	return true,nil
}

// SelectUserFiles : Get the user  first page files
func (this *fileDao) SelectUserFiles(username string,parent_dir,status int64) (userfile []datamodels.UserFile, err error) {
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err := this.mysqlConn.Prepare(
		"select id, file_qetag,file_name,file_size,upload_at," +
			"last_update,is_dir,parent_dir  from tbl_user_file where user_name=? and parent_dir =? and status=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username,parent_dir,status)
	if err != nil {
		return nil, err
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}

	for _,v :=range result{
		temp :=&datamodels.UserFile{}

		common.DataToStructByTagSql(v,temp)
		userfile = append(userfile, *temp)
	}

	return userfile, nil
}

// InsertFile : 文件上传完成,保存meta
func (this *fileDao) InsertFile(fileqetag string, filename string, filesize int64, fileaddr string) (succ bool,err error) {
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err :=this.mysqlConn.Prepare("insert ignore into tbl_file" +
		"(`file_qetag`,`file_name`,`file_size`,`file_addr`,`status`) values(?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement,err:" + err.Error())
		return false,err
	}
	defer stmt.Close()
	ret, err := stmt.Exec(fileqetag, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false,err
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Println("File with hash been upload before", fileqetag)
		}
		return true,err
	}
	return false,err
}

func (this *fileDao) DeleteUserFile(username, fileqetag string)(succ bool,err error){

	return false,errors.New("didn't impl")
}

func (this *fileDao) SelectUserFilesByQetag(username,fileqetag string,parent_dir,status int64) (userfile *datamodels.UserFile, err error){
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err := this.mysqlConn.Prepare(
		"select id, file_qetag,file_name,file_size,upload_at," +
			"last_update,is_dir,parent_dir  from tbl_user_file where user_name=? and parent_dir =? and status=? and file_qetag=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username,parent_dir,status,fileqetag)
	if err != nil {
		return nil, err
	}

	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return nil, err
	}
	userfile =&datamodels.UserFile{}
	common.DataToStructByTagSql(result,userfile)

	return userfile, nil
}


func (this *fileDao) UpdateUserFileStatus(status,id int64)(succ bool,err error){
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err :=this.mysqlConn.Prepare("update tbl_user_file set status = ?,last_update=? where" +
		" id=? ")
	if err != nil {
		fmt.Println("Failed to prepare statement,err:" + err.Error())
		return false,err
	}
	defer stmt.Close()
	_, err = stmt.Exec(status,time.Now(), id)
	if err != nil {
		fmt.Println(err.Error())
		return false,err
	}

	return true,nil
}

func (this *fileDao) UpdateUserFileName(id int64,name string)(succ bool,err error){
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err :=this.mysqlConn.Prepare("update tbl_user_file set file_name = ?,last_update=? where" +
		" id=? ")
	if err != nil {
		fmt.Println("Failed to prepare statement,err:" + err.Error())
		return false,err
	}
	defer stmt.Close()
	_, err = stmt.Exec(name,time.Now(), id)
	if err != nil {
		fmt.Println(err.Error())
		return false,err
	}

	return true,nil
}

func (this *fileDao) SelectUserDirs(username string) (userfile []datamodels.UserFile, err error){
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err := this.mysqlConn.Prepare(
		"select id, file_qetag,file_name,parent_dir from tbl_user_file where user_name=? and status=1 and is_dir=1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		return nil, err
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}

	for _,v :=range result{
		temp :=&datamodels.UserFile{}

		common.DataToStructByTagSql(v,temp)
		userfile = append(userfile, *temp)
	}

	return userfile, nil
}

func (this *fileDao) UpdateUserFileParentDir(id ,parent_dir int64)(succ bool,err error){
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err :=this.mysqlConn.Prepare("update tbl_user_file set parent_dir = ?,last_update=? where" +
		" id=? ")
	if err != nil {
		fmt.Println("Failed to prepare statement,err:" + err.Error())
		return false,err
	}
	defer stmt.Close()
	_, err = stmt.Exec(parent_dir,time.Now(), id)
	if err != nil {
		fmt.Println(err.Error())
		return false,err
	}

	return true,nil
}

func (this *fileDao) InsertShareFile(qetag, pwd string) (succ bool, err error) {
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err :=this.mysqlConn.Prepare("insert ignore into tbl_user_share_file" +
		"(`file_qetag`,`create_at`,`share_pwd`) values(?,?,?)")
	if err != nil {
		fmt.Println("Failed to prepare statement,err:" + err.Error())
		return false,err
	}
	defer stmt.Close()
	ret, err := stmt.Exec(qetag, time.Now(), pwd)
	if err != nil {
		fmt.Println(err.Error())
		return false,err
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Println("File with hash been shared")
		}
		return true,err
	}
	return false,err
}

func (this *fileDao) SelectShareFileBy(qetag string) (share *datamodels.FileShareModel, err error) {
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err :=this.mysqlConn.Prepare("insert ignore into tbl_user_share_file" +
		"(`file_qetag`,`create_at`,`share_pwd`) values(?,?,?)")
	if err != nil {
		fmt.Println("Failed to prepare statement,err:" + err.Error())
		return nil,err
	}
	defer stmt.Close()
	row, err := stmt.Query(qetag)
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}


	result := common.GetResultRow(row)
	if len(result) == 0 {
		return nil, err
	}

	share =&datamodels.FileShareModel{}
	common.DataToStructByTagSql(result,share)
    return share,nil

}