package parser

import (
	"github.com/coopernurse/gorp"
)

type (
	Parser struct {
		Id        int64
		Name      string
		Index     string //web address with links
		Type      string //i.e., "news" "article" "video" "photo
		LastVisit int64
		Rules     string //dom selectors in json i.e., {"date":"article.date"}

		Created int64
		Updated int64
		Deleted int64
		Version int64
	}

	Prsr struct {
		Opts Opts
	}

	Opts struct {
		db    *gorp.DbMap
		Delay int64
	}
)
