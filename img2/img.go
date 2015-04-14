package img

import (
	"github.com/coopernurse/gorp"
)

type Img struct {
	Id          int64
	Name        string
	Slug        string
	Description string

	Hash int64

	Type string

	//size
	Width  int
	Height int

	//coords
	Lon float64
	Lat float64

	//user Id
	Owner int64

	//photograph user Id
	Author int64

	Created int64
	Updated int64
	Deleted int64
	Version int64
}

type Api struct {
	path      string
	nesting   uint8
	dbname    string
	defautloc string
	sthost    string
	ssl       bool

	Db       *gorp.DbMap
	WebpAddr string
}
