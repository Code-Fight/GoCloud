package datamodels

type FileModel struct {
	FileQetag string `sql:"file_qetag"`
	FileName string`sql:"file_name"`
	FileSize int64`sql:"file_size"`
	Location string`sql:"file_addr"`
	UploadAt string`sql:"update_at"`
}
