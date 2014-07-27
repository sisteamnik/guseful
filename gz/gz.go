package gz

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Gz(str string) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write([]byte(str))
	w.Close()
	return b.Bytes()
}

func U(b []byte) string {
	rdr := bytes.NewReader(b)
	r, _ := gzip.NewReader(rdr)
	bts, _ := ioutil.ReadAll(r)
	r.Close()
	return string(bts)
}
