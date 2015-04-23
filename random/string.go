package random

import (
	"crypto/rand"
	"encoding/hex"
)

func String(len int) string {
	buffer := make([]byte, len)
	if _, err := rand.Read(buffer); err != nil {
		panic(err)
	}

	return hex.EncodeToString(buffer)
}
