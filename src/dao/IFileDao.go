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
	SelectFile(filehash string) (file *datamodels.FileModel,err error)
	InsertFile(filehash string, filename string, filesize int64, fileaddr string) (succ bool,err error)
	SelectUserFileMetas(username string, limit int) (userfile []datamodels.UserFile, err error)
	InsertUserFile(username, filehash, filename string, filesize int64) (succ bool,err error)
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

func (this *fileDao) SelectFile(filehash string) (file *datamodels.FileModel,err error) {
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err := this.mysqlConn.Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file " +
			"where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()


	row, errRow := stmt.Query(filehash)

	if errRow != nil {
		return nil, errRow
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return nil, err
	}

	common.DataToStructByTagSql(result, file)

	return file,nil
}

// InsertUserFile : 更新用户文件表
func (this *fileDao) InsertUserFile(username, filehash, filename string, filesize int64) (succ bool,err error) {
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err := this.mysqlConn.Prepare(
		"insert ignore into tbl_user_file (`user_name`,`file_sha1`,`file_name`," +
			"`file_size`,`upload_at`) values (?,?,?,?,?)")
	if err != nil {
		return false,err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		return false,err
	}
	return true,nil
}

// SelectUserFileMetas : 批量获取用户文件信息
func (this *fileDao) SelectUserFileMetas(username string, limit int) (userfile []datamodels.UserFile, err error) {
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err := this.mysqlConn.Prepare(
		"select file_sha1,file_name,file_size,upload_at," +
			"last_update from tbl_user_file where user_name=? limit ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}

	for _,v :=range result{
		temp :=datamodels.UserFile{}
		common.DataToStructByTagSql(v,temp)
		userfile = append(userfile, temp)
	}

	return userfile, nil
}


// InsertFile : 文件上传完成,保存meta
func (this *fileDao) InsertFile(filehash string, filename string, filesize int64, fileaddr string) (succ bool,err error) {
	if err = this.Conn(); err != nil {
		return
	}
	stmt, err :=this.mysqlConn.Prepare("insert ignore into tbl_file" +
		"(`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`) values(?,?,?,?,1)")
	if err != nil {
		fmt.Println("Failed to prepare statement,err:" + err.Error())
		return false,err
	}
	defer stmt.Close()
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false,err
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Println("File with hash been upload before", filehash)
		}
		return true,err
	}
	return false,err
}