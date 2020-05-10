package controllers

import (
	"github.com/kataras/iris/v12"
	"gocloud/services"
)

type FileController struct {
	Ctx iris.Context
	Service services.IFileService
}



