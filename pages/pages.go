package pages

type Page struct {
	Id          int64
	Title       string
	Slug        string
	Description string
	Body        string
	Owner       int64
	PhotoId     int64

	Created int64
	Updated int64
	Deleted int64
	Version int64
}
