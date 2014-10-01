package tags

type (
	Tag struct {
		Id          int64
		Title       string
		Description string
		ImgId       int64

		Owner int64

		Created int64
		Deleted int64
		Updated int64
		Version int64
	}

	Tags struct {
		Id       int64
		ItemId   int64
		ItemType string

		Owner int64

		Created int64
		Deleted int64
		Updated int64
		Version int64
	}
)
