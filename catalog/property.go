package catalog

type PropertyName struct {
	Id            int64
	Name          string
	IsShortable   bool
	ShortTemplate string
	ParentId      int64
	Position      int64
}

type PropertyValue struct {
	Id       int64
	Value    string
	ParentId int64
}

type PropertiesBind struct {
	BindId          int64
	PropertyId      int64
	PropertyValueId int64
}

type Property struct {
	Id        int64
	Name      string
	Value     string
	Shortable bool
	Short     string
	ShortTemplate string

	Position  int64

	ParentId  int64
}

type ProductProperty struct {
	Id   int64
	Name string
	Position int64
	Childs []Property
}