package datamodels

type FileShareModel struct {
	ID        int    `sql:"id"'`
	ShareId string `sql:"share_id"`
	UserFileId int`sql:"user_file_id"`
	CreateAt  string `sql:"create_at"`
	SharePwd  string `sql:"share_pwd"`
	ShareTime int `sql:"share_time"`
}
