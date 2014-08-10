package img

import (
	"bytes"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/sisteamnik/guseful/md5"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"time"
)

func (a *Api) ExistName(name string) (bool, error) {
	id, err := a.Db.SelectInt("select Id from Img where Name = ?", name)
	if err != nil {
		return false, err
	}
	if id != 0 {
		return true, nil
	}
	return false, nil
}

func NewApi(db *gorp.DbMap, path string, net uint8, dbname,
	defaultloc string) *Api {
	a := new(Api)
	a.dbname = dbname
	a.Db = db
	a.path = path
	a.neting = net
	a.defautloc = defaultloc
	return a
}

func randName() string {
	return md5.Hash(fmt.Sprintf("%v", time.Now().UnixNano()))
}

func (a *Api) Create(data []byte, name string, descr string) (Img, error) {
	named := true
	bts := bytes.NewReader(data)
	img, _, err := image.Decode(bts)
	if err != nil {
		return Img{}, err
	}

	tm := time.Now()

	if name == "" {
		name = randName()
		named = false
	}

	for {
		found, err := a.ExistName(name)
		if err != nil {
			return Img{}, err
		}
		if found {
			name = randName()
		}
		if !found {
			break
		}
	}

	path := a.path

	for i := uint8(0); i < a.neting; i++ {
		path += "/" + name[i:i+1]
	}

	os.MkdirAll(path, 0777)

	out, err := os.Create(path + "/" + name + ".jpg")
	if err != nil {
		return Img{}, err
	}
	defer out.Close()
	jpeg.Encode(out, img, nil)

	im := Img{
		Name:        name,
		Description: descr,
		Named:       named,
		Created:     tm.UnixNano(),
		Updated:     tm.UnixNano(),
	}

	err = a.Db.Insert(&im)
	if err != nil {
		return Img{}, err
	}
	return im, nil
}
