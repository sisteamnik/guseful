package img

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/disintegration/imaging"
	"github.com/jteeuwen/imghash"
	"github.com/sisteamnik/guseful/chpu"
	"github.com/sisteamnik/guseful/md5"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	defaultloc string, sthost string, ssl bool) *Api {
	a := new(Api)
	a.dbname = dbname
	a.Db = db
	a.path = path
	a.nesting = net

	//it relative path of site where stored images start with "/"
	a.defautloc = defaultloc
	return a
}

func randName() string {
	return md5.Hash(fmt.Sprintf("%v", time.Now().UnixNano()))
}

func (a *Api) Create(data []byte, name, descr, imgtype string, lon, lat float64) (Img, error) {
	bts := bytes.NewReader(data)
	img, _, err := image.Decode(bts)
	if err != nil {
		return Img{}, err
	}

	tm := time.Now()

	if name == "" {
		name = randName()
	}

	name = chpu.Chpu(name)

	title := name

	if len(name) <= int(a.nesting) {
		name = randName()
	}

	//TODO bad code
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

	path := a.getPath(name)

	os.MkdirAll(path, 0777)

	out, err := os.Create(path + "/" + name + ".jpg")
	if err != nil {
		return Img{}, err
	}
	defer out.Close()
	jpeg.Encode(out, img, nil)

	//todo remove this(fly generated)
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
		Name:        title,
		Slug:        name,
		Description: descr,
		Created:     tm.UnixNano(),
		Updated:     tm.UnixNano(),
		Width:       img.Bounds().Dx(),
		Height:      img.Bounds().Dy(),
		Hash:        int64(imghash.Average(img)),
		Type:        imgtype,
	}

	err = a.Db.Insert(&im)
	if err != nil {
		fmt.Println(err)
		return Img{}, err
	}
	return im, nil
}

func (api *Api) GetName(a int64, args ...interface{}) string {
	return api.getName(a, args...)
}

//args[0] - bool webp
//args[1] - int xsize
//args[2] - int ysize
//args[3] - string resize type ("thumb","fit","crop")
func (api *Api) getName(a int64, args ...interface{}) string {
	if a == 0 {
		return ""
	}

	path, err := api.Db.SelectStr("select Slug from Img where Id = ?", a)
	if err != nil {
		panic(err)
	}
	p := api.defautloc + "/" + path
	fmt.Println(path)
	dest := api.getPath(path) + "/" + path + ".jpg"

	var xsize, ysize int
	if len(args) >= 3 {
		resize := "crop"
		if len(args) >= 4 {
			resize = args[3].(string)
		}
		xsize = args[1].(int)
		ysize = args[2].(int)
		p = setSize(p, xsize, ysize)
		api.Resize(dest, xsize, ysize, resize)
	}
	if len(args) >= 1 {
		if args[0].(bool) {
			api.WebPCopy(setSize(dest, xsize, ysize))
			return p + ".webp"
		}
	}

	return p + ".jpg"
}

func (api *Api) RealLoc(a int64) string {
	path, err := api.Db.SelectStr("select Slug from Img where Id = ?", a)
	if err != nil {
		panic(err)
	}
	return api.getPath(path) + "/" + path + ".jpg"
}

func (api *Api) Resize(imgloc string, xsize, ysize int, resizeType string) bool {
	dest := setSize(imgloc, xsize, ysize)
	if _, err := os.Stat(dest); err == nil {
		return true
	}
	bts, err := ioutil.ReadFile(imgloc)
	if err != nil {
		fmt.Println(err)
		return false
	}
	rdr := bytes.NewReader(bts)
	i, _, err := image.Decode(rdr)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var fsimg *image.NRGBA
	switch resizeType {
	case "fit":
		fsimg = imaging.Fit(i, xsize, ysize, imaging.Lanczos)
	case "thumb":
		fsimg = imaging.Thumbnail(i, xsize, ysize, imaging.Lanczos)
	default:
		fsimg = imaging.Resize(i, xsize, ysize, imaging.Lanczos)
	}
	out, err := os.Create(dest)
	if err != nil {
		return false
	}
	defer out.Close()
	jpeg.Encode(out, fsimg, nil)
	return true
}

func (api *Api) WebPCopy(imgloc string) bool {
	if api.WebpAddr != "" {
		dest, _ := getExt(imgloc)
		if _, err := os.Stat(dest + ".webp"); err == nil {
			return true
		}
		cmd := exec.Command(api.WebpAddr, imgloc, "-o",
			dest+".webp")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	return false
}

func (api *Api) getPath(name string) string {
	path := api.path

	fmt.Println(path)

	for i := uint8(0); i < api.nesting; i++ {
		fmt.Println(path)
		path += "/" + name[i:i+1]
	}
	return path
}

func getExt(s string) (string, string) {
	name := s
	ext := ""
	ar := strings.Split(s, ".")
	if len(ar) >= 2 {
		nar := ar[0 : len(ar)-1]
		name = strings.Join(nar, "")
		ext = ar[len(ar)-1]
	}
	return name, ext
}

func SetSize(s string, x, y int) string {
	return setSize(s, x, y)
}

func setSize(s string, x, y int) string {
	dest, ext := getExt(s)
	if x != 0 && y != 0 {
		dest = dest + "_" + strconv.Itoa(x) + "x" + strconv.Itoa(y)
	}
	if ext == "" {
		return dest
	}
	return dest + "." + ext
}
