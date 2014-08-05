package news

type (
	New struct {
		Id          int64
		Title       string
		Url         string
		Description string
		Body        string
		ImgId       int64
		Views       int64
		OwnerId     int64
		Published   bool

		Deleted int64
		Created int64
		Updated int64
		Version int64

		ImgUrl string   `db:"-"`
		Images []string `db:"-"`
		Tags   []Tag    `db:"-"`
	}

	Tag struct {
		Id          int64
		Title       string
		Url         string
		Description string
		ImgId       int64

		Created int64
		Updated int64
		Version int64

		ImgUrl string `db:"-"`
	}

	Tags struct {
		NewId int64
		TagId int64
	}
)
