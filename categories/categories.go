package categories

type Category struct {
	Id     int64
	Name   string
	Parent int64
	Childs []Category `db:"-"`
}
