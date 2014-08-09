package md5

import (
	"crypto/md5"
	"fmt"
)

func Hash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
