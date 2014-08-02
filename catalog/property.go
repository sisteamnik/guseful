package catalog

type (
	PropertyName struct {
		Id       int64
		Value    string
		Position int64
		UnitId   int64
		GroupId  int64
	}

	PropertyValue struct {
		Id     int64
		Value  string
		NameId int64
	}

	PropertyUnit struct {
		Id      int64
		Value   string
		IsTrait bool
	}

	PropertyGroup struct {
		Id   int64
		Name string
	}

	PropertyBind struct {
		Id        int64
		ProductId int64
		NameId    int64
		ValueId   int64
	}

	Property struct {
		Id      int64
		NameId  int64
		ValueId int64
		UnitId  int64
		IsTrait bool
		GroupId int64

		Value    string
		Name     string
		Unit     string
		Position int64
	}

	ProductProperty struct {
		Id       int64
		Name     string
		Position int64
		Childs   []Property
	}

	ProductVendor struct {
		Id   int64
		Name string
	}

	ProductCollection struct {
		Id   int64
		Name string
	}

	ProductModel struct {
		Id   int64
		Name string
	}

	ProductCategory struct {
		Id       int64
		Name     string
		ParentId int64
	}
)
