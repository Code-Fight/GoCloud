package datamodels

type UserFile struct {
	ID          int    `sql:"id"`
	UserName    string `sql:"user_name"`
	FileQetag   string `sql:"file_qetag"`
	FileName    string `sql:"file_name"`
	FileSize    int64  `sql:"file_size"`
	UploadAt    string `sql:"upload_at"`
	LastUpdated string `sql:"last_update"`
	IsDir       int64  `sql:"is_dir"`
	ParentDir   int64  `sql:"parent_dir"`
}
