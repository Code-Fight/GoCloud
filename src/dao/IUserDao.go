package dao

import (
	"gocloud/datamodels"
)

type IUserDao interface {
	Conn() error
	Insert(*datamodels.UserModel) (int64, error)
	SelectByEmail(email string) (user *datamodels.UserModel,err error)
}

