package common

import (
	"os"
)

var(
	User_Pwd_Sha1_Salt = "(*32%$#"
	Local_Storage_Mount = ""
)

func init() {
	var mounterr error
	Local_Storage_Mount,mounterr = os.Getwd()
	if mounterr !=nil{
		panic("The Local_Storage_Mount ERROR")
	}
}

