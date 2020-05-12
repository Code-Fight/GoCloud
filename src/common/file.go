package common

import "os"


// Exists return true if the path exist else return false
func  Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func MergeFile()  {
	
}


