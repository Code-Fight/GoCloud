package services

import (
	"gocloud/common"
	"gocloud/dao"
	"gocloud/datamodels"
)

type ILoginService interface {
	Signup(*datamodels.UserModel)(int64,error)
	Signin(*datamodels.UserModel)(*datamodels.UserModel,error)
}

type LoginService struct {
	dao dao.IUserDao
}

func NewLoginService(d dao.IUserDao) ILoginService {
	return &LoginService{d}
}

func (this *LoginService) Signup(user *datamodels.UserModel) (int64, error) {
	user.Userpwd = common.Sha1([]byte(user.Userpwd+common.User_Pwd_Sha1_Salt))
	return this.dao.Insert(user)
}

func (this *LoginService) Signin(user *datamodels.UserModel) (*datamodels.UserModel, error)  {
	return this.dao.SelectByEmail(user.Email)
}
