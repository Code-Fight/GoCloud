package common

import (
	"crypto/sha1"
	"encoding/hex"
)



// Sha1 return the data sha1 string
func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}