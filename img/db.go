package img

import (
	"bytes"
	"github.com/coopernurse/gorp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func ExistName(Db *gorp.DbMap, name string) (bool, error) {
	id, err := Db.SelectInt("select Id from Img where Name = ?", name)
	if err != nil {
		return false, err
	}
	if id != 0 {
		return true, nil
	}
	return false, nil
}

func Create(Db *gorp.DbMap, data []byte, name string) (Img, error) {
	bts := bytes.NewReader(data)
	img, _, err := image.Decode(bts)
	if err != nil {
		return Img{}, err
	}

	return Img{}, nil
}
