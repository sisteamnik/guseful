package eav

type (
	Db interface{}

	Eav struct {
		db Db
	}

	AttributeGroup struct {
		AttrId    int64
		GroupId   int64
		SortOrder int64
	}

	AttributeName struct {
		Id       int64
		AttrName string
	}
	AttribyteValue struct {
		AttributeNameId int64
		Value           string
		PriceAffectId   int64
		UnitId          int64
	}
	PriceAffect struct {
		Id       int64
		Addition int64
		Portion  float64
	}
	AttributeUnit struct {
		Id    int64
		Value string
	}
)

func setAttribute(name string, value string, unit string, category int64) {

}

func getAttribute() {

}
