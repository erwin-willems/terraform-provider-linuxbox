package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// StringToMd5 is responsible for generating an md5 sum of string (s)
func StringToMd5(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

// SetToMd5 is responsible for generating an md5 sum of a set of strings (s)
func SetToMd5(s []string) string {
	hasher := md5.New()
	for _, v := range s {
		hasher.Write([]byte(v))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
