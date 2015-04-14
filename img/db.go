package img

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/disintegration/imaging"
	"github.com/sisteamnik/guseful/md5"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"strconv"
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
	defaultloc string, sizes []Size) *Api {
	a := new(Api)
	a.dbname = dbname
	a.Db = db
	a.path = path
	a.nesting = net
	a.defautloc = defaultloc
	a.sizes = sizes
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

	if name == "" || len(name) <= int(a.nesting) {
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

	for i := uint8(0); i < a.nesting; i++ {
		path += "/" + name[i:i+1]
	}

	os.MkdirAll(path, 0777)
	for _, v := range a.sizes {
		save_path := path + "/" + name + "_" + strconv.Itoa(v.Width) + "x" +
			strconv.Itoa(v.Height)
		save_path_with_ex := save_path + ".jpg"
		if v.Crop == "thumb" {
			c := imaging.Thumbnail(img, v.Width, v.Height, imaging.Lanczos)
			imaging.Save(c, save_path_with_ex)
		} else if v.Crop == "fit" {
			c := imaging.Fit(img, v.Width, v.Height, imaging.Lanczos)
			imaging.Save(c, save_path_with_ex)
		}
		if a.WebpAddr != "" {
			cmd := exec.Command(a.WebpAddr, save_path_with_ex, "-o",
				save_path+".webp")
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
				return Img{}, errors.New("Webp Not Working")
			}
		}

	}
	out, err := os.Create(path + "/" + name + ".jpg")
	if err != nil {
		return Img{}, err
	}
	defer out.Close()
	jpeg.Encode(out, img, nil)

	if a.WebpAddr != "" {
		cmd := exec.Command(a.WebpAddr, path+"/"+name+".jpg", "-o",
			path+"/"+name+".webp")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			return Img{}, errors.New("Webp Not Working")
		}
	}

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
