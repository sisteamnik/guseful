package menu

type (
	Menu struct {
		Id    int64
		Title string

		Items []MenuItem `db:"-"`
	}
	MenuItem struct {
		Id       int64
		Title    string
		Url      string
		Position int64
		ParentId int64
		MenuId   int64

		Childs []MenuItem `db:"="`
	}
)
