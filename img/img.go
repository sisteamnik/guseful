package img

import (
	"github.com/coopernurse/gorp"
)

type (
	Img struct {
		Id          int64
		Name        string
		Description string

		Named bool

		Deleted int64
		Updated int64
		Created int64
		Version int64
	}

	Api struct {
		path      string
		neting    uint8
		dbname    string
		defautloc string
		sizes     []Size

		Db *gorp.DbMap
	}
	Size struct {
		Width  int
		Height int
		Crop   string
	}
)
