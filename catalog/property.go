package catalog

type Property struct {
	Id              int64
	ProductId       int64
	PropertyId      int64
	PropertyValueId int64
}

type PropertyName struct {
	Id   int64
	Name string
}

type PropertyValue struct {
	Id    int64
	Value string
}
