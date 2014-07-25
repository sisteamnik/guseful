package website

type (
	Site struct {
		Id       int64
		Domain   string
		Www      bool
		Https    bool
		Encoding string

		Allow bool

		ChangeFreq changefreq

		Description string
	}

	changefreq string

	SitePage struct {
		Id     int64
		SiteId int64
		Url    string
	}
)

const (
	WEEKLY changefreq = "weekly"
	DAILY             = "daily"
)
