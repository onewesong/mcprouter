package util

import (
	"crypto/md5"
	"fmt"
)

// MD5 returns the MD5 hash of the given string
func MD5(key string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}
