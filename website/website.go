package website

import (
	"errors"
	"net/url"
	"strings"
)

type (
	Site struct {
		Id       int64
		Domain   string
		Www      bool
		Https    bool
		Encoding string
		Exist    bool

		Allow bool

		ChangeFreq string
		Pages      int64

		Description string
	}

	SitePage struct {
		Id     int64
		SiteId int64
		Url    string
		Error  int64

		Visited int64
		Body    []byte
	}
)

const (
	WEEKLY string = "weekly"
	DAILY         = "daily"
)
