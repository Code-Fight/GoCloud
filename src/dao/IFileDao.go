package dao

import (
	"database/sql"
	"fmt"
	"gocloud/common"
	"gocloud/datamodels"
	"time"
)

type IFileDao interface {
	Conn() error
	SelectFile(fileqetag string) (file *datamodels.FileModel,err error)
	InsertFile(fileqetag string, filename string, filesize int64, fileaddr string) (succ bool,err error)
	SelectUserFiles(username string) (userfile []datamodels.UserFile, err error)
	InsertUserFile(username, fileqetag, filename string, filesize int64) (succ bool,err error)
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
		"select file_qetag,file_addr,file_name,file_size from tbl_file " +
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
func (this *fileDao) InsertUserFile(username, fileqetag, filename string, filesize int64) (succ bool,err error) {
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err := this.mysqlConn.Prepare(
		"insert ignore into tbl_user_file (`user_name`,`file_qetag`,`file_name`," +
			"`file_size`,`upload_at`) values (?,?,?,?,?)")
	if err != nil {
		return false,err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, fileqetag, filename, filesize, time.Now())
	if err != nil {
		return false,err
	}
	return true,nil
}

// SelectUserFiles : Get the user  first page files
func (this *fileDao) SelectUserFiles(username string) (userfile []datamodels.UserFile, err error) {
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err := this.mysqlConn.Prepare(
		"select file_qetag,file_name,file_size,upload_at," +
			"last_update,is_dir,parent_dir  from tbl_user_file where user_name=? and parent_dir =0")
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