package datamodels

type FileShareModel struct {
	ID        int    `sql:"id"'`
	FileQetag string `sql:"file_qetag"`
	CreateAt  string `sql:"create_at"`
	SharePwd  string `sql:"share_pwd"`
}
