package dao

import (
	"database/sql"
	"gocloud/common"
	"gocloud/datamodels"
)

type IUserDao interface {
	Conn() error
	Insert(*datamodels.UserModel) (int64, error)
	SelectByEmail(email string) (user *datamodels.UserModel,err error)
}


type userDao struct {
	table     string
	mysqlConn *sql.DB
}

func NewUserDao(table string, sql *sql.DB) IUserDao {
	return &userDao{table: table, mysqlConn: sql}
}

func (this *userDao) Conn() error {
	if this.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		this.mysqlConn = mysql
	}
	if this.table == "" {
		this.table = "tbl_user"
	}
	return nil
}

func (this *userDao) Insert(user *datamodels.UserModel) (ID int64, err error) {
	if err = this.Conn(); err != nil {
		return
	}

	sql := "insert  into tbl_user (`user_name`,`user_pwd`,`email`) values (?,?,?)"

	stmt, errStmt := this.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if errStmt != nil {
		return ID, errStmt
	}

	result, errResult := stmt.Exec(user.Username, user.Userpwd, user.Email)
	if errResult != nil {
		return ID, errResult
	}
	return result.LastInsertId()
}

func (this *userDao) SelectByEmail(email string) (user *datamodels.UserModel,err error){
	if errConn := this.Conn(); errConn != nil {
		return &datamodels.UserModel{}, errConn
	}

	sql := "Select * From "+this.table+" where email= ?"
	stmt, errStmt := this.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if errStmt != nil {
		return nil, errStmt
	}

	row, errRow := stmt.Query(email)

	if errRow != nil {
		return &datamodels.UserModel{}, errRow
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.UserModel{}, err
	}

	user = &datamodels.UserModel{}
	common.DataToStructByTagSql(result, user)
	return
}



