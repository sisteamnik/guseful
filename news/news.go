package news

import (
	"github.com/sisteamnik/guseful/rate"
	"github.com/sisteamnik/guseful/tags"
)

type (
	New struct {
		Id          int64
		Title       string
		Slug        string
		Description string
		Body        []byte
		ImgId       int64
		Views       int64
		OwnerId     int64
		Published   bool
		Source      string

		Deleted int64
		Created int64
		Updated int64
		Version int64

		ImgUrl string     `db:"-"`
		Tags   []tags.Tag `db:"-"`
		Images []int64    `db:"-"`
		Rate   rate.Rate  `db:"-"`
	}

	NewImages struct {
		NewId   int64
		ImageId int64
	}
)
