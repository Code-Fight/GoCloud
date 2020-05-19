package services

import (
	"gocloud/datamodels"
	"testing"
)

func TestCreateTree(t *testing.T) {
	user_dir := []datamodels.UserFileModel{
		{FileName: "1",IsDir: 1,ID: 1,ParentDir: 0},
		{FileName: "2",IsDir: 1,ID: 2,ParentDir: 0},
		{FileName: "3",IsDir: 1,ID: 3,ParentDir: 0},
		{FileName: "11",IsDir: 1,ID: 11,ParentDir: 1},
		{FileName: "111",IsDir: 1,ID: 111,ParentDir: 11},
		{FileName: "22",IsDir: 1,ID: 22,ParentDir: 2},
	}
	dirs := &[]map[string]interface{}{}
	root := getNode(user_dir,0)
	createTree(dirs,user_dir,root,100)
}
