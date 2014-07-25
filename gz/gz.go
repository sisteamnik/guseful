package gz

import (
	"bytes"
	"compress/gzip"
)

func Gz(str string) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write([]byte(str))
	w.Close()
	return b.Bytes()
}
