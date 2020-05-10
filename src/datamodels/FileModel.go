package datamodels

import "database/sql"

type FileModel struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}
